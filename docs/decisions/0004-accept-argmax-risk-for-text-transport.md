# ADR 0004: Accept the ARG_MAX risk in Alfred's text transport

## Status

Accepted on 2026-09-03.

## Decision

Keep passing clipboard and selection text to the Workflow's own script
invocations as Alfred exposes them today — the Universal Action's selected
text as a command-line argument (`list selection "$1"`), and clipboard text
as an Alfred item variable that becomes an environment variable for the
`run` step (`os.Getenv("text")` in `main.go`). The Workflow will not add a
temporary-file transport for these two Alfred-glue steps.

This is distinct from [ADR 0002](0002-file-based-cli-invocation.md), which
already routes text through a private temporary file before invoking the
pinned CLI. That transport is unaffected by this decision.

## Rationale

macOS caps combined argv+environ size per process at `ARG_MAX` (`sysctl
kern.argmax`, ~1 MiB on current macOS). Both the selected-text argument and
the clipboard-text environment variable are subject to this cap when Alfred
spawns the Workflow's script. Realistic input for a Unicode-invisible-text
checker — a clipboard paste or a text selection a user is inspecting by
hand — is almost always well under that limit; the failure mode this ADR
accepts is specific to unusually large inputs (hundreds of thousands of
words in one selection or clipboard entry).

Closing this gap fully would require routing both entry points through a
private temp file ahead of the existing Script Filter/Universal Action
steps, touching `workflow/info.plist`, `internal/action`, and both
`cmd/clean-invisible-text-alfred` script invocations — a redesign of the
Alfred-glue layer disproportionate to a failure mode this unlikely to be
hit in practice.

## Consequences

- For the Universal Action selection path, an oversized selection can
  prevent Alfred from even spawning the `list selection` script. There is
  no Script Filter JSON surface at that point, so the Workflow cannot show
  its own Error state for this specific failure — it surfaces as an opaque
  Alfred-level failure instead. See docs/specification.md Failure behavior.
- For the clipboard/keyword path, an oversized clipboard could similarly
  exceed `ARG_MAX` when Alfred exports the `text` variable to the `run`
  step's environment.
- In both cases the original clipboard is never touched, since the failure
  happens before the Workflow's own Clean logic runs.
- If a future change needs to support much larger inputs, revisit this
  ADR — a temp-file transport for the Alfred-glue layer is the identified
  path.

## Related

[Issue #33](https://github.com/y-marui/alfred-clean-invisible-text/issues/33)
raised this; Copilot originally flagged it in review of
[PR #32](https://github.com/y-marui/alfred-clean-invisible-text/pull/32#discussion_r3922804631).
