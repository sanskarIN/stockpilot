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

web_version="$(node -e 'console.log(JSON.parse(require("fs").readFileSync("web/package.json", "utf8")).version)' 2>/dev/null || true)"
if [[ "$web_version" != "$expected_version" ]]; then
  echo "web package version mismatch: expected $expected_version, found ${web_version:-<missing>}" >&2
  exit 1
fi

extension_package_version="$(node -e 'console.log(JSON.parse(require("fs").readFileSync("extension/package.json", "utf8")).version)' 2>/dev/null || true)"
if [[ "$extension_package_version" != "$expected_version" ]]; then
  echo "extension package version mismatch: expected $expected_version, found ${extension_package_version:-<missing>}" >&2
  exit 1
fi

extension_manifest_version="$(node -e 'console.log(JSON.parse(require("fs").readFileSync("extension/manifest.json", "utf8")).version)' 2>/dev/null || true)"
if [[ "$extension_manifest_version" != "$expected_version" ]]; then
  echo "extension manifest version mismatch: expected $expected_version, found ${extension_manifest_version:-<missing>}" >&2
  exit 1
fi

IFS=. read -r major minor patch <<< "$expected_version"
expected_android_code=$((10#$major * 1000000 + 10#$minor * 1000 + 10#$patch))
android_version_name="$(sed -n 's/^[[:space:]]*versionName = "\([0-9][0-9.]*\)"/\1/p' android/app/build.gradle.kts | head -n1)"
android_version_code="$(sed -n 's/^[[:space:]]*versionCode = \([0-9][0-9]*\)/\1/p' android/app/build.gradle.kts | head -n1)"
if [[ "$android_version_name" != "$expected_version" ]]; then
  echo "Android versionName mismatch: expected $expected_version, found ${android_version_name:-<missing>}" >&2
  exit 1
fi
if [[ "$android_version_code" != "$expected_android_code" ]]; then
  echo "Android versionCode mismatch: expected $expected_android_code, found ${android_version_code:-<missing>}" >&2
  exit 1
fi

echo "release consistency check passed for v$expected_version across API, web, extension, and Android"
