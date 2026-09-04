# StockPilot — Work Continuity Log

## Current milestone

Phase 37 — v0.2.6 Warehouse & Location Valuation.

## Repository state

- Default branch: `main`.
- v0.1.0-preview.1 through v0.1.10 release preparation is merged.
- v0.2.0 reporting foundation is merged on `main`.
- v0.2.1 inventory aging implementation is merged on `main`.
- v0.2.2 expiry-risk reporting release is published on GitHub.
- v0.2.3 export-hardening release is published on GitHub.
- v0.2.4 stock movement velocity implementation and release documentation are merged on `main`; stable publication must still be verified before claiming it is released.
- v0.2.5 supplier performance implementation and release documentation are merged on `main`.
- v0.2.6 warehouse/location valuation implementation and release documentation are now merged on `main`.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.2.6

- [x] Added warehouse/location valuation domain contracts.
- [x] Extended the reporting repository contract with warehouse valuation aggregation.
- [x] Added PostgreSQL aggregation across inventory balances, products, locations, and warehouses.
- [x] Added location-level on-hand units, valuation, currency, and distinct active-product counts.
- [x] Added warehouse-level totals grouped by currency.
- [x] Added bounded HTTP JSON reporting with a 1–5000 row limit.
- [x] Added formula-safe CSV export with no-store/no-cache headers.
- [x] Registered `GET /api/v1/reports/warehouse-valuation`.
- [x] Added HTTP regression coverage for defaults, limits, JSON output, and CSV safety.
- [x] Added `docs/RELEASE_NOTES_v0.2.6.md`.
- [x] Updated `CHANGELOG.md`.
- [x] Updated the public metadata endpoint version to `0.2.6`.

## v0.2.6 release gates

- [ ] Confirm the final `main` CI run for the release candidate completes successfully.
- [ ] Confirm Go formatting, vet, unit, race, and server-build verification.
- [ ] Confirm PostgreSQL migration/integration verification.
- [ ] Confirm Web quality/type-check/build verification.
- [ ] Confirm CodeQL/security verification.
- [ ] Confirm configured Android, browser-companion, E2E, and accessibility checks.
- [ ] Confirm restore/rollback and reproducible-artifact gates where applicable.
- [ ] Resolve blocker/critical defects.
- [ ] Verify that GitHub does not already contain `v0.2.6`.
- [ ] Only after applicable gates pass: publish `v0.2.6` as a stable, non-prerelease GitHub Release.

## Release metadata

- Title: `StockPilot v0.2.6 — Warehouse & Location Valuation`
- Tag: `v0.2.6`
- Release type: Stable after verification
- Prerelease: Off
- Latest: On, if intended as the current latest stable release
- Date: 2026-09-04

## Next implementation priority after v0.2.6

1. Replenishment history and recommendation effectiveness metrics.
2. Cursor/streaming support for large report datasets.
3. Supplier/product purchasing trend series.
4. Warehouse/location valuation drill-down by product and lot.
5. Broader E2E, accessibility, Android, restore/rollback, and reproducible-artifact hardening.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
