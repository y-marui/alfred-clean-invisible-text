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

Add this directory as a subtree, pulling from this repository:

```bash
git subtree add --prefix=docs/alfred-workflow-notes \
  https://github.com/y-marui/alfred-workflow-template main --squash
```

To pull later updates, add a Makefile target mirroring this repo's
`update-charter` (see [`Makefile`](../../Makefile)):

```makefile
update-workflow-notes:
	git remote | grep -q '^alfred-workflow-notes$$' || \
	  git remote add alfred-workflow-notes https://github.com/y-marui/alfred-workflow-template
	git fetch alfred-workflow-notes
	@STASHED=0; \
	if ! git diff --quiet || ! git diff --cached --quiet || [ -n "$$(git ls-files --others --exclude-standard)" ]; then \
		git stash push -u -m "update-workflow-notes"; \
		STASHED=1; \
	fi; \
	git subtree pull --prefix=docs/alfred-workflow-notes alfred-workflow-notes main --squash; \
	if [ "$$STASHED" = "1" ]; then git stash pop; fi
```

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
