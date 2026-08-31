# ADR 0003: v1.0 compatibility and CLI upgrade policy

## Status

Accepted on 2026-09-01 (see [Issue #6](https://github.com/y-marui/alfred-clean-invisible-text/issues/6)).

## Context

`docs/specification.md` deferred the minimum Alfred and macOS versions to
"before the first release" (Issue #6). The CLI upgrade behavior — how a
user gets a newer pinned `go-clean-invisible-text` — was implemented
(`docs/dependency-policy.md`, `.github/workflows/check-cli-update.yml`) but
never written down as a committed contract.

## Decision

### Minimum Alfred version

Alfred 5 or later. The Workflow uses Universal Actions and Script Filter
JSON features available since Alfred 4, but targeting only the current
major version keeps the support matrix small for a solo-maintained project
and matches how most active Alfred users are already on Alfred 5.

### Minimum macOS version

Tracks the floor `go-clean-invisible-text` documents for its official
release binaries in
[ADR 0002](https://github.com/y-marui/go-clean-invisible-text/blob/main/docs/decisions/0002-v1-compatibility-and-support-policy.md),
rather than a second, independently maintained number. As of the CLI
version currently pinned in
[`internal/cliasset/pinned.txt`](../../internal/cliasset/pinned.txt) (built
with `go 1.27.0`), that floor is **macOS 13 (Ventura) or later**. Bumping
the pinned CLI version to one built with a newer Go toolchain can raise
this floor; `docs/dependency-policy.md`'s pin-update process is where that
gets caught, not this ADR.

### CLI upgrade behavior

Confirms the behavior already implemented in
`docs/dependency-policy.md`: the Workflow embeds a pinned, checksum-verified
CLI binary per architecture at packaging time and never downloads anything
at ordinary runtime. A user gets a newer CLI only by installing a newer
`.alfredworkflow` release. `.github/workflows/check-cli-update.yml` already
automates noticing a new upstream CLI release and opening a pull request to
bump the pin — no further runtime auto-update mechanism is planned.

## Consequences

- `docs/specification.md`'s compatibility line points here instead of
  saying the decision is still open.
- `README.md`/`README-jp.md` state the minimum Alfred/macOS versions so
  users can check compatibility before installing.
- A future CLI pin bump that raises the effective macOS floor (per
  upstream's own ADR 0002) is a user-facing change and belongs in
  `CHANGELOG.md`, not just the pin-update PR body.
- Alfred Gallery submission requirements are tracked separately in
  [docs/alfred-gallery-readiness.md](../alfred-gallery-readiness.md), since
  that checklist covers more than version support.
