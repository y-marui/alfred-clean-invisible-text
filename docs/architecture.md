# Architecture

## Overview

An Alfred Workflow (Go): `cmd/clean-invisible-text-alfred` is the single
universal (amd64+arm64) binary `workflow/info.plist` invokes. It writes
selected/clipboard text to a temporary file, invokes the pinned
`clean-invisible-text` binary against it, and prints Alfred Script Filter
JSON (or, for Copy report, writes the clipboard directly).
`scripts/build-workflow.sh` packages all of this into a signed, notarised
`.alfredworkflow`, verified inside Alfred's own Workflow debugger on
Apple Silicon. Testing on real Intel hardware remains open but is
optional/best-effort, not a blocker (issue #4) — Intel Macs are
increasingly rare, and the universal binary is already verified via
`lipo` at build time.

## Entry Points

- `cmd/clean-invisible-text-alfred` — subcommands `list`, `run`,
  `copy-report`, one per `workflow/info.plist` node (see the package doc
  comment in `main.go`)

Two Alfred triggers reach it, per
[docs/specification.md](specification.md#entry-points), both fully wired
in `workflow/info.plist` with no manual setup required: a keyword (`cit`,
clipboard) and a Universal Action trigger (selected text) connected to the
keyword-less Script Filter node. See
[docs/alfred-workflow-notes/workflow-object-schema.md](alfred-workflow-notes/workflow-object-schema.md) for how the
Universal Action Trigger object's plist form was reverse-engineered (Alfred
doesn't document it) and a gotcha worth knowing if this wiring is ever
touched: two of this workflow's Script Filter nodes have no keyword, and
only one of them is the correct connection target.

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
