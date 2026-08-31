#!/usr/bin/env bash
# Check whether go-clean-invisible-text has a newer release than the one
# pinned in internal/cliasset/pinned.txt. If so, verify the new release's
# darwin binaries via GitHub build provenance attestation (there is no
# known-good checksum to compare against yet for a release that has never
# been pinned — attestation is the trust anchor for this bootstrap step)
# and rewrite pinned.txt with the new version and checksums.
#
# Intended to run only from check-cli-update.yml (CI), which opens a pull
# request from the result for a human to review and merge. See
# docs/dependency-policy.md.
#
# Requires: gh (authenticated), shasum or sha256sum.
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
UPSTREAM_REPO="y-marui/go-clean-invisible-text"
MANIFEST="${REPO_ROOT}/internal/cliasset/pinned.txt"

CURRENT_VERSION=$(grep -E '^version=' "$MANIFEST" | cut -d= -f2-)
if [ -z "$CURRENT_VERSION" ]; then
  echo "error: could not read a pinned version from ${MANIFEST}" >&2
  exit 1
fi

LATEST_VERSION=$(gh release view --repo "$UPSTREAM_REPO" --json tagName -q .tagName)
if [ -z "$LATEST_VERSION" ]; then
  echo "error: could not determine the latest ${UPSTREAM_REPO} release" >&2
  exit 1
fi

echo "pinned: ${CURRENT_VERSION}, latest: ${LATEST_VERSION}"

if [ "$LATEST_VERSION" = "$CURRENT_VERSION" ]; then
  echo "up to date"
  if [ -n "${GITHUB_OUTPUT:-}" ]; then
    echo "updated=false" >>"$GITHUB_OUTPUT"
  fi
  exit 0
fi

sha256_of() {
  if command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | cut -d' ' -f1
  else
    sha256sum "$1" | cut -d' ' -f1
  fi
}

WORKDIR=$(mktemp -d)
trap 'rm -rf "$WORKDIR"' EXIT

amd64_sum=""
arm64_sum=""
for arch_key in darwin-amd64 darwin-arm64; do
  asset="clean-invisible-text-${arch_key}"
  echo "Fetching ${asset} @ ${LATEST_VERSION}..."
  gh release download "$LATEST_VERSION" --repo "$UPSTREAM_REPO" \
    --pattern "$asset" --dir "$WORKDIR" --clobber

  echo "  verifying build provenance attestation..."
  if ! gh attestation verify "${WORKDIR}/${asset}" --repo "$UPSTREAM_REPO" >/dev/null; then
    echo "error: build provenance attestation failed for ${asset} @ ${LATEST_VERSION}" >&2
    echo "  refusing to pin an unverified release; pinned.txt was not changed." >&2
    exit 1
  fi

  sum=$(sha256_of "${WORKDIR}/${asset}")
  echo "  attestation OK, sha256 ${sum}"
  if [ "$arch_key" = "darwin-amd64" ]; then
    amd64_sum="$sum"
  else
    arm64_sum="$sum"
  fi
done

cat >"$MANIFEST" <<EOF
# Pinned go-clean-invisible-text release for this Workflow.
# Update via \`make fetch-cli\` after manual verification, or review and
# merge the pull request opened by check-cli-update.yml — see
# docs/dependency-policy.md.
version=${LATEST_VERSION}
darwin-amd64=${amd64_sum}
darwin-arm64=${arm64_sum}
EOF

echo "pinned.txt updated to ${LATEST_VERSION}"
if [ -n "${GITHUB_OUTPUT:-}" ]; then
  echo "updated=true" >>"$GITHUB_OUTPUT"
  echo "version=${LATEST_VERSION}" >>"$GITHUB_OUTPUT"
fi
