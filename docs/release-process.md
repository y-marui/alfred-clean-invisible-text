# Release Process

`.github/workflows/release.yml` builds the packaged `.alfredworkflow` and
publishes it as a GitHub Release.

## Cutting a release

1. Bump `workflow/info.plist`'s `version` key to `X.Y.Z` (used by
   `scripts/build-workflow.sh` to name the output file, and checked against
   the tag below).
2. Update `CHANGELOG.md`: move `[Unreleased]` entries under a new
   `## [vX.Y.Z] - YYYY-MM-DD` heading.
3. Tag and push: `git tag vX.Y.Z && git push origin vX.Y.Z`.
4. The workflow verifies the tag matches `workflow/info.plist`'s version,
   runs `make fetch-cli && make build-workflow`, generates `checksums.txt`
   (SHA-256), attests build provenance for the `.alfredworkflow` via
   [`actions/attest-build-provenance`](https://github.com/actions/attest-build-provenance),
   and publishes a GitHub Release with both files attached plus
   auto-generated release notes.

`workflow_dispatch` runs the same build without publishing anything (the
version check, provenance attestation, and release-creation steps only run
on an actual `refs/tags/*` push) — use it to validate the build after
changing the workflow, before ever pushing a version tag.

## Cutting a release without Actions

If Actions can't run (e.g. a billing/spending-limit issue), `scripts/release.sh`
(wired as `make release`) reproduces the same build and checksum steps
locally and publishes the same GitHub Release, prompting for confirmation
before pushing the tag or creating the release. It does not attest build
provenance — that needs the Actions runner's OIDC token — so a
locally-published release's artifact won't verify with `gh attestation verify`.

## Verifying a downloaded `.alfredworkflow`

```bash
shasum -a 256 -c checksums.txt
gh attestation verify clean-invisible-text-alfred-<version>.alfredworkflow --repo y-marui/alfred-clean-invisible-text
```

## Code signing and notarization

The packaged entrypoint binary is signed (Developer ID, hardened runtime)
and notarised for a tagged release, gated behind five GitHub Actions
secrets; the embedded `go-clean-invisible-text` CLI binaries are signed
upstream as of the pinned v1.1.1. See
[docs/alfred-gallery-readiness.md](alfred-gallery-readiness.md) for details
and the one-time secret setup.

## Related

- [docs/dependency-policy.md](dependency-policy.md) — how the embedded CLI
  binary is pinned and verified (a separate, ongoing process from cutting a
  Workflow release)
- [docs/decisions/0003-v1-compatibility-and-upgrade-policy.md](decisions/0003-v1-compatibility-and-upgrade-policy.md) —
  minimum Alfred/macOS versions and the CLI upgrade policy this release
  process implements
