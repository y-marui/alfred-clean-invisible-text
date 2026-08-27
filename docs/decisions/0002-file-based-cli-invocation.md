# ADR 0002: Invoke `check`/`explain`/`fix` against a temporary file

## Status

Accepted on 2026-08-27.

## Decision

For Check, Reveal, and Clean, the Workflow writes the input text (selected text or clipboard plain text) to a single-use local temporary file and invokes `check`/`explain`/`fix --json` against it. The Workflow does not attempt to pipe this text through standard input.

## Rationale

The CLI's `check`, `explain`, and `fix` subcommands take `FILE...` arguments; they do not read standard input ([CLI Contract](https://github.com/y-marui/go-clean-invisible-text/blob/main/docs/cli.md)). Only `clean` reads standard input, and it only ever streams cleaned bytes back — it never reports findings, counts, or locations. Reveal and Copy report require the structured findings that only `--json` on `check`/`explain`/`fix` provides, so a temporary file is the only way to reach that data for arbitrary in-memory text.

## Consequences

- Selected or clipboard text — potentially sensitive — is briefly written to local disk. Its creation, permissions, and guaranteed cleanup (including on CLI failure) are specified in [issue #3](https://github.com/y-marui/alfred-clean-invisible-text/issues/3), not here.
- Clean uses `fix --json` rather than the `clean` stdin/stdout command, so that a single CLI invocation yields both the cleaned content and the findings used for the Warning state and Copy report.
- If a future CLI release adds standard-input support to `check`/`explain`/`fix`, this ADR should be revisited; until then, the Workflow's local file lifecycle is part of its privacy contract.
