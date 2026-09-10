#!/usr/bin/env bash
set -euo pipefail

expected_version="${1:-$(tr -d '[:space:]' < VERSION)}"

if [[ ! "$expected_version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "invalid release version: $expected_version" >&2
  exit 1
fi

version_file="$(tr -d '[:space:]' < VERSION)"
if [[ "$version_file" != "$expected_version" ]]; then
  echo "VERSION mismatch: expected $expected_version, found $version_file" >&2
  exit 1
fi

meta_version="$(grep -oE 'version[^\"]*\"[[:space:]]*:[[:space:]]*\"[0-9]+\.[0-9]+\.[0-9]+\"' internal/httpapi/api.go | head -n1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' || true)"
if [[ "$meta_version" != "$expected_version" ]]; then
  echo "API metadata mismatch: expected $expected_version, found ${meta_version:-<missing>}" >&2
  exit 1
fi

if ! grep -Fq "response.Version != \"$expected_version\"" internal/httpapi/meta_version_test.go; then
  echo "metadata regression test does not assert $expected_version" >&2
  exit 1
fi

echo "release consistency check passed for v$expected_version"
