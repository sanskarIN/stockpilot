# What Changed

## Current milestone

Phase 56 — v0.6.0 client version-drift hardening.

## Completed

- Fixed a concrete release-integrity defect where the web package still reported `0.1.0` while the repository/API release version was `0.6.0`.
- Aligned the browser extension package and Manifest V3 version with `0.6.0`.
- Aligned the Android app `versionName` with `0.6.0` and assigned deterministic `versionCode` `6000`.
- Strengthened `scripts/check-release-consistency.sh` so supported client surfaces are checked against `VERSION` instead of allowing silent version drift.
- Kept the change migration-free and independent of production credentials/data.

## Verification state

The fix is prepared in `fix/v0.6.0-version-drift` and submitted as PR #71. Repository inspection confirms the previously discovered client-version drift has been corrected in the changed files.

Full release confidence is still gated on GitHub Actions and applicable Android/extension/web/Go/security checks. Do not claim zero bugs or publish/move a release tag solely from static inspection.

## Next steps

1. Wait for PR #71 CI/security checks.
2. If a check fails, inspect the failing job and fix the concrete defect before merging.
3. Re-run the complete applicable verification set after each fix.
4. Review the final PR diff for unintended changes.
5. Merge only after required gates are green.

## Engineering rule

A green version-consistency check proves only release metadata consistency. It does not prove the entire application is bug-free. Functional, integration, browser, Android, extension, database, and security gates remain authoritative.