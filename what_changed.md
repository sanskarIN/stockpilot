# StockPilot — Work Continuity Log

## Current milestone

Phase 31 — v0.2.1 Inventory Aging foundation, with stable-release hardening for v0.2.0.

## Repository state

- Default branch: `main`.
- v0.1.0-preview.1 through v0.1.10 release preparation is merged.
- v0.2.0 reporting foundation is merged on `main`.
- v0.2.1 begins inventory aging with deterministic domain-level age buckets.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.2.1 preparation

- [x] Added inventory-aging domain report contracts.
- [x] Added deterministic aging buckets: 0-30, 31-60, 61-90, 91-180, and 181+ days.
- [x] Added boundary regression coverage for aging buckets.
- [x] Updated roadmap to reflect the v0.2.1 reporting task.

## v0.2.0 stable-release status

- [ ] Run complete Go formatting, vet, unit, race, and integration verification.
- [ ] Run complete Web, Android, browser-companion, and CodeQL verification where configured.
- [ ] Complete browser end-to-end and accessibility checks required for stable quality.
- [ ] Complete production restore and migration rollback rehearsal.
- [ ] Verify reproducible release artifacts.
- [ ] Resolve all blocker/critical defects.
- [ ] Only after all gates pass: publish `v0.2.0` as a stable, non-prerelease GitHub Release.

## Stable release policy

`v0.2.0` must not be labeled stable merely because the feature work is merged. Stable status requires the repository's release gates to pass. The connected GitHub workspace does not expose a trustworthy local shell execution result, so unexecuted checks are intentionally not claimed as complete.

## Next exact development tasks

1. Implement inventory aging repository/query support using authoritative inventory/movement data.
2. Add HTTP reporting and bounded CSV export for aging.
3. Add web aging report UI with filtering and refresh/error/session handling.
4. Add PostgreSQL integration coverage for aging calculations.
5. Add expiry-risk reporting with configurable windows.
6. Add stock movement history and velocity analytics.
7. Add supplier purchasing totals and lead-time measurements.
8. Add warehouse/location valuation breakdowns.
9. Add replenishment effectiveness metrics.
10. Add cursor/streaming report readers and large-export lifecycle endpoints.
11. Expand browser E2E, Android instrumentation, accessibility, restore, compatibility, and reproducible-artifact gates.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
