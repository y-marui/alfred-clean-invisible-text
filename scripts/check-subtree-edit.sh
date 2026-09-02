#!/usr/bin/env bash
# Block commits that directly edit files under an installed read-only
# git-subtree (e.g. docs/dev-charter/, docs/alfred-workflow-notes/). The
# only sanctioned way to change such a tree is `git subtree add`/`pull
# --squash`, and those bypass the normal commit hooks entirely
# (git-subtree builds the squash commit with `git commit-tree` rather
# than `git commit`), so this hook never sees a legitimate subtree
# update — anything it does see staged under the prefix is a direct edit
# and gets rejected.
#
# Configure per pre-commit hook entry via env vars:
#   SUBTREE_PREFIX     - path to the subtree (default: docs/dev-charter)
#   SUBTREE_UPSTREAM    - "owner/repo" to point contributors at (default: dev-charter)
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

CHANGED=$(git diff --cached --name-only -- "$PREFIX" || true)
[ -n "$CHANGED" ] || exit 0

echo "error: ${PREFIX}/ 配下は直接編集禁止です。"
echo "  変更が必要な場合は ${UPSTREAM} リポジトリに Issue を立て、git subtree pull で取り込んでください。"
exit 1
