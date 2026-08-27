# File Map

_Last updated: 2026-08-27_

## Privacy-safe clipboard handling (issue #3)

| File | Role | Key Dependencies |
|---|---|---|
| `internal/clipboard/clipboard.go` | `ReadPlainText`/`WritePlainText` against the macOS general pasteboard; distinguishes "no text on clipboard" (`ErrUnsupported`) from "empty text" (`ErrEmpty`) via `osascript -e 'clipboard info'`, since `pbpaste` alone can't tell them apart | `os/exec` (`pbcopy`, `pbpaste`, `osascript`) |
| `internal/tempinput/tempinput.go` | `WithTempFile` — owner-only-permission, single-use temp file, guaranteed removal | `os` |

No Alfred-facing entry point exists yet (issue #4).
