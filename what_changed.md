# StockPilot — Work Continuity Log

## Current milestone

Phase 22 — v0.1.3 product catalog CSV export and release preparation.

## Repository state

- Default branch: `main`.
- v0.1.0-preview.1 release preparation is merged.
- v0.1.1 maintenance-release preparation is merged.
- v0.1.2 CSV serialization foundation is merged.
- The transactional CSV product-import workflow remains in `main` from PR #38.
- v0.1.3 now moves the CSV foundation into the first bounded application-level export endpoint.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.1.3 preparation

- [x] Added `internal/httpapi/catalog_export.go` with the product catalog CSV endpoint.
- [x] Registered `GET /api/v1/products/export.csv`.
- [x] Reused existing product filters: `q`, `categoryId`, `supplierId`, `activeOnly`, `limit`, and `offset`.
- [x] Added default and maximum export bounds, with a hard maximum of 5,000 rows per request.
- [x] Added deterministic catalog CSV columns and UTC timestamp serialization.
- [x] Enabled formula-safe serialization for downloadable catalog values.
- [x] Added download response headers with a deterministic filename.
- [x] Added `internal/httpapi/catalog_export_test.go` for export-bound normalization.
- [x] Added `docs/RELEASE_NOTES_v0.1.3.md`.
- [x] Began the v0.1.3 changelog entry; the connected GitHub contents operation reported a stale SHA conflict while updating `CHANGELOG.md`, so the existing changelog was intentionally left untouched rather than overwriting concurrent repository state.

## v0.1.3 release gates

- [ ] Confirm the exact v0.1.3 release commit on `main` after all intended changes are merged.
- [ ] Run `gofmt`, `go vet ./...`, normal tests, and race-enabled Go tests.
- [ ] Run web typecheck and production build.
- [ ] Run PostgreSQL migration/readiness smoke testing.
- [ ] Run Android lint/tests/build and release-networking/security checks.
- [ ] Run browser-companion manifest and unit checks.
- [ ] Run configured CodeQL checks for Go and JavaScript/TypeScript.
- [ ] Verify catalog export route, filters, bounds, schema, formula safety, response headers, and repository error handling.
- [ ] Complete deployed authorization/CSRF and authentication/session smoke checks.
- [ ] Verify no secrets or credentials are included in exported data or artifacts.
- [ ] Create immutable `v0.1.3` tag on the verified commit.
- [ ] Publish the GitHub Release with the prepared notes and mark it pre-release until the full gate set passes.
- [ ] Perform post-release smoke testing against the published tag/artifacts.

## Publication details

- Version: `v0.1.3`
- Release title: `StockPilot v0.1.3 — Product Catalog CSV Export`
- Git tag: `v0.1.3`
- Release class: normal pre-1.0 feature release
- Pre-release: yes until every release gate passes
- Release notes: `docs/RELEASE_NOTES_v0.1.3.md`

## Next exact development tasks after v0.1.3

1. Add inventory-balance CSV export.
2. Add low-stock and reorder-suggestion CSV export.
3. Add lot-inventory and expiry-risk CSV export.
4. Add purchase-order and receiving export contracts.
5. Add export-specific audit coverage for sensitive datasets.
6. Add deterministic streaming behavior for large exports.
7. Add web download controls and accessible export feedback.
8. Add integration coverage around authorization and export headers.
9. Continue toward v0.2.x analytics and operational reporting.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
