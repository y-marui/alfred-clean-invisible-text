package cliinvoke

import (
	"os"
	"path/filepath"
	"testing"
)

func fakeCLIPath(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs("testdata/fakecli.sh")
	if err != nil {
		t.Fatalf("Abs: %v", err)
	}
	return abs
}

func TestRun_Clean(t *testing.T) {
	t.Setenv("FAKECLI_SCENARIO", "clean")
	f := writeTempFile(t, "hello")

	got, err := Run(fakeCLIPath(t), Check, f, false)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if got.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", got.ExitCode)
	}
	if len(got.File.Findings) != 0 {
		t.Errorf("Findings = %v, want none", got.File.Findings)
	}
	if got.State() != StateClean {
		t.Errorf("State() = %v, want StateClean", got.State())
	}
}

func TestRun_Cleaned(t *testing.T) {
	t.Setenv("FAKECLI_SCENARIO", "cleaned")
	f := writeTempFile(t, "hello")

	got, err := Run(fakeCLIPath(t), Fix, f, false)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if got.ExitCode != 1 {
		t.Errorf("ExitCode = %d, want 1", got.ExitCode)
	}
	if len(got.File.Findings) != 1 {
		t.Fatalf("Findings = %v, want 1", got.File.Findings)
	}
	if got.File.Findings[0].Action != "remove" {
		t.Errorf("Findings[0].Action = %q, want %q", got.File.Findings[0].Action, "remove")
	}
	if got.State() != StateCleaned {
		t.Errorf("State() = %v, want StateCleaned", got.State())
	}
}

func TestRun_Warning(t *testing.T) {
	t.Setenv("FAKECLI_SCENARIO", "warning")
	f := writeTempFile(t, "hello")

	got, err := Run(fakeCLIPath(t), Fix, f, false)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if got.State() != StateWarning {
		t.Errorf("State() = %v, want StateWarning", got.State())
	}
}

func TestRun_FileError(t *testing.T) {
	t.Setenv("FAKECLI_SCENARIO", "file-error")
	f := writeTempFile(t, "hello")

	got, err := Run(fakeCLIPath(t), Check, f, false)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if got.File.Error == nil {
		t.Fatal("File.Error = nil, want a message")
	}
	if got.State() != StateError {
		t.Errorf("State() = %v, want StateError", got.State())
	}
}

func TestRun_ProcessError(t *testing.T) {
	t.Setenv("FAKECLI_SCENARIO", "process-error")
	f := writeTempFile(t, "hello")

	_, err := Run(fakeCLIPath(t), Check, f, false)
	if err == nil {
		t.Fatal("Run() error = nil, want an error (no JSON on stdout to parse)")
	}
}

func TestRun_MalformedJSON(t *testing.T) {
	t.Setenv("FAKECLI_SCENARIO", "malformed-json")
	f := writeTempFile(t, "hello")

	_, err := Run(fakeCLIPath(t), Check, f, false)
	if err == nil {
		t.Fatal("Run() error = nil, want a JSON parse error")
	}
}

func TestRun_KeepWarningsFlag(t *testing.T) {
	// fakecli.sh doesn't branch on --keep-warnings, so this only verifies
	// Run() doesn't error when the flag is passed; the real CLI's handling
	// of --keep-warnings is out of scope here (ADR 0002 / dependency-policy.md).
	t.Setenv("FAKECLI_SCENARIO", "clean")
	f := writeTempFile(t, "hello")

	if _, err := Run(fakeCLIPath(t), Fix, f, true); err != nil {
		t.Fatalf("Run() with keepWarnings=true: %v", err)
	}
}

func TestState_String(t *testing.T) {
	cases := map[State]string{
		StateClean:   "Clean",
		StateCleaned: "Cleaned",
		StateWarning: "Warning",
		StateError:   "Error",
		State(99):    "Unknown",
	}
	for state, want := range cases {
		if got := state.String(); got != want {
			t.Errorf("State(%d).String() = %q, want %q", state, got, want)
		}
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "input.txt")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}
