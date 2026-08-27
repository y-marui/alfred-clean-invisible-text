#!/usr/bin/env bash
# Download the pinned go-clean-invisible-text release's darwin binaries,
# verify each against internal/cliasset/pinned.txt (the trust anchor checked
# into this repository) and, best-effort, its build provenance attestation,
# then stage them under assets/bin/ (gitignored — see docs/dependency-policy.md
# "Runtime downloads are not performed during ordinary text processing").
#
# Requires: gh (authenticated), shasum or sha256sum.
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
UPSTREAM_REPO="y-marui/go-clean-invisible-text"
MANIFEST="${REPO_ROOT}/internal/cliasset/pinned.txt"
OUT_DIR="${REPO_ROOT}/assets/bin"

VERSION=$(grep -E '^version=' "$MANIFEST" | cut -d= -f2-)
if [ -z "$VERSION" ]; then
  echo "error: could not read a pinned version from ${MANIFEST}" >&2
  exit 1
fi

sha256_of() {
  if command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | cut -d' ' -f1
  else
    sha256sum "$1" | cut -d' ' -f1
  fi
}

mkdir -p "$OUT_DIR"
WORKDIR=$(mktemp -d)
trap 'rm -rf "$WORKDIR"' EXIT

status=0
for arch_key in darwin-amd64 darwin-arm64; do
  pinned_sum=$(grep -E "^${arch_key}=" "$MANIFEST" | cut -d= -f2-)
  if [ -z "$pinned_sum" ]; then
    echo "error: no pinned checksum for ${arch_key} in ${MANIFEST}" >&2
    status=1
    continue
  fi

  asset="clean-invisible-text-${arch_key}"
  echo "Fetching ${asset} @ ${VERSION}..."
  gh release download "$VERSION" --repo "$UPSTREAM_REPO" \
    --pattern "$asset" --dir "$WORKDIR" --clobber

  actual_sum=$(sha256_of "${WORKDIR}/${asset}")
  if [ "$actual_sum" != "$pinned_sum" ]; then
    echo "error: checksum mismatch for ${asset}: got ${actual_sum}, pinned ${pinned_sum}" >&2
    echo "  Do not use this binary. If ${VERSION} was intentionally re-cut, update ${MANIFEST} only after independently re-verifying the new release." >&2
    status=1
    continue
  fi
  echo "  checksum OK (${actual_sum})"

  if command -v gh >/dev/null 2>&1 && gh auth status >/dev/null 2>&1; then
    if gh attestation verify "${WORKDIR}/${asset}" --repo "$UPSTREAM_REPO" >/dev/null 2>&1; then
      echo "  attestation OK"
    else
      echo "error: build provenance attestation failed for ${asset}" >&2
      status=1
      continue
    fi
  else
    echo "  warning: gh not authenticated; skipping attestation verification" >&2
  fi

  install -m 0755 "${WORKDIR}/${asset}" "${OUT_DIR}/${asset}"
done

if [ "$status" -ne 0 ]; then
  echo "error: one or more binaries failed verification; assets/bin/ was not fully updated" >&2
  exit 1
fi

echo "All pinned binaries verified and staged in ${OUT_DIR}."
