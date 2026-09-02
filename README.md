# Clean Invisible Text for Alfred

> **This is the reference (English) version.**
> The canonical (Japanese) version is [README-jp.md](README-jp.md).

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui?style=social)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-donate-yellow.svg)](https://www.buymeacoffee.com/y.marui)

An Alfred Workflow for reviewing and cleaning dangerous invisible
Unicode characters locally.

> **Status:** [released](https://github.com/y-marui/alfred-clean-invisible-text/releases/latest)
> (signed and notarised); both entry points verified working on Apple
> Silicon — see the
> [roadmap](https://github.com/y-marui/alfred-clean-invisible-text/issues/1)
> for what's left (Intel hardware testing, optional/best-effort).

The Workflow is a thin macOS frontend for
[go-clean-invisible-text](https://github.com/y-marui/go-clean-invisible-text).
It does not implement Unicode cleaning rules independently — see
[ADR 0001](docs/decisions/0001-separate-cli-and-workflow.md).

## Requirements

- Alfred 5 or later
- macOS 13 (Ventura) or later, Intel or Apple Silicon — tracks the pinned
  CLI's own floor; see [ADR 0003](docs/decisions/0003-v1-compatibility-and-upgrade-policy.md)

## Setup

Download the signed
[latest release](https://github.com/y-marui/alfred-clean-invisible-text/releases/latest)
`.alfredworkflow` and double-click it to load it into Alfred. Nothing else
to configure — both entry points work immediately. (Building from source
is a contributor task; see [DEVELOPING.md](DEVELOPING.md).)

## Usage

Inspect and clean invisible Unicode characters in the clipboard's text via
the `cit` keyword.

![Check, Reveal, or Clean chooser](images/keyword-chooser.png)

Inspect and clean invisible Unicode characters in the currently selected
text via the Universal Action.

![Universal Action entry](images/universal-action.png)

Either entry point leads to the same choice (see
[docs/specification.md](docs/specification.md) for the full interaction
model, states, and accessibility notes):

* Check — inspect text and show a one-line summary, findings or not
* Reveal — show every finding (code point, name, category, location) without writing changes

  ![Reveal showing one finding](images/reveal-finding.png)
* Clean — run the CLI cleaner and replace the clipboard's plain-text content after success

  ![Clean result in the Warning state](images/clean-warning.png)

On a Check/Reveal/Clean result row:

* <kbd>↩︎</kbd> Copy a report of the findings, excluding the original text
* <kbd>⌘</kbd><kbd>↩︎</kbd> Copy a report that includes the original text
* <kbd>⇧</kbd><kbd>↩︎</kbd> Re-run Clean keeping unclassified characters instead of removing them (Warning state only)

## Documentation

- [docs/specification.md](docs/specification.md) — Alfred Workflow specification
- [docs/dependency-policy.md](docs/dependency-policy.md) — how the CLI is pinned and verified
- [docs/release-process.md](docs/release-process.md) — how a Workflow release is cut and published
- [docs/alfred-gallery-readiness.md](docs/alfred-gallery-readiness.md) — Alfred Gallery submission checklist
- [docs/decisions/](docs/decisions/) — architecture decision records (ADRs)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)

---
*This document has a Japanese canonical version [README-jp.md](README-jp.md). Update both in the same commit when editing.*
