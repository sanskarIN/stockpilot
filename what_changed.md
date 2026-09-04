# StockPilot — Work Continuity Log

## Current milestone

Phase 25 — v0.1.6 purchasing and receiving CSV export and release preparation.

## Repository state

- Default branch: `main`.
- v0.1.0-preview.1 release preparation is merged.
- v0.1.1 maintenance-release preparation is merged.
- v0.1.2 CSV serialization foundation is merged.
- The transactional CSV product-import workflow remains in `main` from PR #38.
- v0.1.3 added the first bounded application-level product catalog CSV export.
- v0.1.4 extended the export surface to inventory balances, low-stock data, and reorder suggestions.
- v0.1.5 extended exports to lot inventory and expiry-risk filtering.
- v0.1.6 now extends exports to purchase-order lines and current receiving progress.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.1.6 preparation

- [x] Added `domain.PurchaseOrderExportRow` as a dedicated flattened export model.
- [x] Added `repository.Orders.ListOrderExportRows` to keep export retrieval separate from the normal order graph API.
- [x] Added a PostgreSQL CTE-based export query that selects purchase orders first and then returns deterministic line rows.
- [x] Added `GET /api/v1/orders/export.csv`.
- [x] Added optional purchase-order `status` filtering.
- [x] Added default 500 and maximum 5,000 selected-purchase-order export bounds.
- [x] Added line-level `quantity`, `received`, `remaining`, `unitCostMinor`, and `lineTotalMinor` fields.
- [x] Added UTC RFC 3339 timestamp serialization.
- [x] Reused the shared formula-safe CSV serializer.
- [x] Added deterministic browser download filename and CSV content type.
- [x] Registered the new purchase-order export route without changing existing purchasing or receiving endpoints.
- [x] Added focused tests for bounds, status validation, CSV schema, formula safety, totals, and UTC timestamps.
- [x] Added `docs/RELEASE_NOTES_v0.1.6.md`.
- [x] Added the v0.1.6 `CHANGELOG.md` entry.

## v0.1.6 release gates

- [ ] Run `gofmt`, `go vet ./...`, normal tests, and race-enabled Go tests.
- [ ] Verify every production and test implementation of `repository.Orders` implements `ListOrderExportRows`.
- [ ] Verify purchase-order export status filtering for every supported status.
- [ ] Verify invalid status values return HTTP 400.
- [ ] Verify negative offsets normalize to zero and limits above 5,000 clamp safely.
- [ ] Verify deterministic schema, order/line ordering, formula safety, UTC timestamps, filename, and content type.
- [ ] Verify received and remaining quantities against the authoritative purchase-order state.
- [ ] Run PostgreSQL migration/readiness smoke testing.
- [ ] Run web typecheck and production build.
- [ ] Run Android lint/tests/build and release-networking/security checks.
- [ ] Run browser-companion manifest and unit checks.
- [ ] Run configured CodeQL checks for Go and JavaScript/TypeScript.
- [ ] Complete deployed authorization/CSRF and authentication/session smoke checks.
- [ ] Verify no credentials, secrets, or payment information are included in exported data.
- [ ] Create immutable `v0.1.6` tag on the verified commit.
- [ ] Publish the GitHub Release with the prepared notes and keep it pre-release until every gate passes.
- [ ] Perform post-release smoke testing against the published tag/artifacts.

## Publication details

- Version: `v0.1.6`
- Release title: `StockPilot v0.1.6 — Purchasing & Receiving CSV Export`
- Git tag: `v0.1.6`
- Release class: normal pre-1.0 feature release
- Pre-release: yes until every release gate passes
- Release notes: `docs/RELEASE_NOTES_v0.1.6.md`

## Known verification limitation in this workspace

- The repository was updated directly through the connected GitHub integration.
- A local `go test ./...` cannot be treated as completed unless executed in an environment with the repository and required dependencies available.
- The release must remain pending until GitHub Actions or a local developer environment confirms the full build/test matrix.

## Next exact development tasks after v0.1.6

1. Add export-specific audit coverage for sensitive datasets.
2. Add a dedicated receipt-event/history export backed by the stock-movement/audit model where the data contract supports it.
3. Add deterministic streaming behavior for large exports.
4. Add web download controls and accessible export feedback.
5. Add integration coverage around authorization, CSRF, and export headers.
6. Add export job lifecycle/status endpoints for large datasets.
7. Add richer expiry-risk classification and operational alerts.
8. Continue toward v0.2.x analytics and operational reporting.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
