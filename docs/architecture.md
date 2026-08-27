# Architecture

## Overview

An Alfred Workflow (Go) that writes selected/clipboard text to a temporary
file, invokes a pinned `clean-invisible-text` binary against it, and
reflects the result in Alfred's native UI. No Alfred-facing entry point
exists yet (issue #4) — only the privacy-safe clipboard/temp-file primitives
so far.

## Entry Points

None yet. Planned: a Universal Action (selected text) and an Alfred keyword
(clipboard) — see [docs/specification.md](specification.md#entry-points).

## Directory Structure

| Directory | Role |
|---|---|
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
