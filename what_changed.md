# What Changed

## Current milestone

Phase 57 — v0.6.1 complete release verification gate.

## Completed

- Advanced the repository release version from `0.6.0` to `0.6.1`.
- Aligned `/api/v1/meta` and its regression test with `0.6.1`.
- Aligned the web package, browser extension package/Manifest V3, and Android application with `0.6.1`.
- Assigned Android version code `6001` using the deterministic semantic-version mapping already enforced by release consistency validation.
- Expanded `make release-verify` so the top-level release gate covers backend, web, browser extension, and Android checks instead of stopping after the web build.
- Added the dedicated `docs/releases/v0.6.1.md` verification record.
- Updated the changelog and README to describe the new release gate and current version.
- Kept the release migration-free and independent of production credentials/data.

## Verification state

The v0.6.1 work is prepared on branch `release/v0.6.1` from the current `main` baseline. Repository-level inspection confirms the intended version surfaces and release-gate changes are present.

The GitHub connector does not execute arbitrary local shell commands in this workflow, so the final release state must still be validated by GitHub Actions and/or a local checkout before the branch is merged and a release tag is published.

Do not claim zero bugs or full verification solely from static repository inspection.

## Release gate

The intended complete local verification command is:

```text
make release-verify
```

This now expands to:

```text
make release-check
make release-readiness
make vet
make test
make build
make web-build
make extension-check
make extension-test
make android-lint
make android-test
make android-build
```

## Next steps

1. Create a pull request from `release/v0.6.1` into `main`.
2. Wait for all applicable CI, browser, Android, database, and CodeQL/security checks.
3. Fix any concrete failing check before merge; do not bypass a failing gate.
4. Re-run the complete applicable verification set after each fix.
5. Review the final diff and release metadata for unintended changes.
6. Merge only after the required gates are green.
7. Publish/tag `v0.6.1` only from the verified final main commit.

## Engineering rule

A passing release-consistency check proves version metadata consistency. A passing release-readiness check proves repository release documentation and module-verification requirements are satisfied. Neither alone proves the application is bug-free. Functional, integration, database, browser, Android, and security gates remain authoritative.
