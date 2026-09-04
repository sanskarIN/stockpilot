# StockPilot — Work Continuity Log

## Current milestone

Phase 33 — v0.2.2 Expiry-Risk Reporting and final-release preparation.

## Repository state

- Default branch: `main`.
- v0.1.0-preview.1 through v0.1.10 release preparation is merged.
- v0.2.0 reporting foundation is merged on `main`.
- v0.2.1 inventory aging implementation is merged on `main`.
- v0.2.2 final release documentation is now merged on `main`.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.2.2

- [x] Added expiry-risk reporting scope and release documentation.
- [x] Added configurable expiry-risk windows.
- [x] Added expired, critical, warning, and safe risk classifications.
- [x] Added bounded report/export requirements.
- [x] Added formula-safe CSV requirements.
- [x] Added authorization and regression-test requirements.
- [x] Added `docs/RELEASE_NOTES_v0.2.2.md`.
- [x] Updated `CHANGELOG.md` for v0.2.2.

## Final release gates

- [ ] Confirm the latest `main` CI run completes successfully on the final release commit.
- [ ] Confirm Go formatting, vet, unit, race, and integration verification.
- [ ] Confirm PostgreSQL migration/integration verification.
- [ ] Confirm Web quality/build verification.
- [ ] Confirm Android and browser-companion verification where configured.
- [ ] Confirm CodeQL/security verification where configured.
- [ ] Confirm browser E2E and accessibility checks where configured.
- [ ] Complete production restore and migration rollback rehearsal.
- [ ] Verify reproducible server, web, Android, and extension release artifacts.
- [ ] Resolve all blocker/critical defects.
- [ ] Only after every applicable gate passes: publish `v0.2.2` as a stable, non-prerelease GitHub Release.

## Release metadata

- Title: `StockPilot v0.2.2 — Expiry-Risk Reporting`
- Tag: `v0.2.2`
- Release type: Stable
- Prerelease: Off
- Latest: On
- Date: 2026-09-04

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
