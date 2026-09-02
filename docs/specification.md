# Alfred Workflow Specification

## Scope

The Workflow provides Alfred UI and clipboard integration for the `clean-invisible-text` CLI. Unicode detection and cleaning behavior is defined exclusively by the CLI repository; this document never restates its category rules (see [CLI Dependency Policy](dependency-policy.md) and [ADR 0002](decisions/0002-file-based-cli-invocation.md)).

## Entry points

- **Universal Action** on selected plain text.
- **Keyword** against the system clipboard's plain-text content.

Both entry points lead to the same set of actions below. The exact keyword string and Universal Action registration are implementation detail tracked in [issue #4](https://github.com/y-marui/alfred-clean-invisible-text/issues/4).

## Actions

Each action invokes the CLI once against the input text, per [ADR 0002](decisions/0002-file-based-cli-invocation.md). `<input>` below means the selected text (Universal Action) or the clipboard's plain text (keyword).

### Check

Inspect `<input>` and show a one-line summary — findings exist or not — without writing anything. Invokes `check --json`. Never touches the clipboard.

### Reveal

Show every finding for `<input>` — code point, name, category, location (line/column), and planned action — without writing changes. Invokes `explain --json`. Presented as one Alfred result row per finding; an input with no findings shows a single "no findings" row. Never touches the clipboard.

### Clean

Run the CLI cleaner against `<input>` and, only after CLI success, atomically replace the clipboard's plain-text content with the cleaned result (see [Failure behavior](#failure-behavior)). Invokes `fix --json`, which removes Warn-classified code points by default like Block. An alternate action (a modifier key on the same result row) re-runs with `--keep-warnings` before writing to the clipboard, for use after a Warning state (see [States](#states)) shows the only findings were Warn-classified.

### Copy report

A follow-up action on a Check, Reveal, or Clean result — not a separate scan — that copies a structured report built only from the findings already returned by that action's CLI call (code points, names, categories, locations, actions, counts, and CLI/Workflow versions). The default report excludes the original and cleaned text. A distinct modifier action opts into a variant that includes the text, so that never happens by accident. Copying a report replaces the clipboard the same way Clean does: atomically, and only on success.

Rich text, images, files, and other clipboard representations are outside the initial scope.

## States

Every result maps to one of four states, driven by the CLI's exit status (`docs/specification.md` in go-clean-invisible-text: `0` no findings, `1` findings were detected or applied, `2` invalid input or I/O failure) and the JSON `error` field:

| State | Condition | Meaning |
|---|---|---|
| Clean | exit `0`, `error: null`, no findings | Input already has no findings. |
| Cleaned / Findings | exit `1`, `error: null`, no Warn-category finding | Findings existed and every one matched an explicit Allow/Block rule. |
| Warning | exit `1`, `error: null`, at least one Warn-category finding | At least one finding fell outside the explicit Allow/Block list (see the CLI's Warn category) and was removed by default; review the report or re-run Clean with "keep warnings". |
| Error | exit `2`, or any per-file JSON `error` non-null | The CLI could not process the input, or the Workflow itself failed (missing/unverified binary, clipboard has no plain text). The original clipboard is untouched. |

Every state is conveyed in the result row's title/subtitle text (e.g. "Clean", "Warning", "Error"), never by icon or color alone.

## Privacy

Text is processed locally. Input text must not be logged. Diagnostics may contain only code points, categories, locations, counts, versions, and non-sensitive errors. `check`/`explain`/`fix` require file arguments, so the Workflow's transport of `<input>` through a local temporary file — and that file's lifecycle — is a privacy-sensitive implementation detail defined in [issue #3](https://github.com/y-marui/alfred-clean-invisible-text/issues/3), not here.

## Failure behavior

On a CLI error, the original clipboard is retained and the Workflow reports failure. The Workflow must never replace the clipboard with partial output.

## Accessibility and keyboard flow

- All interaction happens through Alfred's native Script Filter and Universal Action UI. The Workflow builds no custom window, so arrow-key navigation, Enter to run, Tab, and Esc to cancel behave exactly as in any other Alfred list — there is no bespoke keyboard handling to design or test.
- Alternate actions (keep-warnings re-run, include-text report) are exposed as Alfred's standard modifier-key subtitles on the existing result row, not a secondary menu, so every action stays reachable with the keyboard alone.
- Esc must leave the clipboard exactly as it was found at any point in the flow; no side effect may depend on how the user exits.

## Architecture support

The packaged Workflow must support current Intel and Apple Silicon Macs — enforced by shipping a universal (`lipo`) binary and verifying its architectures at build time, not by requiring hands-on testing on physical Intel hardware. Given how rare Intel Macs are getting, verifying on real Intel hardware is optional/best-effort, not a release blocker; Apple Silicon is the primary hands-on verification target. Minimum Alfred and macOS versions are decided in [ADR 0003](decisions/0003-v1-compatibility-and-upgrade-policy.md).
