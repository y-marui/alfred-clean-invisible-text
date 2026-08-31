# CLI Dependency Policy

The Workflow embeds released macOS binaries from go-clean-invisible-text.

- The CLI version is pinned for each Workflow release.
- SHA-256 checksums are verified during packaging.
- Intel and Apple Silicon binaries come from the same upstream CLI release.
- Runtime downloads are not performed during ordinary text processing.
- The packaged CLI version is included in diagnostics.
- Unicode rules are never copied into Workflow scripts.

## Implementation

- The pin (release tag and per-architecture SHA-256 checksums) lives in
  [internal/cliasset/pinned.txt](../internal/cliasset/pinned.txt), the single
  trust anchor for both packaging and runtime.
- `make fetch-cli` (`scripts/fetch-cli-binaries.sh`) downloads the pinned
  release's darwin binaries, verifies each against `pinned.txt` and its
  GitHub build provenance attestation, and stages them under `assets/bin/`
  (gitignored — never committed, never fetched at ordinary runtime).
- `internal/cliasset.ResolvePath` selects the correct architecture's staged
  binary at runtime and re-verifies its checksum before returning its path,
  rather than trusting the packaging-time check alone.
- `internal/cliasset.Version` exposes the pinned release tag for diagnostics.
- Updating the pin: re-run the verification in
  [scripts/fetch-cli-binaries.sh](../scripts/fetch-cli-binaries.sh) against
  the new release, then update `pinned.txt` only after it passes.
- [.github/workflows/check-cli-update.yml](../.github/workflows/check-cli-update.yml)
  runs [scripts/check-cli-update.sh](../scripts/check-cli-update.sh) weekly
  (and on manual dispatch) to check for a newer upstream release. Because a
  release that has never been pinned has no known-good checksum to compare
  against, this script verifies the new release's darwin binaries against
  their GitHub build provenance attestation instead, then opens a pull
  request updating `pinned.txt`. A human still reviews and merges that pull
  request — the workflow only proposes the update, it never merges it.
