# StockPilot — Work Continuity Log

## Current milestone

Phase 20 — v0.1.1 maintenance-release preparation.

## Repository state

- Default branch: `main`.
- The v0.1.0-preview.1 release preparation has been merged.
- The transactional CSV product-import workflow remains in `main` from PR #38.
- The final preview verification run passed Go quality, PostgreSQL migration smoke testing, Web quality, and Go/JavaScript/TypeScript CodeQL.
- v0.1.1 is intentionally a maintenance release focused on release discipline, verification, and operational clarity rather than a large new feature payload.
- This continuation preserves focused, reviewable commits instead of creating meaningless commits solely to increase the commit count.

## Completed for v0.1.1 preparation

- [x] Added `docs/RELEASE_NOTES_v0.1.1.md` with scope, verification gates, upgrade guidance, limitations, and artifact requirements.
- [x] Added `docs/RELEASE_RUNBOOK_v0.1.1.md` covering source freeze, automated checks, manual checks, tagging, GitHub publication, artifacts, post-release verification, and rollback.
- [x] Recorded the v0.1.1 maintenance-release entry in `CHANGELOG.md`.
- [x] Documented the release process and publication requirements without changing the core application behavior.

## Required release gates

- [ ] Confirm the exact v0.1.1 release commit on `main` after all intended release-preparation changes are merged.
- [ ] Verify required GitHub Actions checks are green for that exact commit.
- [ ] Complete the backup/restore drill from `docs/RESTORE_DRILL.md`.
- [ ] Complete authentication, authorization, CSRF, and session-expiry smoke tests.
- [ ] Complete responsive and keyboard-navigation review.
- [ ] Complete Android lint/test, reproducible APK, HTTPS, encrypted-session, and device smoke checks.
- [ ] Complete browser-companion manifest, permission-scope, credential-storage, installation, and handoff checks.
- [ ] Review production secrets, TLS, CORS, security headers, backup retention, request IDs, and rollback procedure.
- [ ] Generate reproducible release artifacts from the exact tagged commit where applicable.
- [ ] Generate and publish SHA-256 checksums for binary artifacts.
- [ ] Create the immutable `v0.1.1` Git tag on the verified commit.
- [ ] Publish the GitHub Release as the latest stable release after every required gate passes.
- [ ] Complete post-release smoke testing.

## Publication details

- Version: `v0.1.1`
- Release title: `StockPilot v0.1.1 — Maintenance Release`
- Git tag: `v0.1.1`
- Release class: maintenance release
- Pre-release: no, provided all release gates pass
- Latest release: yes, after publication verification
- Release notes: `docs/RELEASE_NOTES_v0.1.1.md`
- Runbook: `docs/RELEASE_RUNBOOK_v0.1.1.md`

## Known limitations

- The current connected GitHub write surface does not expose release/tag creation, so the final tag and Release publication may need to be performed from the GitHub UI/API.
- CSV inventory/report export is the next product-development milestone after the maintenance release.
- Advanced analytics, expanded concurrent-inventory integration coverage, migration compatibility coverage, browser E2E, Android instrumentation, accessibility, restore automation, artifact automation, and backup-retention examples remain pending.

## Next exact tasks after v0.1.1

1. Publish v0.1.1 from the exact verified `main` commit.
2. Begin the v0.2.0 development branch for bounded/streaming CSV inventory and report export.
3. Add inventory, movement, purchasing, and audit/report export authorization and audit coverage.
4. Add inventory aging, configurable expiry-risk, movement velocity, supplier analytics, and replenishment analytics.
5. Expand concurrent inventory PostgreSQL integration and migration compatibility testing.
6. Add browser E2E, Android instrumentation, accessibility, restore verification, artifact checks, and backup-retention examples.
7. Prepare later release-candidate and stable gates for `v1.0.0`.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
