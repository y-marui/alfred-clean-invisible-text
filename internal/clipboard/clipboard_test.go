package clipboard

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
)

// TestMain saves the real clipboard's plain-text content once before any
// test in this package runs, and restores it once after they all finish, so
// a local run doesn't clobber the developer's clipboard. Non-text clipboard
// content present before the suite runs cannot be restored this way — an
// accepted limitation of testing against the real pasteboard.
func TestMain(m *testing.M) {
	original, _ := exec.Command("pbpaste").Output()
	code := m.Run()
	_ = WritePlainText(string(original))
	os.Exit(code)
}

// requireMacOS skips the test unless running on macOS with pbcopy and
// pbpaste on PATH.
func requireMacOS(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "darwin" {
		t.Skip("clipboard package is macOS-only")
	}
	for _, bin := range []string{"pbcopy", "pbpaste"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("%s not found", bin)
		}
	}
}

func TestWritePlainText_RoundTrip(t *testing.T) {
	requireMacOS(t)

	want := "line1\nline2​ZWSP"
	if err := WritePlainText(want); err != nil {
		t.Fatalf("WritePlainText: %v", err)
	}

	got, err := exec.Command("pbpaste").Output()
	if err != nil {
		t.Fatalf("pbpaste: %v", err)
	}
	if string(got) != want {
		t.Errorf("pbpaste = %q, want %q", got, want)
	}
}
