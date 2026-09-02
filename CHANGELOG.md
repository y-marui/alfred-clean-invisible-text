# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Fixed

- `workflow/info.plist`: the Universal Action trigger is now wired in
  source and works immediately on import — no more one-time manual setup
  in Alfred's own UI. The earlier assumption that this object couldn't be
  authored reproducibly was never actually tested; exporting a
  manually-wired workflow and diffing its plist showed it needed only one
  static object and one connection, no machine-specific state. See
  [docs/alfred-workflow-notes/workflow-object-schema.md](docs/alfred-workflow-notes/workflow-object-schema.md) for the
  reverse-engineered plist schema this relies on, including a gotcha this
  fix's own testing hit: two of this workflow's Script Filter nodes have
  no keyword set, and wiring the Universal Action to the wrong one (the
  action executor instead of the Check/Reveal/Clean chooser) silently
  produces "Nothing to do" for every selection.
- `internal/action`: every Error result row now sets `valid: false`, so
  Alfred no longer accepts Enter on it. Previously the row was silently
  actionable with no `action` variable set for the next step to dispatch
  on, so pressing Enter was a harmless no-op rather than a rejected
  keypress — safe, but not obviously so from the UI.

### Changed

- `README.md`/`README-jp.md` restructured to follow the [Alfred Gallery
  style guide](https://alfred.app/submit/styleguide/): `## Usage` now
  opens with "via the `cit` keyword" / "via the Universal Action" phrasing
  with a screenshot per entry point (`images/*.png`), modifier keys are
  documented as `<kbd>` tag bullets instead of a table, and `## Setup`
  covers only downloading and double-clicking the signed release.
  Build-from-source instructions moved to
  [DEVELOPING.md](DEVELOPING.md#building-the-workflow-package), since the
  style guide excludes installation instructions from `## Setup`.
- `docs/workflow-object-schema.md` moved to
  [docs/alfred-workflow-notes/](docs/alfred-workflow-notes/), a `git
  subtree` pulled from
  [y-marui/alfred-workflow-template](https://github.com/y-marui/alfred-workflow-template)
  (mirroring `docs/dev-charter/`'s existing pattern), since this content
  is generic Alfred-workflow knowledge shared across y-marui's other
  Alfred workflow projects, not specific to this one. Update it via `make
  update-workflow-notes`; direct edits under the prefix are blocked by
  pre-commit, same as `docs/dev-charter/`.

## [v1.0.0] - 2026-09-02

### Added

- Alfred Workflow specification: entry points, the Check/Reveal/Clean/Copy
  report actions, the Clean/Cleaned/Warning/Error state model, and
  accessibility/keyboard flow ([docs/specification.md](docs/specification.md)).
- [ADR 0001](docs/decisions/0001-separate-cli-and-workflow.md): keep the CLI
  and Alfred Workflow in separate repositories.
- [ADR 0002](docs/decisions/0002-file-based-cli-invocation.md): invoke
  `check`/`explain`/`fix` against a temporary file, since only `clean` reads
  standard input.
- `internal/clipboard`: reads/writes the macOS pasteboard's plain-text
  representation only, distinguishing "no text on the clipboard" from "empty
  text" without ever logging clipboard content.
- `internal/tempinput`: the private, owner-only-permission, single-use temp
  file `check`/`explain`/`fix` require as input, with guaranteed removal.
- `internal/action`, `internal/cliinvoke`, `internal/scriptfilter`, and
  `cmd/clean-invisible-text-alfred`: implements Check/Reveal/Clean/Copy
  report against the pinned CLI, printing Alfred Script Filter JSON. Copy
  report (excludes text by default, cmd-modifier includes it) and Clean's
  keep-warnings re-run (shift-modifier, Warning state only) are modifier-key
  alternate actions on the same result row, per
  docs/specification.md Accessibility and keyboard flow.
- `workflow/info.plist` and `scripts/build-workflow.sh`: wires the Keyword
  (`cit`, clipboard) path end-to-end and builds a universal (amd64+arm64)
  `.alfredworkflow` bundling the pinned CLI for both architectures. The
  Universal Action (selected text) needs a one-time manual step in Alfred's
  own UI — see README.md Setup — since that object isn't something this
  project can generate reproducibly from source; its downstream node
  already exists in `info.plist`, ready for that connection. Not yet
  verified inside Alfred's own Workflow debugger or on real Intel hardware.
- `internal/cliasset` and `scripts/fetch-cli-binaries.sh`: pins
  go-clean-invisible-text v1.1.1 with per-architecture SHA-256 checksums,
  downloads and verifies (checksum + build provenance attestation +
  Developer ID codesign) the darwin binaries into `assets/bin/` (gitignored,
  not fetched at ordinary runtime), and resolves/re-verifies the correct
  architecture's binary at runtime.
- [ADR 0003](docs/decisions/0003-v1-compatibility-and-upgrade-policy.md):
  minimum Alfred 5 / macOS 13 (Ventura), and confirms the embedded-CLI
  upgrade behavior (no runtime downloads; new CLI versions ship only via a
  new `.alfredworkflow` release).
- [docs/release-process.md](docs/release-process.md) and
  `.github/workflows/release.yml`: tag-triggered build and publish of the
  packaged `.alfredworkflow` (with checksums and build provenance
  attestation) as a GitHub Release.
- [docs/alfred-gallery-readiness.md](docs/alfred-gallery-readiness.md):
  checklist against the Alfred Gallery submission requirements and style
  guide, plus a one-time setup runbook for the GitHub Actions secrets below.
- `scripts/build-workflow.sh`, `scripts/notarize-binary.sh`, and
  `.github/workflows/release.yml`: sign and notarise the packaged entrypoint
  binary for a tagged release (Developer ID codesign with hardened runtime,
  App Store Connect API key for `notarytool`), gated behind five repository
  secrets so an ordinary `make build-workflow` stays unsigned. The embedded
  `go-clean-invisible-text` CLI binaries are now signed with a Developer ID
  upstream as of v1.1.1 (github.com/y-marui/go-clean-invisible-text#31);
  `scripts/fetch-cli-binaries.sh` now hard-fails if a fetched binary is
  missing a Developer ID signature, alongside its checksum/attestation
  checks.
