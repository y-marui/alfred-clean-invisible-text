# Architecture

## Overview

An Alfred Workflow (Go): `cmd/clean-invisible-text-alfred` is the single
binary the packaged workflow invokes. It writes selected/clipboard text to a
temporary file, invokes the pinned `clean-invisible-text` binary against it,
and prints Alfred Script Filter JSON (or, for Copy report, writes the
clipboard directly). The `workflow/info.plist` wiring itself (issue #4) is
not built yet — only the fully-tested Go layer that plist will invoke.

## Entry Points

- `cmd/clean-invisible-text-alfred` — subcommands `list`, `run`,
  `copy-report`, one per planned `workflow/info.plist` node (see the
  package doc comment in `main.go`)

Two Alfred triggers reach it, per
[docs/specification.md](specification.md#entry-points): a Universal Action
(selected text) and a keyword (clipboard) — not wired yet (issue #4).

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
| `scripts/fetch-cli-binaries.sh` | Downloads and verifies the pinned CLI release, stages it into `assets/bin/` (gitignored) |
| `docs/` | Specification, dependency policy, ADRs |
| `docs/dev-charter/` | Shared dev-charter (`git subtree`) |

## Key Dependencies

| Library / Module | Purpose |
|---|---|
| [go-clean-invisible-text](https://github.com/y-marui/go-clean-invisible-text) | Pinned, checksum-verified CLI binary that performs all Unicode detection/cleaning ([docs/dependency-policy.md](dependency-policy.md)) |
| `pbcopy`/`pbpaste`/`osascript` (macOS system binaries) | Pasteboard read/write and plain-text type detection (`internal/clipboard`) — no third-party Go module |
