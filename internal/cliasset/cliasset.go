// Package cliasset tracks the pinned go-clean-invisible-text release this
// Workflow bundles (docs/dependency-policy.md) and selects/verifies the
// correct architecture's binary at runtime.
package cliasset

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed pinned.txt
var pinnedManifest string

// Pinned is the pinned CLI version and its per-architecture SHA-256
// checksums.
type Pinned struct {
	Version   string
	Checksums map[string]string // e.g. "darwin-arm64" -> lowercase hex SHA-256
}

// Loaded is the manifest embedded in this binary at build time
// (pinned.txt). A parse failure here is a build-time bug in a checked-in
// file, not a runtime condition to recover from.
var Loaded = mustParse(pinnedManifest)

func mustParse(s string) Pinned {
	p, err := parse(s)
	if err != nil {
		panic(err)
	}
	return p
}

func parse(s string) (Pinned, error) {
	p := Pinned{Checksums: map[string]string{}}
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return Pinned{}, fmt.Errorf("cliasset: malformed manifest line %q", line)
		}
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if key == "version" {
			p.Version = value
			continue
		}
		p.Checksums[key] = value
	}
	if p.Version == "" {
		return Pinned{}, fmt.Errorf("cliasset: manifest has no version")
	}
	return p, nil
}

// Version is the pinned CLI release tag, suitable for display in
// diagnostics (docs/dependency-policy.md: "The packaged CLI version is
// included in diagnostics").
func Version() string {
	return Loaded.Version
}

// BinaryName is the upstream release asset name for a GOOS/GOARCH pair,
// e.g. "clean-invisible-text-darwin-arm64". This is also the filename
// scripts/fetch-cli-binaries.sh stages each verified binary under.
func BinaryName(goos, goarch string) string {
	return fmt.Sprintf("clean-invisible-text-%s-%s", goos, goarch)
}

// CurrentArchKey returns the pinned-manifest key for the running process's
// architecture ("darwin-amd64" or "darwin-arm64"), or an error if this
// Workflow is not running on a supported Mac architecture.
func CurrentArchKey() (string, error) {
	if runtime.GOOS != "darwin" {
		return "", fmt.Errorf("cliasset: unsupported OS %q (this Workflow only runs on macOS)", runtime.GOOS)
	}
	switch runtime.GOARCH {
	case "amd64", "arm64":
		return "darwin-" + runtime.GOARCH, nil
	default:
		return "", fmt.Errorf("cliasset: unsupported architecture %q", runtime.GOARCH)
	}
}

// Verify checks that the file at path matches the pinned checksum for
// archKey. It never includes file content in an error, only checksums.
func (p Pinned) Verify(archKey, path string) error {
	want, ok := p.Checksums[archKey]
	if !ok {
		return fmt.Errorf("cliasset: no pinned checksum for %q", archKey)
	}
	got, err := sha256File(path)
	if err != nil {
		return fmt.Errorf("cliasset: hashing %s: %w", path, err)
	}
	if !strings.EqualFold(got, want) {
		return fmt.Errorf("cliasset: checksum mismatch for %s: got %s, want %s", archKey, got, want)
	}
	return nil
}

// ResolvePath returns the path to the pinned CLI binary for the current
// architecture within assetsDir, verified against the pinned checksum. It
// re-verifies on every call rather than trusting a prior packaging-time
// check, since the file lives outside this binary once bundled.
func (p Pinned) ResolvePath(assetsDir string) (string, error) {
	archKey, err := CurrentArchKey()
	if err != nil {
		return "", err
	}
	path := filepath.Join(assetsDir, BinaryName(runtime.GOOS, runtime.GOARCH))
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("cliasset: pinned binary not found: %w", err)
	}
	if err := p.Verify(archKey, path); err != nil {
		return "", err
	}
	return path, nil
}

// ResolvePath resolves the pinned CLI binary for the current architecture
// using the embedded manifest (Loaded).
func ResolvePath(assetsDir string) (string, error) {
	return Loaded.ResolvePath(assetsDir)
}

func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
