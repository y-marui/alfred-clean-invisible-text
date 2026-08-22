# ADR 0001: Keep the CLI and Alfred Workflow separate

## Status

Accepted on 2026-08-22.

## Decision

Maintain the cross-platform Go CLI in `y-marui/go-clean-invisible-text` and the Alfred frontend in this repository.

## Rationale

The CLI serves Windows, macOS, Raspberry Pi, pre-commit, and standard streams. Alfred packaging, UI, clipboard behavior, and macOS architecture selection have a separate lifecycle.

## Consequences

The Workflow pins and verifies a released CLI binary. Unicode behavior is specified and tested only in the CLI repository. Cross-repository work is linked through GitHub Issues and a shared Project.
