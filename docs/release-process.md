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

## Verifying a downloaded `.alfredworkflow`

```bash
shasum -a 256 -c checksums.txt
gh attestation verify clean-invisible-text-alfred-<version>.alfredworkflow --repo y-marui/alfred-clean-invisible-text
```

## What this does not cover yet

Neither the packaged entrypoint binary nor the embedded
`go-clean-invisible-text` CLI binaries are code-signed or notarised — see
[docs/alfred-gallery-readiness.md](alfred-gallery-readiness.md) for what
that requires and why it's blocked on an Apple Developer Program decision.

## Related

- [docs/dependency-policy.md](dependency-policy.md) — how the embedded CLI
  binary is pinned and verified (a separate, ongoing process from cutting a
  Workflow release)
- [docs/decisions/0003-v1-compatibility-and-upgrade-policy.md](decisions/0003-v1-compatibility-and-upgrade-policy.md) —
  minimum Alfred/macOS versions and the CLI upgrade policy this release
  process implements
