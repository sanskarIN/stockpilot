# StockPilot — Work Continuity Log

## Current milestone

Phase 34 — v0.2.3 Export Hardening and stable-release preparation.

## Repository state

- Default branch: `main`.
- v0.1.0-preview.1 through v0.1.10 release preparation is merged.
- v0.2.0 reporting foundation is merged on `main`.
- v0.2.1 inventory aging implementation is merged on `main`.
- v0.2.2 expiry-risk reporting release is published on GitHub.
- v0.2.3 export-hardening changes are merged on `main`.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.2.3

- [x] Corrected audit-export metadata formula sanitization.
- [x] Preserved shared formula-safe CSV serialization.
- [x] Corrected receipt-history export timestamp regression expectations to canonical UTC output.
- [x] Added `docs/RELEASE_NOTES_v0.2.3.md`.
- [x] Updated `CHANGELOG.md` for v0.2.3.
- [x] Updated `ROADMAP.md` for the v0.2.3 milestone.

## v0.2.3 release gates

- [ ] Confirm the final `main` CI run completes successfully for the release candidate.
- [ ] Confirm Go formatting, vet, unit, race, and server-build verification.
- [ ] Confirm PostgreSQL migration/integration verification.
- [ ] Confirm Web quality/build verification.
- [ ] Confirm CodeQL/security verification.
- [ ] Confirm any configured Android, browser-companion, E2E, and accessibility checks.
- [ ] Confirm restore/rollback and reproducible-artifact gates where applicable.
- [ ] Resolve blocker/critical defects.
- [ ] Only after applicable gates pass: publish `v0.2.3` as a stable, non-prerelease GitHub Release.

## Release metadata

- Title: `StockPilot v0.2.3 — Export Hardening`
- Tag: `v0.2.3`
- Release type: Stable
- Prerelease: Off
- Latest: On, if intended as the current latest stable release
- Date: 2026-09-04

## Next implementation priority

1. Stock movement history and velocity report.
2. Supplier purchasing totals and lead-time tracking.
3. Warehouse/location valuation breakdown.
4. Replenishment effectiveness metrics.
5. Cursor/streaming support for large report datasets.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
