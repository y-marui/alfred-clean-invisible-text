// Package tempinput provides the private, ephemeral temporary file that
// check/explain/fix require as input, since only clean reads standard input
// — see ADR 0002 (docs/decisions/0002-file-based-cli-invocation.md).
package tempinput

import (
	"fmt"
	"os"
)

// WithTempFile writes text to a new, single-use file with owner-only
// permissions, calls fn with that file's path, and removes the file
// afterward regardless of whether fn succeeds. text is never included in
// any error this function returns.
func WithTempFile(text string, fn func(path string) error) error {
	f, err := os.CreateTemp("", "alfred-clean-invisible-text-*")
	if err != nil {
		return fmt.Errorf("tempinput: creating temp file: %w", err)
	}
	path := f.Name()
	defer os.Remove(path)

	if err := f.Chmod(0o600); err != nil {
		f.Close()
		return fmt.Errorf("tempinput: setting temp file permissions: %w", err)
	}

	if _, err := f.WriteString(text); err != nil {
		f.Close()
		return fmt.Errorf("tempinput: writing temp file: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("tempinput: closing temp file: %w", err)
	}

	return fn(path)
}
