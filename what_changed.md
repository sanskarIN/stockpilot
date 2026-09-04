# StockPilot — Work Continuity Log

## Current milestone

Phase 39 — v0.2.8 Large-Report Scalability & Purchasing Trends preparation.

## Repository state

- Default branch: `main`.
- v0.1.0-preview.1 through v0.1.10 release preparation is merged.
- v0.2.0 reporting foundation is merged on `main`.
- v0.2.1 inventory aging implementation is merged on `main`.
- v0.2.2 expiry-risk reporting release is published on GitHub.
- v0.2.3 export-hardening release is published on GitHub.
- v0.2.4 stock movement velocity implementation and release documentation are merged on `main`.
- v0.2.5 supplier performance implementation and release documentation are merged on `main`.
- v0.2.6 warehouse/location valuation implementation and release documentation are merged on `main`; GitHub release is published.
- v0.2.7 replenishment-readiness implementation and release documentation are now merged on `main`.
- v0.2.7 stable publication remains gated on verification results; do not claim the GitHub Release is published until confirmed.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## v0.2.7 completed scope

- [x] Connect existing reorder suggestions with authoritative recent outbound movement data.
- [x] Add replenishment readiness domain/report contracts.
- [x] Add recent outbound units and average daily outbound velocity to the replenishment report.
- [x] Add estimated days of cover when outbound velocity is positive.
- [x] Add deterministic advisory readiness/risk classification.
- [x] Add bounded authenticated JSON reporting.
- [x] Add optional formula-safe CSV export.
- [x] Register `GET /api/v1/reports/replenishment-readiness`.
- [x] Update API metadata to `0.2.7`.
- [x] Add domain regression coverage.
- [x] Keep the feature read-only; no automatic purchase orders or inventory mutation.
- [x] Reuse existing authoritative records; no migration was introduced.

## v0.2.7 release gates

- [ ] Confirm Go formatting, vet, unit, race, and server-build verification.
- [ ] Confirm PostgreSQL migration/readiness and integration verification.
- [ ] Confirm Web quality/type-check/build verification.
- [ ] Confirm CodeQL/security verification.
- [ ] Confirm configured Android and browser-companion checks.
- [ ] Confirm E2E/accessibility checks where configured.
- [ ] Confirm restore/rollback and reproducible-artifact gates where applicable.
- [ ] Resolve blocker/critical defects.
- [x] Verify GitHub currently has no `v0.2.7` release.
- [ ] Only after applicable gates pass: publish `v0.2.7` as a stable, non-prerelease GitHub Release.

## v0.2.8 preparation scope

### Large-report scalability

- Design cursor/keyset pagination or streaming for large reporting datasets.
- Preserve deterministic ordering and continuation semantics.
- Bound memory usage and reject unsafe page/window inputs.
- Prefer indexed, server-side iteration over unbounded in-memory aggregation.
- Add regression coverage for first page, continuation, end-of-stream, duplicate prevention, and ordering stability.

### Purchasing trends

- Add supplier/product purchasing trend series from authoritative purchase-order and purchase-order-line records.
- Support bounded date windows and deterministic series ordering.
- Expose ordered units, received units, open units, and purchasing value where the existing schema supports them without currency conversion assumptions.
- Add formula-safe CSV export and non-cacheable download behavior.
- Integrate trends into Reports & Analytics with clear loading, empty, error, and accessibility states.

### v0.2.8 safety rules

- Reporting remains read-only.
- No automatic purchase-order creation or inventory mutation is part of this milestone.
- No schema migration unless implementation evidence proves existing data cannot support the feature.
- Do not publish v0.2.8 until the applicable verification gates are green.

## v0.2.8 implementation order

1. Audit existing report pagination/aggregation patterns and indexes.
2. Define cursor/streaming contracts without breaking current endpoints.
3. Implement server-side continuation and bounds.
4. Define purchasing trend domain contracts.
5. Implement PostgreSQL trend aggregation.
6. Add HTTP endpoints and safe CSV exports.
7. Add Web Reports & Analytics trend views.
8. Add unit, HTTP, PostgreSQL, and integration coverage.
9. Run all repository verification gates.
10. Prepare and publish v0.2.8 only after verification.

## Future milestones after v0.2.8

1. Warehouse/location valuation drill-down by product and lot.
2. Broader E2E, accessibility, Android, restore/rollback, and reproducible-artifact hardening.
3. Pre-1.0 API stability and operational-readiness review.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
