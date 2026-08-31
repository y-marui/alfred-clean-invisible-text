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
| Binaries signed and notarised | ❌ Blocked | See [Code signing](#code-signing-and-notarization-the-main-blocker) below |
| No self-update | ✅ Done | [ADR 0003](decisions/0003-v1-compatibility-and-upgrade-policy.md) — updates only via new `.alfredworkflow` releases |
| No self-installed external software | ✅ Done | `make fetch-cli` runs at build/packaging time only; nothing is downloaded at ordinary runtime (`docs/dependency-policy.md`) |
| Icon ≥ 256×256 px | ✅ Done | `workflow/icon.png` is 512×512 |
| Keyword ≥ 3 characters | ✅ Done | `cit` (exactly 3; revisit if Gallery feedback asks for more headroom) |
| User Configuration over environment variables | ✅ N/A | Nothing today needs user-facing configuration; revisit if that changes |
| English instructions in About/README | ✅ Done | `README.md` is the canonical reference version ([LANGUAGE_POLICY.md](dev-charter/LANGUAGE_POLICY.md)) |
| README follows Gallery style guide | ⚠️ Partial | See [README style guide gaps](#readme-style-guide-gaps) below |
| Screenshots (full Alfred window, shadow, no background) | ❌ Not started | Needs the hands-on Alfred verification already tracked in [Issue #4](https://github.com/y-marui/alfred-clean-invisible-text/issues/4) |

## Code signing and notarization (the main blocker)

Gallery workflows containing compiled binaries must be signed with a Developer
ID and notarised by Apple — "or otherwise try to bypass this macOS security
feature" is explicitly disallowed. Two binaries ship inside this Workflow,
neither signed today:

1. **`cmd/clean-invisible-text-alfred`**, built by `scripts/build-workflow.sh`
   at packaging time. Signing this is entirely within this repository's
   control once a Developer ID is available.
2. **The pinned `go-clean-invisible-text` release binaries**
   (`assets/bin/clean-invisible-text-darwin-{amd64,arm64}`), built and
   published upstream. `codesign -dv` confirms these are unsigned as
   published. Signing these requires either upstream adding Developer ID
   signing to its own release workflow, or this repository re-signing
   (and notarising) the binaries it embeds during packaging — re-signing
   a third-party binary changes its trust story and should be a deliberate
   decision, not a default.

Both require an **Apple Developer Program membership** (currently $99/year,
tied to an Apple ID) to obtain a Developer ID Application certificate and to
notarise via `notarytool`. That is an account/cost decision only
@y-marui can make — this document stops at describing the requirement and
does not enroll, purchase, or configure signing credentials. Once enrolled,
the follow-up work is:

- Add codesigning + notarization steps to `scripts/build-workflow.sh` (or a
  dedicated release-only step) for the entrypoint binary.
- Decide the embedded-CLI-binary signing approach (upstream signs, or this
  repo re-signs) and record that decision as a new ADR — it affects
  `docs/dependency-policy.md`'s trust model (checksum + build provenance
  attestation today), not just packaging mechanics.
- Store the Developer ID certificate/notarization credentials as GitHub
  Actions secrets for `.github/workflows/release.yml` — never commit them.

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
- Enrolling in the Apple Developer Program and any signing/notarization
  implementation — blocked on that account decision above.
