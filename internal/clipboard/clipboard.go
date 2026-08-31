// Package clipboard reads and writes the macOS general pasteboard's
// plain-text representation only. It never logs the text it handles — see
// docs/specification.md Privacy.
package clipboard

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

// ErrEmpty indicates the clipboard has a plain-text representation, but it
// is the empty string.
var ErrEmpty = errors.New("clipboard: plain text is empty")

// ErrUnsupported indicates the clipboard has no plain-text representation
// at all (e.g. it holds only an image).
var ErrUnsupported = errors.New("clipboard: no plain-text content available")

// ReadPlainText returns the general pasteboard's plain-text content. It
// returns ErrUnsupported if the pasteboard has no text representation, and
// ErrEmpty if the text representation is the empty string. Errors never
// include clipboard content.
func ReadPlainText() (string, error) {
	hasText, err := hasTextType()
	if err != nil {
		return "", fmt.Errorf("clipboard: checking content type: %w", err)
	}
	if !hasText {
		return "", ErrUnsupported
	}

	out, err := exec.Command("pbpaste").Output()
	if err != nil {
		return "", fmt.Errorf("clipboard: reading text: %w", err)
	}
	if len(out) == 0 {
		return "", ErrEmpty
	}
	return string(out), nil
}

// WritePlainText replaces the general pasteboard's content with text,
// discarding any other representation the pasteboard held. It either fully
// succeeds or leaves the pasteboard untouched — it never writes partial
// content.
func WritePlainText(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = bytes.NewReader([]byte(text))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clipboard: writing text: %w", err)
	}
	return nil
}

// hasTextType reports whether the general pasteboard currently has a
// plain-text (string) representation. pbpaste alone cannot distinguish "no
// text on the clipboard" from "the empty string is on the clipboard" — both
// exit 0 with empty output — so this checks the pasteboard's declared types
// instead of trying to read its content.
func hasTextType() (bool, error) {
	out, err := exec.Command("osascript", "-e", "clipboard info").Output()
	if err != nil {
		return false, err
	}
	return bytes.Contains(out, []byte("string")), nil
}
