# StockPilot — Work Continuity Log

## Current milestone

Phase 36 — v0.2.5 Supplier Performance and release preparation.

## Repository state

- Default branch: `main`.
- v0.1.0-preview.1 through v0.1.10 release preparation is merged.
- v0.2.0 reporting foundation is merged on `main`.
- v0.2.1 inventory aging implementation is merged on `main`.
- v0.2.2 expiry-risk reporting release is published on GitHub.
- v0.2.3 export-hardening release is published on GitHub.
- v0.2.4 stock movement velocity implementation and release documentation are merged on `main`; stable publication must still be verified before claiming it is released.
- v0.2.5 supplier performance implementation and release documentation are now merged on `main`.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.2.5

- [x] Added supplier performance domain report contracts.
- [x] Extended the reporting repository contract with supplier performance aggregation.
- [x] Added PostgreSQL supplier aggregation using purchase orders, order lines, suppliers, and receive movements.
- [x] Added ordered/received/open units and purchasing value metrics.
- [x] Added observed receipt lead-time measurement.
- [x] Added completed-order and on-time-order metrics.
- [x] Added bounded HTTP JSON reporting with a 1–365 day window and 1–5000 row limit.
- [x] Added formula-safe CSV export with no-store/no-cache headers.
- [x] Registered `GET /api/v1/reports/supplier-performance`.
- [x] Added HTTP normalization regression tests.
- [x] Added web API/types and supplier reliability panel.
- [x] Updated `CHANGELOG.md`.
- [x] Added `docs/RELEASE_NOTES_v0.2.5.md`.

## v0.2.5 release gates

- [ ] Confirm the final `main` CI run for the release candidate completes successfully.
- [ ] Confirm Go formatting, vet, unit, race, and server-build verification.
- [ ] Confirm PostgreSQL migration/integration verification.
- [ ] Confirm Web quality/type-check/build verification.
- [ ] Confirm CodeQL/security verification.
- [ ] Confirm configured Android, browser-companion, E2E, and accessibility checks.
- [ ] Confirm restore/rollback and reproducible-artifact gates where applicable.
- [ ] Resolve blocker/critical defects.
- [ ] Verify that GitHub does not already contain `v0.2.5`.
- [ ] Only after applicable gates pass: publish `v0.2.5` as a stable, non-prerelease GitHub Release.

## Release metadata

- Title: `StockPilot v0.2.5 — Supplier Performance`
- Tag: `v0.2.5`
- Release type: Stable after verification
- Prerelease: Off
- Latest: On, if intended as the current latest stable release
- Date: 2026-09-04

## Next implementation priority after v0.2.5

1. Warehouse/location valuation breakdown.
2. Replenishment history and recommendation effectiveness metrics.
3. Cursor/streaming support for large report datasets.
4. Supplier/product purchasing trend series.
5. Broader E2E, accessibility, Android, restore/rollback, and reproducible-artifact hardening.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
