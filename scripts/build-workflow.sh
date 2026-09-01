#!/usr/bin/env bash
# Build the .alfredworkflow package.
#
# Steps:
#   1. Build cmd/clean-invisible-text-alfred as a universal (amd64+arm64)
#      binary via lipo, since this Workflow — unlike the CLI it wraps — must
#      run natively on both Intel and Apple Silicon from a single bundle.
#   2. If CODESIGN_IDENTITY is set, codesign that binary (and notarize it if
#      NOTARY_KEY_ID is also set). Unset by default, so an ordinary local
#      `make build-workflow` stays unsigned; .github/workflows/release.yml
#      sets these for tagged releases. See docs/alfred-gallery-readiness.md.
#   3. Copy workflow/ (info.plist, icon.png) into the build dir.
#   4. Copy the pinned, checksum-verified CLI binaries from assets/bin/
#      (run `make fetch-cli` first if that's empty).
#   5. Zip into dist/<name>-<version>.alfredworkflow.
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
WORKFLOW_DIR="${REPO_ROOT}/workflow"
DIST_DIR="${REPO_ROOT}/dist"
BUILD_DIR="${REPO_ROOT}/.build"
ASSETS_SRC="${REPO_ROOT}/assets/bin"

if [ ! -d "$ASSETS_SRC" ] || [ -z "$(ls -A "$ASSETS_SRC" 2>/dev/null)" ]; then
  echo "error: ${ASSETS_SRC} is empty. Run 'make fetch-cli' first." >&2
  exit 1
fi

echo "→ Preparing build directory"
rm -rf "$BUILD_DIR"
cp -r "$WORKFLOW_DIR/" "$BUILD_DIR/"

echo "→ Building universal binary (amd64 + arm64)"
GOOS=darwin GOARCH=amd64 go build -o "${BUILD_DIR}/clean-invisible-text-alfred-amd64" ./cmd/clean-invisible-text-alfred
GOOS=darwin GOARCH=arm64 go build -o "${BUILD_DIR}/clean-invisible-text-alfred-arm64" ./cmd/clean-invisible-text-alfred
lipo -create -output "${BUILD_DIR}/clean-invisible-text-alfred" \
  "${BUILD_DIR}/clean-invisible-text-alfred-amd64" \
  "${BUILD_DIR}/clean-invisible-text-alfred-arm64"
rm "${BUILD_DIR}/clean-invisible-text-alfred-amd64" "${BUILD_DIR}/clean-invisible-text-alfred-arm64"
chmod +x "${BUILD_DIR}/clean-invisible-text-alfred"
lipo -info "${BUILD_DIR}/clean-invisible-text-alfred"

if [ -n "${CODESIGN_IDENTITY:-}" ]; then
  echo "→ Signing entrypoint binary (${CODESIGN_IDENTITY})"
  codesign --force --options runtime --timestamp --sign "$CODESIGN_IDENTITY" \
    "${BUILD_DIR}/clean-invisible-text-alfred"
  codesign --verify --strict --verbose=2 "${BUILD_DIR}/clean-invisible-text-alfred"

  if [ -n "${NOTARY_KEY_ID:-}" ]; then
    "${REPO_ROOT}/scripts/notarize-binary.sh" "${BUILD_DIR}/clean-invisible-text-alfred"
  fi
else
  echo "→ Skipping signing (CODESIGN_IDENTITY not set) — unsigned local/dev build"
fi

echo "→ Staging pinned CLI binaries"
mkdir -p "${BUILD_DIR}/assets/bin"
cp "${ASSETS_SRC}"/* "${BUILD_DIR}/assets/bin/"
chmod +x "${BUILD_DIR}"/assets/bin/*

VERSION=$(/usr/libexec/PlistBuddy -c "Print :version" "${BUILD_DIR}/info.plist")
WORKFLOW_NAME=$(/usr/libexec/PlistBuddy -c "Print :name" "${BUILD_DIR}/info.plist" | tr '[:upper:] ' '[:lower:]-')

mkdir -p "$DIST_DIR"
OUTPUT="${DIST_DIR}/${WORKFLOW_NAME}-${VERSION}.alfredworkflow"
rm -f "$OUTPUT" # ensure a clean zip (zip -r updates rather than replaces)

echo "→ Packaging: ${OUTPUT}"
(cd "$BUILD_DIR" && zip -r "$OUTPUT" . -x "*.DS_Store" --quiet)

echo "✓ Build complete: ${OUTPUT}"
