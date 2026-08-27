package tempinput

import (
	"errors"
	"os"
	"testing"
)

func TestWithTempFile_ContentAndPermissions(t *testing.T) {
	const want = "line1\nline2​ZWSP"

	var gotPath string
	var gotContent []byte
	err := WithTempFile(want, func(path string) error {
		gotPath = path
		info, statErr := os.Stat(path)
		if statErr != nil {
			return statErr
		}
		if perm := info.Mode().Perm(); perm != 0o600 {
			t.Errorf("temp file permissions = %o, want 0600", perm)
		}
		var readErr error
		gotContent, readErr = os.ReadFile(path)
		return readErr
	})
	if err != nil {
		t.Fatalf("WithTempFile: %v", err)
	}
	if string(gotContent) != want {
		t.Errorf("temp file content = %q, want %q", gotContent, want)
	}

	if _, statErr := os.Stat(gotPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Errorf("temp file %q still exists after WithTempFile returned", gotPath)
	}
}

func TestWithTempFile_RemovesFileOnFnError(t *testing.T) {
	wantErr := errors.New("boom")

	var gotPath string
	err := WithTempFile("text", func(path string) error {
		gotPath = path
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("WithTempFile() error = %v, want %v", err, wantErr)
	}

	if _, statErr := os.Stat(gotPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Errorf("temp file %q still exists after fn returned an error", gotPath)
	}
}
