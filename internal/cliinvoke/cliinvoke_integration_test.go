package cliinvoke

import (
	"os"
	"testing"

	"github.com/y-marui/alfred-clean-invisible-text/internal/cliasset"
)

// TestRun_RealPinnedCLI exercises Run against the actual pinned CLI binary
// staged by `make fetch-cli` (internal/cliasset), rather than the
// testdata/fakecli.sh stand-in used elsewhere in this package. It skips
// itself when assets/bin/ hasn't been populated, so `go test ./...` stays
// runnable offline — see DEVELOPING.md.
func TestRun_RealPinnedCLI(t *testing.T) {
	binaryPath, err := cliasset.ResolvePath("../../assets/bin")
	if err != nil {
		t.Skipf("pinned CLI binary not staged (run `make fetch-cli`): %v", err)
	}

	clean := writeTempFile(t, "hello world")
	got, err := Run(binaryPath, Check, clean, false)
	if err != nil {
		t.Fatalf("Run(check) on clean input: %v", err)
	}
	if got.State() != StateClean {
		t.Errorf("State() = %v, want StateClean for plain ASCII input", got.State())
	}

	dirty := writeTempFile(t, "hello​world") // ZERO WIDTH SPACE
	got, err = Run(binaryPath, Fix, dirty, false)
	if err != nil {
		t.Fatalf("Run(fix) on dirty input: %v", err)
	}
	if got.State() != StateCleaned {
		t.Errorf("State() = %v, want StateCleaned for a ZWSP finding", got.State())
	}
	if !got.File.Changed {
		t.Error("File.Changed = false, want true after fix removed the ZWSP")
	}
	cleaned, err := os.ReadFile(dirty)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(cleaned) != "helloworld" {
		t.Errorf("fixed file content = %q, want %q", cleaned, "helloworld")
	}
}
