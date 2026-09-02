#!/usr/bin/env bash
# Block commits that directly edit files under an installed read-only
# git-subtree (e.g. docs/dev-charter/, docs/alfred-workflow-notes/). The
# only sanctioned way to change such a tree is `git subtree add`/`pull`/
# `merge --squash`.
#
# `git subtree` builds its "Squashed content" commit via `git commit-tree`,
# which never touches this hook — but the commit that actually joins that
# squashed history into the current branch (what `subtree add`/`pull`/
# `merge` do last) is a real merge commit made via the normal commit
# machinery, which DOES run pre-commit hooks. (This was learned the hard
# way: an earlier version of this comment claimed subtree operations
# bypass hooks "entirely" and got the working tree stuck mid-merge the
# first time this hook ever saw a real, non-no-op update — see
# alfred-clean-invisible-text#26.) Skip the check whenever a merge is in
# progress (`.git/MERGE_HEAD` exists) so that commit goes through; a
# hand-crafted `git add` + `git commit` under the prefix outside of a
# merge is still caught.
#
# Configure per pre-commit hook entry via env vars:
#   SUBTREE_PREFIX      - path to the subtree (default: docs/dev-charter)
#   SUBTREE_UPSTREAM    - repo name to point contributors at (default: dev-charter)
#   SUBTREE_UPDATE_HINT - command to suggest for pulling the update (default:
#                         "git subtree pull"); use the project's own wrapper
#                         (e.g. "make update-workflow-notes") when the raw
#                         git-subtree command isn't the right one to hand out
#                         (e.g. the shared content lives in a subdirectory of
#                         the upstream repo rather than at its root).
#
# Local-only safety net: this checks the staged diff (`git diff --cached`),
# which is what a real `git commit` sees. CI's `pre-commit run --all-files`
# runs against an already-committed working tree with nothing staged, so
# this hook is a no-op there — it does not catch a bad edit that already
# landed in a PR. Enforcing it in CI would need a separate check against
# the PR's base branch, not this script.
set -euo pipefail

PREFIX="${SUBTREE_PREFIX:-docs/dev-charter}"
UPSTREAM="${SUBTREE_UPSTREAM:-dev-charter}"
UPDATE_HINT="${SUBTREE_UPDATE_HINT:-git subtree pull}"

MERGE_HEAD_PATH=$(git rev-parse --git-path MERGE_HEAD 2>/dev/null || true)
[ -n "$MERGE_HEAD_PATH" ] && [ -f "$MERGE_HEAD_PATH" ] && exit 0

CHANGED=$(git diff --cached --name-only -- "$PREFIX" || true)
[ -n "$CHANGED" ] || exit 0

echo "error: ${PREFIX}/ 配下は直接編集禁止です。"
echo "  変更が必要な場合は ${UPSTREAM} リポジトリに Issue を立て、${UPDATE_HINT} で取り込んでください。"
exit 1
