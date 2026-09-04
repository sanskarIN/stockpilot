# StockPilot — Work Continuity Log

## Current milestone

Phase 32 — v0.2.1 Inventory Aging implementation, with stable-release hardening gates for v0.2.1.

## Repository state

- Default branch: `main`.
- v0.1.0-preview.1 through v0.1.10 release preparation is merged.
- v0.2.0 reporting foundation is merged on `main`.
- v0.2.1 inventory aging implementation is now merged on `main`.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.2.1

- [x] Added inventory-aging domain report contracts.
- [x] Added deterministic aging buckets: 0-30, 31-60, 61-90, 91-180, and 181+ days.
- [x] Added boundary regression coverage for aging buckets.
- [x] Added repository contract for inventory aging.
- [x] Added PostgreSQL inventory aging query using authoritative movement history with a balance timestamp fallback.
- [x] Added bounded HTTP inventory aging report endpoint.
- [x] Added formula-safe CSV output through the aging report endpoint.
- [x] Added HTTP coverage for aging export limit normalization.
- [x] Added typed web aging models and API client.
- [x] Added inventory aging panel and export action to Reports & Analytics.
- [x] Added v0.2.1 release notes.

## v0.2.1 stable-release gates

- [ ] Run complete Go formatting, vet, unit, race, and integration verification.
- [ ] Run complete Web, Android, browser-companion, and CodeQL verification where configured.
- [ ] Complete browser end-to-end and accessibility checks required for stable quality.
- [ ] Complete production restore and migration rollback rehearsal.
- [ ] Verify reproducible server, web, Android, and extension release artifacts.
- [ ] Resolve all blocker/critical defects.
- [ ] Only after all gates pass: publish `v0.2.1` as a stable, non-prerelease GitHub Release.

## Aging semantics

Inventory aging is currently defined as elapsed whole days since the most recent persisted stock movement for a positive product/location/lot balance. When no matching movement is available, the balance `updated_at` timestamp is used. This definition is intentionally explicit so future FIFO/layer-based aging can be introduced as a separate, compatible reporting model rather than silently changing existing results.

## Next exact development tasks

1. Add PostgreSQL integration fixtures and assertions for aging boundaries and movement selection.
2. Add access/export authorization regression coverage for the aging report.
3. Add configurable expiry-risk reporting windows.
4. Add stock movement history and velocity analytics.
5. Add supplier purchasing totals and lead-time measurements.
6. Add warehouse/location valuation breakdowns.
7. Add replenishment effectiveness metrics.
8. Add cursor/streaming report readers and large-export lifecycle endpoints.
9. Expand browser E2E, Android instrumentation, accessibility, restore, compatibility, and reproducible-artifact gates.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
