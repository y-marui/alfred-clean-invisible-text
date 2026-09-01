# Alfred Gallery Readiness

Tracks this Workflow's compliance with the
[Alfred Gallery submission requirements](https://alfred.app/submit/) and
[style guide](https://alfred.app/submit/styleguide/), per
[Issue #6](https://github.com/y-marui/alfred-clean-invisible-text/issues/6).
This is a checklist against an external, occasionally-changing policy —
re-read the linked pages before acting on stale entries here.

## Submission process

Alfred does not accept Gallery submissions directly. The documented path is:
share the workflow on the [Alfred Forum](https://www.alfredforum.com/) first;
once it is "generally stable and trusted by a number of users," the Alfred
team may invite an official Gallery submission. There is no self-service
form. Per Issue #6, forum posting stays out of scope here unless explicitly
requested — this document only tracks technical/documentation readiness so
that submission is not blocked on our side whenever that step happens.

## Checklist

| Requirement | Status | Notes |
|---|---|---|
| Binaries signed and notarised | ⚠️ Partial | Entrypoint binary: automated in CI, pending secret setup. Embedded CLI binaries: still unsigned upstream. See [Code signing](#code-signing-and-notarization) below |
| No self-update | ✅ Done | [ADR 0003](decisions/0003-v1-compatibility-and-upgrade-policy.md) — updates only via new `.alfredworkflow` releases |
| No self-installed external software | ✅ Done | `make fetch-cli` runs at build/packaging time only; nothing is downloaded at ordinary runtime (`docs/dependency-policy.md`) |
| Icon ≥ 256×256 px | ✅ Done | `workflow/icon.png` is 512×512 |
| Keyword ≥ 3 characters | ✅ Done | `cit` (exactly 3; revisit if Gallery feedback asks for more headroom) |
| User Configuration over environment variables | ✅ N/A | Nothing today needs user-facing configuration; revisit if that changes |
| English instructions in About/README | ✅ Done | `README.md` is the canonical reference version ([LANGUAGE_POLICY.md](dev-charter/LANGUAGE_POLICY.md)) |
| README follows Gallery style guide | ⚠️ Partial | See [README style guide gaps](#readme-style-guide-gaps) below |
| Screenshots (full Alfred window, shadow, no background) | ❌ Not started | Needs the hands-on Alfred verification already tracked in [Issue #4](https://github.com/y-marui/alfred-clean-invisible-text/issues/4) |

## Code signing and notarization

Gallery workflows containing compiled binaries must be signed with a Developer
ID and notarised by Apple — "or otherwise try to bypass this macOS security
feature" is explicitly disallowed. Two binaries ship inside this Workflow:

1. **`cmd/clean-invisible-text-alfred`**, built by `scripts/build-workflow.sh`
   at packaging time. `.github/workflows/release.yml` now signs and
   notarises it automatically for a tagged release, gated behind the GitHub
   Actions secrets described below — until those secrets exist, the step
   fails loudly rather than silently publishing an unsigned binary.
2. **The pinned `go-clean-invisible-text` release binaries**
   (`assets/bin/clean-invisible-text-darwin-{amd64,arm64}`), built and
   published upstream, still unsigned as of the currently pinned version
   (`internal/cliasset/pinned.txt`). `codesign -dv` confirms this;
   `scripts/fetch-cli-binaries.sh` now reports the signing status of every
   fetched binary (informational only — it does not yet fail the build,
   since upstream hasn't shipped a signed release to require). The chosen
   direction is for **`go-clean-invisible-text` to sign its own darwin
   release binaries** (matches [ADR 0001](decisions/0001-separate-cli-and-workflow.md)'s
   separate-lifecycle split, and keeps this repo's checksum/attestation
   trust model in `docs/dependency-policy.md` intact — re-signing a
   third-party binary here would change the bytes `pinned.txt` checksums).
   Tracking that as upstream work is still open; once it ships,
   `scripts/fetch-cli-binaries.sh`'s codesign check should become a hard
   failure alongside its checksum/attestation checks.

### Apple Developer Program status

@y-marui is enrolled, and a Developer ID Application certificate already
exists in the local Keychain (Team ID `7TEQWKRRX7`, confirmed via
`security find-identity -v -p codesigning`). What's left is exporting that
certificate and generating notarization credentials, then registering both
as GitHub Actions secrets — steps that need to happen on @y-marui's machine
and Apple/GitHub accounts, not something this repository's automation can do
on its own.

### One-time setup (manual, @y-marui)

`.github/workflows/release.yml` expects five repository secrets. None of
them should ever be committed or pasted into an issue/PR/chat — set them
directly via `gh secret set` or the GitHub UI.

**1. Export the Developer ID Application certificate as `.p12`:**

Keychain Access → login keychain → My Certificates → find
"Developer ID Application: Yukihiro Marui (7TEQWKRRX7)" → expand it, select
both the certificate and its private key → right-click → Export 2 items… →
save as `certificate.p12`, choosing an export password (this becomes
`MACOS_CERTIFICATE_PASSWORD` below).

```bash
base64 -i certificate.p12 | gh secret set MACOS_CERTIFICATE_P12 \
  --repo y-marui/alfred-clean-invisible-text
gh secret set MACOS_CERTIFICATE_PASSWORD \
  --repo y-marui/alfred-clean-invisible-text  # paste the export password when prompted
rm certificate.p12  # don't leave the exported key sitting on disk
```

**2. Generate an App Store Connect API key for notarization:**

[appstoreconnect.apple.com](https://appstoreconnect.apple.com/) → Users and
Access → Integrations → App Store Connect API → generate a key with the
"Developer" role. Apple lets you download the `.p8` private key file only
once — save it somewhere safe until it's registered below, then delete it.
Note the Key ID and Issuer ID shown on the same page.

```bash
base64 -i AuthKey_XXXXXXXXXX.p8 | gh secret set NOTARY_API_KEY \
  --repo y-marui/alfred-clean-invisible-text
gh secret set NOTARY_KEY_ID --repo y-marui/alfred-clean-invisible-text     # the Key ID
gh secret set NOTARY_ISSUER_ID --repo y-marui/alfred-clean-invisible-text  # the Issuer ID
rm AuthKey_XXXXXXXXXX.p8
```

Once all five secrets (`MACOS_CERTIFICATE_P12`, `MACOS_CERTIFICATE_PASSWORD`,
`NOTARY_API_KEY`, `NOTARY_KEY_ID`, `NOTARY_ISSUER_ID`) exist, the next tagged
release automatically signs and notarises the entrypoint binary — no further
workflow changes needed. `workflow_dispatch` runs never touch these secrets
(see the `if: startsWith(github.ref, 'refs/tags/')` guards in
`release.yml`), so iterating on the workflow doesn't burn notarization
requests.

## README style guide gaps

The Gallery style guide expects specific phrasing patterns the current
`README.md`/`README-jp.md` don't yet fully follow:

- A `## Usage` section whose opening lines follow the "via the `keyword`" /
  "via the Universal Action" phrasing convention.
- Modifier-key documentation immediately below a screenshot, using `<kbd>`
  tags (`* <kbd>⌘</kbd><kbd>Y</kbd> ...`) rather than the current table.
- A `## Setup` section (already present) limited to manual steps only —
  Gallery installations skip the `make fetch-cli && make build-workflow`
  instructions entirely, so that section will need Gallery-specific wording
  once distribution moves off "build from source."

Restructuring the README to this format is worth doing close to actual
submission (it depends on final screenshots), not speculatively now.

## Out of scope here

- Posting to the Alfred Forum — optional per Issue #6, do only on explicit
  request.
- Screenshots — need a real Alfred window on the user's machine; tracked
  under Issue #4's hands-on verification, not here.
- Exporting the certificate, generating the API key, and registering GitHub
  secrets — manual steps only @y-marui can perform (see
  [One-time setup](#one-time-setup-manual-y-marui) above); this repository's
  automation only consumes those secrets once they exist.
- Adding Developer ID signing to `go-clean-invisible-text`'s own release
  workflow — separate repository, separate lifecycle (ADR 0001); not
  implemented from here.
