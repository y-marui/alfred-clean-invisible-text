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
  `internal/clipboard/clipboard_test.go`) so a local `make test` run doesn't
  clobber the developer's clipboard.
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

## Commit Messages

[Conventional Commits](https://www.conventionalcommits.org/) format
(`feat:`, `fix:`, `docs:`, `refactor:`, `chore:`, `test:`).
