# File Map

_Last updated: 2026-08-27_

## Privacy-safe clipboard handling (issue #3)

| File | Role | Key Dependencies |
|---|---|---|
| `internal/clipboard/clipboard.go` | `ReadPlainText`/`WritePlainText` against the macOS general pasteboard; distinguishes "no text on clipboard" (`ErrUnsupported`) from "empty text" (`ErrEmpty`) via `osascript -e 'clipboard info'`, since `pbpaste` alone can't tell them apart | `os/exec` (`pbcopy`, `pbpaste`, `osascript`) |
| `internal/tempinput/tempinput.go` | `WithTempFile` — owner-only-permission, single-use temp file, guaranteed removal | `os` |

## Pinned CLI binaries (issue #5)

| File | Role | Key Dependencies |
|---|---|---|
| `internal/cliasset/pinned.txt` | The pinned CLI release tag and per-architecture SHA-256 checksums — the trust anchor | — |
| `internal/cliasset/cliasset.go` | Parses `pinned.txt` (via `go:embed`); resolves and re-verifies the current architecture's staged binary at runtime | `pinned.txt` |
| `scripts/fetch-cli-binaries.sh` | Downloads the pinned release's darwin binaries, verifies checksum + attestation, stages them under `assets/bin/` | `internal/cliasset/pinned.txt`, `gh` |

No Alfred-facing entry point exists yet (issue #4).
