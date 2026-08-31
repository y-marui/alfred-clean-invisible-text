# Developing

## Build and Test

```bash
make build   # go build ./...
make test    # go test ./...
make lint    # gofmt -l . && go vet ./...
make fmt     # gofmt -w .
```

See the [Makefile](Makefile) for the exact commands each target runs.

## Requirements

- macOS (this Workflow only ever runs on macOS; the packages under
  `internal/` shell out to `pbcopy`/`pbpaste`/`osascript`)
- Go (see [go.mod](go.mod) for the minimum version)
- [pre-commit](https://pre-commit.com/) for the security and documentation
  hooks in [.pre-commit-config.yaml](.pre-commit-config.yaml)

## Project Layout

See [docs/architecture.md](docs/architecture.md) for directory structure and
[docs/file-map.md](docs/file-map.md) for file-level dependencies.

## Conventions

- The normative source for Workflow behavior is
  [docs/specification.md](docs/specification.md). A behavior change starts
  with a specification update (or an ADR under
  [docs/decisions/](docs/decisions/)), not the other way around.
- Unicode detection/cleaning behavior is never implemented here — see
  [docs/dependency-policy.md](docs/dependency-policy.md).
- `internal/clipboard` and `internal/tempinput` must never log or wrap
  clipboard/input text content into an error or log line — see
  [docs/specification.md](docs/specification.md) Privacy. Tests that exercise
  the real macOS pasteboard save and restore it (see `TestMain` in
  `internal/clipboard/clipboard_test.go`, and `saveAndRestoreClipboard` in
  `internal/action/action_test.go`) so a local `make test` run doesn't
  clobber the developer's clipboard.
- The real clipboard is a single OS resource shared across packages
  (`internal/clipboard`, `internal/action`). `go test`'s default
  cross-package parallelism races them against each other, so `make test`
  and CI both run `go test -p 1 ./...` (packages sequentially) — dropping
  `-p 1` reintroduces an intermittent, hard-to-reproduce-locally failure.
- Package names and file layout follow standard Go conventions
  (`internal/<package>/<file>.go`, `_test.go` alongside the code it tests).
- No comments that restate what the code does; comments explain non-obvious
  *why* only (see [docs/dev-charter/CODE_STYLE.md](docs/dev-charter/CODE_STYLE.md)).

## Pinned CLI Binaries

`make fetch-cli` downloads, checksum- and attestation-verifies, and stages
the pinned `go-clean-invisible-text` release into `assets/bin/` (gitignored).
Requires an authenticated `gh`. See
[docs/dependency-policy.md](docs/dependency-policy.md) for the trust model
and how to update the pin.

## Building the Workflow

```bash
make fetch-cli       # stage the pinned CLI (see above) — required first
make build-workflow  # → dist/*.alfredworkflow
```

`scripts/build-workflow.sh` builds `cmd/clean-invisible-text-alfred` as a
universal (amd64+arm64) binary via `lipo`, so a single package runs
natively on both Intel and Apple Silicon — verify with `lipo -info` on the
built binary if you touch that step. `workflow/info.plist` is the Alfred
object graph; edit it directly (there's no builder). Regenerate
`workflow/icon.png` (a placeholder) with
`go run scripts/tools/generate-icon.go workflow/icon.png`.

Because Alfred always runs a workflow's scripts with the working directory
set to the workflow's own bundle, `info.plist` invokes the binary via a
relative path (`./clean-invisible-text-alfred ...`) and the binary resolves
the pinned CLI via the relative `assets/bin/` — this only works when run
from inside the bundle (as verified below), not from an arbitrary `cwd`.

**Not yet done, and not something I can verify without Alfred itself**:
loading the built package into Alfred and testing via its Workflow
debugger, and running it on real Intel hardware (built and manually
verified here only via `lipo`/direct execution on Apple Silicon — see
[docs/file-map.md](docs/file-map.md) for exactly what was checked). The
Universal Action trigger also needs a one-time manual step in Alfred's own
UI — see [README.md](README.md) Setup.

## Commit Messages

[Conventional Commits](https://www.conventionalcommits.org/) format
(`feat:`, `fix:`, `docs:`, `refactor:`, `chore:`, `test:`).
