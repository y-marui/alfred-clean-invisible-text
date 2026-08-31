# Architecture

## Overview

An Alfred Workflow (Go): `cmd/clean-invisible-text-alfred` is the single
universal (amd64+arm64) binary `workflow/info.plist` invokes. It writes
selected/clipboard text to a temporary file, invokes the pinned
`clean-invisible-text` binary against it, and prints Alfred Script Filter
JSON (or, for Copy report, writes the clipboard directly).
`scripts/build-workflow.sh` packages all of this into a `.alfredworkflow`.
Not yet done: verification inside Alfred's own Workflow debugger, testing
on real Intel and Apple Silicon Macs, and a signed release (issue #4/#6).

## Entry Points

- `cmd/clean-invisible-text-alfred` — subcommands `list`, `run`,
  `copy-report`, one per `workflow/info.plist` node (see the package doc
  comment in `main.go`)

Two Alfred triggers reach it, per
[docs/specification.md](specification.md#entry-points): a keyword (`cit`,
clipboard) is fully wired in `workflow/info.plist`. The Universal Action
(selected text) needs a one-time manual step in Alfred's own UI — see
[README.md](../README.md) Setup — since Alfred's Universal Action object
isn't something this project can generate reproducibly from source; its
downstream Script Filter node already exists in `workflow/info.plist`
(the one with no keyword set), ready for that connection.

## Directory Structure

| Directory | Role |
|---|---|
| `cmd/clean-invisible-text-alfred/` | The binary Alfred invokes; dispatches to `internal/action` and prints Script Filter JSON |
| `internal/action/` | Check/Reveal/Clean/Copy report orchestration and the Alfred result rows for each |
| `internal/cliinvoke/` | Runs the pinned CLI's `check`/`explain`/`fix --json`; classifies the Clean/Cleaned/Warning/Error state |
| `internal/scriptfilter/` | Alfred Script Filter JSON response types |
| `internal/clipboard/` | Reads/writes the macOS pasteboard's plain-text representation only; never logs content |
| `internal/tempinput/` | The private, single-use temp file `check`/`explain`/`fix` require as input ([ADR 0002](decisions/0002-file-based-cli-invocation.md)) |
| `internal/cliasset/` | Pinned CLI version/checksums and runtime binary selection ([docs/dependency-policy.md](dependency-policy.md)) |
| `workflow/` | `info.plist` (the Alfred object graph), `icon.png` |
| `scripts/fetch-cli-binaries.sh` | Downloads and verifies the pinned CLI release, stages it into `assets/bin/` (gitignored) |
| `scripts/build-workflow.sh` | Builds the universal binary and packages `workflow/` + `assets/bin/` into `dist/*.alfredworkflow` |
| `docs/` | Specification, dependency policy, ADRs |
| `docs/dev-charter/` | Shared dev-charter (`git subtree`) |

## Key Dependencies

| Library / Module | Purpose |
|---|---|
| [go-clean-invisible-text](https://github.com/y-marui/go-clean-invisible-text) | Pinned, checksum-verified CLI binary that performs all Unicode detection/cleaning ([docs/dependency-policy.md](dependency-policy.md)) |
| `pbcopy`/`pbpaste`/`osascript` (macOS system binaries) | Pasteboard read/write and plain-text type detection (`internal/clipboard`) — no third-party Go module |
