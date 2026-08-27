# Clean Invisible Text for Alfred

> **This is the reference (English) version.**
> The canonical (Japanese) version is [README-jp.md](README-jp.md).

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/alfred-clean-invisible-text/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui?style=social)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-donate-yellow.svg)](https://www.buymeacoffee.com/y.marui)

A planned Alfred Workflow for reviewing and cleaning dangerous invisible
Unicode characters locally.

> **Status:** specification and roadmap. The Workflow is not implemented yet
> — see the [roadmap](https://github.com/y-marui/alfred-clean-invisible-text/issues/1).

The Workflow is a thin macOS frontend for
[go-clean-invisible-text](https://github.com/y-marui/go-clean-invisible-text).
It does not implement Unicode cleaning rules independently — see
[ADR 0001](docs/decisions/0001-separate-cli-and-workflow.md).

## Setup

Not available yet. Once the first release ships, installation will be:
double-click a signed `.alfredworkflow` to load it into Alfred (issue #4
tracks building and verifying that artifact).

## Usage

Planned actions, triggered via a Universal Action on selected text or a
keyword against the clipboard (see
[docs/specification.md](docs/specification.md) for the full interaction
model, states, and accessibility notes):

| Action | Description |
|---|---|
| Check | Inspect text and show a one-line summary — findings or not |
| Reveal | Show every finding (code point, name, category, location) without writing changes |
| Clean | Run the CLI cleaner and replace the clipboard's plain-text content after success |
| Copy report | Copy a structured report of findings, excluding the original text by default |

## Documentation

- [docs/specification.md](docs/specification.md) — Alfred Workflow specification
- [docs/dependency-policy.md](docs/dependency-policy.md) — how the CLI is pinned and verified
- [docs/decisions/](docs/decisions/) — architecture decision records (ADRs)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)

---
*This document has a Japanese canonical version [README-jp.md](README-jp.md). Update both in the same commit when editing.*
