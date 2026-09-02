# Alfred Workflow Notes

Cross-project reference material for building Alfred 5 Script Filter
workflows — knowledge that isn't specific to any single workflow project.

`alfred-workflow-template` is the **canonical** source for this directory.
Edit files under here directly in that repository. Other Alfred workflow
projects (including ones generated from this template) pull it in via
`git subtree`, the same distribution pattern this repo itself uses for
[`docs/dev-charter/`](../dev-charter/README.md).

## Contents

| Document | Description |
|---|---|
| [configuration-builder.md](configuration-builder.md) | Alfred's Configuration Builder (`userconfigurationconfig`) — widget types and their config keys |
| [workflow-object-schema.md](workflow-object-schema.md) | Reverse-engineered schema of `info.plist` objects (Script Filter, Universal Action Trigger, etc.) — Alfred does not publish this |

## Installing in a consuming project

This directory lives at `docs/alfred-workflow-notes/` *inside* the
`alfred-workflow-template` repo, not at that repo's root (unlike
`docs/dev-charter/`, whose upstream repo root *is* the shared content). A
plain `git subtree add`/`pull` against `alfred-workflow-template` pulls in
that repo's entire root — `.github/`, `pyproject.toml`, `uv.lock`, everything
— not just this subdirectory. Split this subdirectory's history out first,
then add/merge just that.

`git subtree split` (without `--branch`) prints the split commit's SHA
without creating any local branch — no throwaway ref that could collide
with a branch the consuming repo already has, and nothing to clean up
afterwards.

First-time install:

```bash
set -e
git remote add alfred-workflow-notes https://github.com/y-marui/alfred-workflow-template
git fetch alfred-workflow-notes
SPLIT_SHA=$(git subtree split --prefix=docs/alfred-workflow-notes alfred-workflow-notes/main)
git subtree add --prefix=docs/alfred-workflow-notes "$SPLIT_SHA" --squash
```

To pull later updates, add a Makefile target mirroring this repo's
`update-charter` (see [`Makefile`](../../Makefile)) — still a single
`make update-workflow-notes` one-liner for the caller, the split just
happens inside the target:

```makefile
update-workflow-notes:
	git remote | grep -q '^alfred-workflow-notes$$' || \
	  git remote add alfred-workflow-notes https://github.com/y-marui/alfred-workflow-template
	git fetch alfred-workflow-notes
	@set -e; \
	STASHED=0; \
	if ! git diff --quiet || ! git diff --cached --quiet || [ -n "$$(git ls-files --others --exclude-standard)" ]; then \
		git stash push -u -m "update-workflow-notes"; \
		STASHED=1; \
	fi; \
	SPLIT_SHA=$$(git subtree split --prefix=docs/alfred-workflow-notes alfred-workflow-notes/main); \
	git subtree merge --prefix=docs/alfred-workflow-notes "$$SPLIT_SHA" --squash; \
	if [ "$$STASHED" = "1" ]; then git stash pop; fi
```

`set -e` makes a failed `git subtree split` (or any other step) abort the
recipe immediately instead of silently falling through to `merge` with a
stale or empty `SPLIT_SHA`. If the recipe aborts after stashing, the stash
is left in place rather than popped — recoverable manually (`git stash
list`), and safer than popping onto a tree a failed step may have left in
an unexpected state.

And a pre-commit hook blocking direct edits under the installed subtree,
mirroring this repo's `scripts/check-charter-subtree-edit.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail
PREFIX="docs/alfred-workflow-notes"
CHANGED=$(git diff --cached --name-only -- "$PREFIX" || true)
[ -n "$CHANGED" ] || exit 0
echo "error: ${PREFIX}/ must not be edited directly."
echo "  Open an issue against alfred-workflow-template and pull the update via git subtree pull."
exit 1
```

No CI automation pulls updates on a schedule — same choice `dev-charter`
makes. Pull manually when you know something changed.

If a change is needed, open an issue against `alfred-workflow-template`
(this repo) rather than editing the consuming project's copy — edits made
downstream are lost on the next `git subtree pull` and never reach the
other projects sharing this directory.
