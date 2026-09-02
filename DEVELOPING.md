# Developing

## Build and Test

```bash
make build   # go build ./...
make test    # go test ./...
make lint    # gofmt -l . && go vet ./...
make fmt     # gofmt -w .
```

See the [Makefile](Makefile) for the exact commands each target runs.

## Building the Workflow Package

Most users should just download the signed release (see README.md Setup).
To build the `.alfredworkflow` package from source instead — e.g. to pick
up a newer CLI pin before it's released, or to test a local change:

```bash
git clone https://github.com/y-marui/alfred-clean-invisible-text
cd alfred-clean-invisible-text
make fetch-cli       # downloads and verifies the pinned CLI binaries
make build-workflow  # → dist/*.alfredworkflow
```

Double-click the resulting `.alfredworkflow` to load it into Alfred. A
local build is unsigned (see
[docs/alfred-gallery-readiness.md](docs/alfred-gallery-readiness.md) for
what signing requires); only tagged releases built by CI are signed and
notarised.

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

A weekly [check-cli-update.yml](.github/workflows/check-cli-update.yml)
workflow opens a pull request when a newer upstream release exists; the pin
itself is otherwise never updated automatically.

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

Both entry points (the `cit` keyword and the Universal Action trigger) are
fully wired in `workflow/info.plist` and load with no manual setup — see
[docs/alfred-workflow-notes/workflow-object-schema.md](docs/alfred-workflow-notes/workflow-object-schema.md) for how
the Universal Action Trigger object's plist form was reverse-engineered.
Both have been verified working, including via Alfred's own Workflow
debugger, on Apple Silicon (see [docs/file-map.md](docs/file-map.md) for
exactly what was checked). Running on real Intel hardware remains
optional/best-effort, not a blocker — Intel Macs are increasingly rare,
and the universal binary is already verified via `lipo` at build time.

## Commit Messages

[Conventional Commits](https://www.conventionalcommits.org/) format
(`feat:`, `fix:`, `docs:`, `refactor:`, `chore:`, `test:`).
