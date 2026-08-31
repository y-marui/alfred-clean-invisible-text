package cliasset

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	p, err := parse("# comment\nversion = v1.2.3\ndarwin-amd64 = abc123\n\ndarwin-arm64=def456\n")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if p.Version != "v1.2.3" {
		t.Errorf("Version = %q, want %q", p.Version, "v1.2.3")
	}
	want := map[string]string{"darwin-amd64": "abc123", "darwin-arm64": "def456"}
	for k, v := range want {
		if p.Checksums[k] != v {
			t.Errorf("Checksums[%q] = %q, want %q", k, p.Checksums[k], v)
		}
	}
}

func TestParse_Errors(t *testing.T) {
	cases := map[string]string{
		"malformed line":  "not-a-key-value-pair",
		"missing version": "darwin-amd64=abc123",
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := parse(input); err == nil {
				t.Error("parse() error = nil, want an error")
			}
		})
	}
}

func TestLoadedManifest(t *testing.T) {
	// The embedded pinned.txt must itself parse and cover both Mac
	// architectures — this is what mustParse() would have panicked on at
	// package init if it didn't.
	if Loaded.Version == "" {
		t.Error("Loaded.Version is empty")
	}
	for _, key := range []string{"darwin-amd64", "darwin-arm64"} {
		if Loaded.Checksums[key] == "" {
			t.Errorf("Loaded.Checksums[%q] is empty", key)
		}
	}
	if Version() != Loaded.Version {
		t.Errorf("Version() = %q, want %q", Version(), Loaded.Version)
	}
}

func TestBinaryName(t *testing.T) {
	got := BinaryName("darwin", "arm64")
	want := "clean-invisible-text-darwin-arm64"
	if got != want {
		t.Errorf("BinaryName() = %q, want %q", got, want)
	}
}

func TestCurrentArchKey(t *testing.T) {
	key, err := CurrentArchKey()
	if runtime.GOOS != "darwin" {
		if err == nil {
			t.Fatal("CurrentArchKey() error = nil on a non-macOS test runner, want an error")
		}
		return
	}
	if err != nil {
		t.Fatalf("CurrentArchKey(): %v", err)
	}
	if key != "darwin-amd64" && key != "darwin-arm64" {
		t.Errorf("CurrentArchKey() = %q, want darwin-amd64 or darwin-arm64", key)
	}
}

func writeTempBinary(t *testing.T, dir, name string, content []byte) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, content, 0o755); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}

func TestPinned_Verify(t *testing.T) {
	dir := t.TempDir()
	content := []byte("fake binary content")
	path := writeTempBinary(t, dir, "bin", content)

	sum, err := sha256File(path)
	if err != nil {
		t.Fatalf("sha256File: %v", err)
	}

	p := Pinned{Checksums: map[string]string{"test-arch": sum}}

	if err := p.Verify("test-arch", path); err != nil {
		t.Errorf("Verify() with matching checksum: %v", err)
	}
	if err := p.Verify("test-arch", path); err != nil {
		t.Errorf("Verify() called again (must not mutate state): %v", err)
	}

	p2 := Pinned{Checksums: map[string]string{"test-arch": "0000000000000000000000000000000000000000000000000000000000000000"}}
	if err := p2.Verify("test-arch", path); err == nil {
		t.Error("Verify() with mismatched checksum: error = nil, want an error")
	}

	if err := p.Verify("unknown-arch", path); err == nil {
		t.Error("Verify() with unknown arch key: error = nil, want an error")
	}
}

func TestPinned_ResolvePath(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("ResolvePath resolves a darwin-specific binary name")
	}

	archKey, err := CurrentArchKey()
	if err != nil {
		t.Fatalf("CurrentArchKey: %v", err)
	}

	dir := t.TempDir()
	content := []byte("fake pinned CLI binary")
	name := BinaryName(runtime.GOOS, runtime.GOARCH)
	writeTempBinary(t, dir, name, content)

	sum, err := sha256File(filepath.Join(dir, name))
	if err != nil {
		t.Fatalf("sha256File: %v", err)
	}
	p := Pinned{Version: "vTest", Checksums: map[string]string{archKey: sum}}

	got, err := p.ResolvePath(dir)
	if err != nil {
		t.Fatalf("ResolvePath: %v", err)
	}
	want := filepath.Join(dir, name)
	if got != want {
		t.Errorf("ResolvePath() = %q, want %q", got, want)
	}

	// Corrupting the file after it was staged must be caught, not trusted
	// from a prior packaging-time check.
	if err := os.WriteFile(want, []byte("tampered"), 0o755); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if _, err := p.ResolvePath(dir); err == nil {
		t.Error("ResolvePath() after tampering: error = nil, want a checksum error")
	} else if !strings.Contains(err.Error(), "checksum mismatch") {
		t.Errorf("ResolvePath() after tampering: error = %v, want a checksum mismatch error", err)
	}
}

func TestPinned_ResolvePath_MissingFile(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("ResolvePath resolves a darwin-specific binary name")
	}
	if _, err := Loaded.ResolvePath(t.TempDir()); err == nil {
		t.Error("ResolvePath() on an empty directory: error = nil, want an error")
	}
}
