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
// content present before the suite runs cannot be restored — an accepted
// limitation of testing against the real pasteboard.
func TestMain(m *testing.M) {
	original, err := ReadPlainText()
	code := m.Run()
	// Restore for both a real value and ErrEmpty (an empty string is still
	// a text state we can faithfully write back). ErrUnsupported and other
	// errors mean the original content wasn't text, so it cannot be
	// restored through WritePlainText — leave the clipboard as the suite
	// left it in that case.
	if err == nil || err == ErrEmpty {
		_ = WritePlainText(original)
	}
	os.Exit(code)
}

// requireMacOS skips the test unless running on macOS with pbcopy, pbpaste,
// and osascript on PATH.
func requireMacOS(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "darwin" {
		t.Skip("clipboard package is macOS-only")
	}
	for _, bin := range []string{"pbcopy", "pbpaste", "osascript"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("%s not found", bin)
		}
	}
}

func TestReadPlainText_RoundTrip(t *testing.T) {
	requireMacOS(t)

	want := "line1\nline2​ZWSP"
	if err := WritePlainText(want); err != nil {
		t.Fatalf("WritePlainText: %v", err)
	}

	got, err := ReadPlainText()
	if err != nil {
		t.Fatalf("ReadPlainText: %v", err)
	}
	if got != want {
		t.Errorf("ReadPlainText() = %q, want %q", got, want)
	}
}

func TestReadPlainText_Empty(t *testing.T) {
	requireMacOS(t)

	if err := WritePlainText(""); err != nil {
		t.Fatalf("WritePlainText: %v", err)
	}

	_, err := ReadPlainText()
	if err != ErrEmpty {
		t.Errorf("ReadPlainText() error = %v, want ErrEmpty", err)
	}
}

func TestReadPlainText_Unsupported(t *testing.T) {
	requireMacOS(t)

	// Put a non-text-only representation on the pasteboard via AppleScript
	// (an icon file coerced to «class icns»), matching how the discriminator
	// was verified manually: a pasteboard with only an image representation.
	icon := "/System/Library/CoreServices/Automator Application Stub.app/Contents/Resources/ApplicationStub.icns"
	cmd := exec.Command("osascript", "-e",
		`set the clipboard to (read (POSIX file "`+icon+`") as «class icns»)`)
	if err := cmd.Run(); err != nil {
		t.Skipf("could not set image clipboard: %v", err)
	}

	_, err := ReadPlainText()
	if err != ErrUnsupported {
		t.Errorf("ReadPlainText() error = %v, want ErrUnsupported", err)
	}
}
