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
| Binaries signed and notarised | ✅ Done | Both binaries verified signed/notarised in the actual [v1.0.0 release](https://github.com/y-marui/alfred-clean-invisible-text/releases/tag/v1.0.0) (`codesign -dvvv`, `gh attestation verify`). See [Code signing](#code-signing-and-notarization) below |
| No self-update | ✅ Done | [ADR 0003](decisions/0003-v1-compatibility-and-upgrade-policy.md) — updates only via new `.alfredworkflow` releases |
| No self-installed external software | ✅ Done | `make fetch-cli` runs at build/packaging time only; nothing is downloaded at ordinary runtime (`docs/dependency-policy.md`) |
| Icon ≥ 256×256 px | ✅ Done | `workflow/icon.png` is 512×512 |
| Keyword ≥ 3 characters | ✅ Done | `cit` (exactly 3; revisit if Gallery feedback asks for more headroom) |
| User Configuration over environment variables | ✅ N/A | Nothing today needs user-facing configuration; revisit if that changes |
| English instructions in About/README | ✅ Done | `README.md` is the canonical reference version ([LANGUAGE_POLICY.md](dev-charter/LANGUAGE_POLICY.md)) |
| README follows Gallery style guide | ✅ Done | See [README style guide gaps](#readme-style-guide-gaps) below |
| Screenshots (full Alfred window, shadow, no background) | ✅ Done | `images/*.png`, real window captures (rounded corners + drop shadow verified pixel-by-pixel, transparent background) |

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
   published upstream. **Signed as of upstream v1.1.1**
   ([go-clean-invisible-text#31](https://github.com/y-marui/go-clean-invisible-text/issues/31),
   resolved via
   [go-clean-invisible-text#32](https://github.com/y-marui/go-clean-invisible-text/pull/32)) —
   this matches [ADR 0001](decisions/0001-separate-cli-and-workflow.md)'s
   separate-lifecycle split and keeps this repo's checksum/attestation trust
   model in `docs/dependency-policy.md` intact (re-signing a third-party
   binary here would change the bytes `pinned.txt` checksums against).
   `internal/cliasset/pinned.txt` is now pinned to v1.1.1, and
   `scripts/fetch-cli-binaries.sh` hard-fails (alongside its
   checksum/attestation checks) if a fetched binary isn't signed by a
   Developer ID authority.

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

Done. `README.md`/`README-jp.md` now follow the patterns from
[alfred.app/submit/styleguide/](https://alfred.app/submit/styleguide/):

- `## Usage` opens with "via the `cit` keyword" / "via the Universal
  Action" phrasing, one screenshot per entry point.
- Modifier-key documentation is a bullet list of `<kbd>` tags
  (`* <kbd>⌘</kbd><kbd>↩︎</kbd> ...`) immediately below the relevant
  screenshot, not a table.
- `## Setup` only covers downloading and double-clicking the signed
  release — no manual configuration exists to document. Build-from-source
  instructions (`make fetch-cli && make build-workflow`) moved to
  [DEVELOPING.md](../DEVELOPING.md), a contributor-facing doc, per the
  style guide's rule against installation instructions in `## Setup`.

## Out of scope here

- Posting to the Alfred Forum — optional per Issue #6, do only on explicit
  request.
- Exporting the certificate, generating the API key, and registering GitHub
  secrets — manual steps only @y-marui can perform (see
  [One-time setup](#one-time-setup-manual-y-marui) above); this repository's
  automation only consumes those secrets once they exist.
- Adding Developer ID signing to `go-clean-invisible-text`'s own release
  workflow — separate repository, separate lifecycle (ADR 0001); not
  implemented from here.
