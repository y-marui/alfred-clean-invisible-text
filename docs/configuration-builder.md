# Configuration Builder

Alfred's Workflow Configuration Sheet (`userconfigurationconfig` in
`workflow/info.plist`) lets a user set workflow variables from a form in
Alfred's UI, without editing the plist directly.

## Current state

None yet. This Workflow has no user-facing settings to expose there:

- The pinned CLI version and its per-architecture checksums
  (`internal/cliasset/pinned.txt`) are a maintainer-controlled trust anchor,
  not something a user should override — see
  [docs/dependency-policy.md](dependency-policy.md).
- The keyword and Universal Action triggers are configured through Alfred's
  own standard object UI (Preferences → Workflows), which doesn't need a
  Workflow Configuration Sheet entry.
- `ALFRED_CLEAN_ASSETS_DIR` (see `cmd/clean-invisible-text-alfred/main.go`)
  is a developer/test override, not a setting end users need.

If a future action genuinely needs a user-adjustable setting (e.g. a default
for `--keep-warnings`), add it here and to `workflow/info.plist`'s
`userconfigurationconfig` in the same change, per the Document Sync Rule in
[AI_CONTEXT.md](../AI_CONTEXT.md).
