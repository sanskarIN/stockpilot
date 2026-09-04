# StockPilot — Work Continuity Log

## Current milestone

Phase 35 — v0.2.4 Stock Movement Velocity and release preparation.

## Repository state

- Default branch: `main`.
- v0.1.0-preview.1 through v0.1.10 release preparation is merged.
- v0.2.0 reporting foundation is merged on `main`.
- v0.2.1 inventory aging implementation is merged on `main`.
- v0.2.2 expiry-risk reporting release is published on GitHub.
- v0.2.3 export-hardening changes are merged on `main`; stable publication remains subject to its final release gates.
- v0.2.4 stock movement velocity implementation and release documentation are merged on `main`.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.2.4

- [x] Added stock movement history domain contracts.
- [x] Added a separate repository reader contract without expanding the mutation-focused Inventory interface.
- [x] Added PostgreSQL aggregation over authoritative `stock_movements` records.
- [x] Added configurable 1–365 day reporting window with 30-day default.
- [x] Added movement count, inbound units, outbound units, net units, and average daily outbound metrics.
- [x] Added bounded HTTP JSON reporting.
- [x] Added formula-safe CSV export with existing privacy-oriented download headers.
- [x] Added web API/types and Reports workspace movement-velocity panel.
- [x] Added HTTP/domain regression coverage for report bounds and contracts.
- [x] Updated `ROADMAP.md`.
- [x] Updated `CHANGELOG.md`.
- [x] Added `docs/RELEASE_NOTES_v0.2.4.md`.

## v0.2.4 release gates

- [ ] Confirm the final `main` CI run completes successfully for the release candidate.
- [ ] Confirm Go formatting, vet, unit, race, and server-build verification.
- [ ] Confirm PostgreSQL migration/integration verification.
- [ ] Confirm Web quality/build verification.
- [ ] Confirm CodeQL/security verification.
- [ ] Confirm any configured Android, browser-companion, E2E, and accessibility checks.
- [ ] Confirm restore/rollback and reproducible-artifact gates where applicable.
- [ ] Resolve blocker/critical defects.
- [ ] Only after applicable gates pass: publish `v0.2.4` as a stable, non-prerelease GitHub Release.

## Release metadata

- Title: `StockPilot v0.2.4 — Stock Movement Velocity`
- Tag: `v0.2.4`
- Release type: Stable after verification
- Prerelease: Off
- Latest: On, if intended as the current latest stable release
- Date: 2026-09-04

## Next implementation priority

1. Supplier purchasing totals and lead-time tracking.
2. Warehouse/location valuation breakdown.
3. Replenishment history and recommendation effectiveness metrics.
4. Cursor/streaming support for large report datasets.
5. Broader E2E, accessibility, Android, restore/rollback, and reproducible-artifact hardening.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
