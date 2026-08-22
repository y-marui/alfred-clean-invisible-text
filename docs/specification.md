# Alfred Workflow Specification

## Scope

The Workflow provides Alfred UI and clipboard integration for the `clean-invisible-text` CLI. Unicode detection and cleaning behavior is defined exclusively by the CLI repository.

## Actions

- Check: inspect selected plain text or clipboard plain text and display a summary.
- Reveal: display code points, names, counts, and locations without writing changes.
- Clean: run the CLI cleaner and replace the clipboard's plain-text content after success.
- Copy report: copy a report that excludes the original text by default.

Rich text, images, files, and other clipboard representations are outside the initial scope.

## Privacy

Text is processed locally. Input text must not be logged. Diagnostics may contain only code points, categories, locations, counts, versions, and non-sensitive errors.

## Failure behavior

On a CLI error, the original clipboard is retained and the Workflow reports failure. The Workflow must never replace the clipboard with partial output.

## Architecture support

The packaged Workflow must support current Intel and Apple Silicon Macs. The exact minimum macOS and Alfred versions will be decided before the first release.
