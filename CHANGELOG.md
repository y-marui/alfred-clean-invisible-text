# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- Alfred Workflow specification: entry points, the Check/Reveal/Clean/Copy
  report actions, the Clean/Cleaned/Warning/Error state model, and
  accessibility/keyboard flow ([docs/specification.md](docs/specification.md)).
- [ADR 0001](docs/decisions/0001-separate-cli-and-workflow.md): keep the CLI
  and Alfred Workflow in separate repositories.
- [ADR 0002](docs/decisions/0002-file-based-cli-invocation.md): invoke
  `check`/`explain`/`fix` against a temporary file, since only `clean` reads
  standard input.
