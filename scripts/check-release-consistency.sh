#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
version_file="$root_dir/VERSION"
api_file="$root_dir/internal/httpapi/api.go"

if [[ ! -f "$version_file" ]]; then
  echo "ERROR: VERSION file not found: $version_file" >&2
  exit 1
fi

if [[ ! -f "$api_file" ]]; then
  echo "ERROR: API source file not found: $api_file" >&2
  exit 1
fi

version="$(tr -d '[:space:]' < "$version_file")"
if [[ ! "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "ERROR: VERSION must contain a semantic version such as 0.5.5; found '$version'" >&2
  exit 1
fi

meta_version="$({
  grep -oE '"version"[[:space:]]*:[[:space:]]*"[0-9]+\.[0-9]+\.[0-9]+"' "$api_file" || true
} | sed -E 's/.*"([0-9]+\.[0-9]+\.[0-9]+)"/\1/' | head -n 1)"

if [[ -z "$meta_version" ]]; then
  echo "ERROR: could not locate the /api/v1/meta version literal in $api_file" >&2
  exit 1
fi

if [[ "$meta_version" != "$version" ]]; then
  echo "ERROR: VERSION ($version) does not match /api/v1/meta ($meta_version)" >&2
  exit 1
fi

echo "Release consistency check passed: VERSION=$version, /api/v1/meta=$meta_version"
