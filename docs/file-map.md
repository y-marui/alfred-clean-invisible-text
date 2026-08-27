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

## Action orchestration and Alfred entry point (issue #4, Go layer)

| File | Role | Key Dependencies |
|---|---|---|
| `internal/scriptfilter/scriptfilter.go` | Alfred Script Filter JSON response types (`Item`, `Mod`, `Response`) | — |
| `internal/cliinvoke/cliinvoke.go` | Runs `check`/`explain`/`fix --json` via a pinned binary; parses the single-file JSON report; classifies exit code + `error` + `action=="warn"` findings into the Clean/Cleaned/Warning/Error `State` | `internal/scriptfilter`-independent (plain JSON) |
| `internal/action/action.go` | `List`/`Check`/`Reveal`/`Clean`/`CopyReport`/`BuildReport` — resolves input (`Request.resolve`, surfacing clipboard `ErrEmpty`/`ErrUnsupported` as an Alfred Error item), runs the CLI via `tempinput`+`cliinvoke`, writes the clipboard on Clean success, and attaches the cmd/shift modifier actions (copy report with text, re-run keeping warnings) | `internal/cliinvoke`, `internal/clipboard`, `internal/tempinput`, `internal/scriptfilter` |
| `cmd/clean-invisible-text-alfred/main.go` | The binary Alfred invokes; `list`/`run`/`copy-report` subcommands map 1:1 to planned `workflow/info.plist` nodes | `internal/action`, `internal/cliasset` |

The `workflow/info.plist` wiring that invokes this binary (Universal Action,
keyword, and the connections between nodes) does not exist yet.
