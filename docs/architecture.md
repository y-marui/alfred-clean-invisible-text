# Architecture

## Overview

No implementation exists yet. Planned shape: an Alfred Workflow (Go) that
writes selected/clipboard text to a temporary file, invokes a pinned
`clean-invisible-text` binary against it, and reflects the result in
Alfred's native UI — see [docs/specification.md](specification.md) and
[ADR 0002](decisions/0002-file-based-cli-invocation.md).

## Entry Points

None yet. Planned: a Universal Action (selected text) and an Alfred keyword
(clipboard) — see [docs/specification.md](specification.md#entry-points).

## Directory Structure

| Directory | Role |
|---|---|
| `docs/` | Specification, dependency policy, ADRs |
| `docs/dev-charter/` | Shared dev-charter (`git subtree`) |

## Key Dependencies

| Library / Module | Purpose |
|---|---|
| [go-clean-invisible-text](https://github.com/y-marui/go-clean-invisible-text) | Pinned, checksum-verified CLI binary that performs all Unicode detection/cleaning ([docs/dependency-policy.md](dependency-policy.md)) |
