# Contributing

## How to Contribute

Use GitHub Issues for roadmap work, temporary tasks, and decisions under
discussion. For large changes, open an issue before submitting a PR. Small
fixes and typos can be submitted directly as a PR.

Persist accepted Alfred behavior in [docs/specification.md](docs/specification.md)
or an ADR under [docs/decisions/](docs/decisions/) before closing the issue.
Closed issues are not the source of truth.

Unicode cleaning behavior belongs in
[go-clean-invisible-text](https://github.com/y-marui/go-clean-invisible-text)
and must not be reimplemented here.

## Development Setup

See [README.md](README.md) for setup instructions. There is no `DEVELOPING.md`
yet — no implementation exists (see [issue #4](https://github.com/y-marui/alfred-clean-invisible-text/issues/4)).

## Code Style

Follow [docs/dev-charter/CODE_STYLE.md](docs/dev-charter/CODE_STYLE.md).

## Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/) format
(e.g. `fix: ...`, `feat: ...`).

## Pull Request Checklist

See [.github/PULL_REQUEST_TEMPLATE.md](.github/PULL_REQUEST_TEMPLATE.md) for
the current checklist.
