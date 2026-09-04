# StockPilot — Work Continuity Log

## Current milestone

Phase 39 — v0.2.8 Large-Report Scalability & Purchasing Trends.

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
- v0.2.7 replenishment-readiness implementation and release documentation are merged on `main`.
- v0.2.7 stable publication remains gated on verification results; do not claim the GitHub Release is published until confirmed.
- v0.2.8 planning is now active on `main`.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## v0.2.8 active scope

### Large-report scalability

- [ ] Audit existing report pagination, aggregation, and database indexes.
- [ ] Define additive cursor/keyset pagination or streaming contracts.
- [ ] Preserve deterministic ordering and opaque continuation semantics.
- [ ] Bound page sizes, date windows, query work, and memory growth.
- [ ] Add first-page, continuation, end-of-stream, ordering, and duplicate-prevention tests.

### Purchasing trends

- [ ] Define supplier/product trend domain contracts from existing purchasing records.
- [ ] Aggregate ordered, received, open units and purchasing value server-side.
- [ ] Keep currencies separated; never introduce implicit FX conversion.
- [ ] Add bounded HTTP reporting endpoints.
- [ ] Add formula-safe CSV exports with non-cacheable download headers.
- [ ] Integrate trend views into Reports & Analytics.
- [ ] Add loading, empty, error, and accessibility states.

## v0.2.8 engineering rules

- Preserve existing endpoint semantics unless a compatibility review explicitly approves a change.
- Prefer new additive endpoints/contracts when pagination semantics would otherwise alter existing clients.
- Reuse authoritative purchase orders, purchase-order lines, products, suppliers, and existing indexes.
- Keep reporting read-only.
- No automatic purchase-order creation or inventory mutation.
- No schema migration unless implementation evidence proves the existing schema insufficient.

## v0.2.8 implementation sequence

1. Audit report query/index patterns.
2. Add reusable bounded pagination primitives.
3. Add cursor/keyset tests.
4. Define purchasing trend contracts.
5. Implement PostgreSQL trend aggregation.
6. Add HTTP endpoints and CSV exports.
7. Integrate Reports & Analytics.
8. Add PostgreSQL/integration/security regression coverage.
9. Run Go, Web, PostgreSQL, CodeQL, Android/browser, E2E/accessibility, rollback, and artifact gates where configured.
10. Publish v0.2.8 only after implementation and verification are complete.

## Future milestones after v0.2.8

1. Warehouse/location valuation drill-down by product and lot.
2. Broader E2E, accessibility, Android, restore/rollback, and reproducible-artifact hardening.
3. Pre-1.0 API stability and operational-readiness review.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
