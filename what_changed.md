# StockPilot — Work Continuity Log

## Current milestone

Phase 38 — v0.2.7 Replenishment Readiness.

## Repository state

- Default branch: `main`.
- v0.1.0-preview.1 through v0.1.10 release preparation is merged.
- v0.2.0 reporting foundation is merged on `main`.
- v0.2.1 inventory aging implementation is merged on `main`.
- v0.2.2 expiry-risk reporting release is published on GitHub.
- v0.2.3 export-hardening release is published on GitHub.
- v0.2.4 stock movement velocity implementation and release documentation are merged on `main`; stable publication must still be verified before claiming it is released.
- v0.2.5 supplier performance implementation and release documentation are merged on `main`.
- v0.2.6 warehouse/location valuation implementation and release documentation are merged on `main`.
- v0.2.7 release scope and release notes are now prepared on `main`; implementation and verification are still required before publication.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## v0.2.7 release scope

- [ ] Connect existing reorder suggestions with authoritative recent outbound movement data.
- [ ] Add replenishment readiness domain/report contracts.
- [ ] Add recent outbound units and average daily outbound velocity to the replenishment report.
- [ ] Add estimated days of cover when outbound velocity is positive.
- [ ] Add deterministic advisory readiness/risk classification.
- [ ] Add bounded authenticated JSON reporting.
- [ ] Add optional formula-safe CSV export.
- [ ] Integrate replenishment readiness into Reports & Analytics.
- [ ] Add unit, repository, HTTP, and integration regression coverage.
- [ ] Keep the feature read-only; do not automatically create purchase orders or mutate inventory.
- [ ] Prefer existing authoritative records and avoid a migration unless implementation evidence requires persisted recommendation snapshots.

## v0.2.7 release gates

- [ ] Confirm the v0.2.7 implementation is complete on `main`.
- [ ] Confirm Go formatting, vet, unit, race, and server-build verification.
- [ ] Confirm PostgreSQL migration/readiness and integration verification.
- [ ] Confirm Web quality/type-check/build verification.
- [ ] Confirm CodeQL/security verification.
- [ ] Confirm configured Android and browser-companion checks.
- [ ] Confirm E2E/accessibility checks where configured.
- [ ] Confirm restore/rollback and reproducible-artifact gates where applicable.
- [ ] Resolve blocker/critical defects.
- [ ] Verify that GitHub does not already contain `v0.2.7`.
- [ ] Only after applicable gates pass: publish `v0.2.7` as a stable, non-prerelease GitHub Release.

## Release metadata

- Title: `StockPilot v0.2.7 — Replenishment Readiness`
- Tag: `v0.2.7`
- Release type: Stable after verification
- Prerelease: Off
- Latest: On, if intended as the current latest stable release
- Date: 2026-09-04
- Release notes: `docs/RELEASE_NOTES_v0.2.7.md`

## Implementation order

1. Replenishment report domain contracts.
2. Repository aggregation over reorder configuration and movement history.
3. HTTP endpoint, bounds, authorization, and CSV export.
4. Web Reports & Analytics integration.
5. Regression/integration coverage.
6. Documentation and release metadata.
7. Full verification gates.
8. Stable GitHub publication.

## Next milestones after v0.2.7

1. Cursor/streaming support for large report datasets.
2. Supplier/product purchasing trend series.
3. Warehouse/location valuation drill-down by product and lot.
4. Broader E2E, accessibility, Android, restore/rollback, and reproducible-artifact hardening.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
