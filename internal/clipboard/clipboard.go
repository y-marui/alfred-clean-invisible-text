// Package clipboard writes the macOS general pasteboard's plain-text
// representation. It never logs the text it handles — see
// docs/specification.md Privacy. Reading the clipboard doesn't need
// package-level code: workflow/info.plist supplies it via Alfred's own
// {clipboard} placeholder.
package clipboard

import (
	"bytes"
	"fmt"
	"os/exec"
)

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
