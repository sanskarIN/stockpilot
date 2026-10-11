#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

EXPECTED_VERSION="$(tr -d '[:space:]' < VERSION)"
VERSION_FILE="VERSION"
CHANGELOG_FILE="CHANGELOG.md"
RELEASE_DOC="docs/releases/v${EXPECTED_VERSION}.md"

fail() {
  echo "release-readiness: $1" >&2
  exit 1
}

[[ -f "$VERSION_FILE" ]] || fail "missing $VERSION_FILE"
[[ -f "$CHANGELOG_FILE" ]] || fail "missing $CHANGELOG_FILE"
[[ -f "$RELEASE_DOC" ]] || fail "missing $RELEASE_DOC"

actual_version="$(tr -d '[:space:]' < "$VERSION_FILE")"
[[ "$actual_version" == "$EXPECTED_VERSION" ]] || fail "VERSION is '$actual_version'; expected '$EXPECTED_VERSION'"

printf '%s\n' "$actual_version" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$' || fail "VERSION is not semantic-version shaped"

grep -Fq "### v${EXPECTED_VERSION} —" "$CHANGELOG_FILE" || fail "CHANGELOG.md has no v${EXPECTED_VERSION} release entry"
grep -Fq "# StockPilot v${EXPECTED_VERSION} —" "$RELEASE_DOC" || fail "release documentation has no v${EXPECTED_VERSION} heading"
grep -Fq 'scripts/check-release-consistency.sh' "$RELEASE_DOC" || fail "release documentation does not reference release consistency validation"
grep -Fq 'make release-check' "$RELEASE_DOC" || fail "release documentation does not reference make release-check"
grep -Fq 'go mod verify' "$RELEASE_DOC" || fail "release documentation does not reference module verification"

if command -v go >/dev/null 2>&1; then
  go mod verify
else
  echo "release-readiness: Go is not installed; skipping local module verification"
fi

echo "release-readiness: v${EXPECTED_VERSION} repository checks passed"
