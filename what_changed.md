# StockPilot — Work Continuity Log

## Current milestone

Phase 23 — v0.1.4 inventory and reorder CSV exports and release preparation.

## Repository state

- Default branch: `main`.
- v0.1.0-preview.1 release preparation is merged.
- v0.1.1 maintenance-release preparation is merged.
- v0.1.2 CSV serialization foundation is merged.
- The transactional CSV product-import workflow remains in `main` from PR #38.
- v0.1.3 added the first bounded application-level product catalog CSV export.
- v0.1.4 now extends the export surface to inventory balances, low-stock data, and reorder suggestions.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.1.4 preparation

- [x] Extended `repository.Inventory` with bounded `ListBalances` pagination.
- [x] Added PostgreSQL inventory-balance listing with deterministic product/location/lot ordering.
- [x] Added `GET /api/v1/inventory/export.csv`.
- [x] Added `GET /api/v1/inventory/low-stock/export.csv`.
- [x] Added `GET /api/v1/inventory/reorder-suggestions/export.csv`.
- [x] Reused the shared formula-safe CSV serializer.
- [x] Added bounded export normalization with a 5,000-row application maximum for inventory-balance exports.
- [x] Added deterministic CSV schemas and UTC RFC 3339 timestamp serialization.
- [x] Added deterministic browser download filenames and CSV content types.
- [x] Restored and extended the HTTP API fake store after the inventory repository interface changed.
- [x] Added focused export tests covering bounds, headers, formula safety, timestamps, low-stock output, and reorder-suggestion output.
- [x] Added `docs/RELEASE_NOTES_v0.1.4.md`.

## v0.1.4 release gates

- [ ] Run `gofmt`, `go vet ./...`, normal tests, and race-enabled Go tests.
- [ ] Verify the repository interface change is implemented by every production and test store.
- [ ] Run PostgreSQL migration/readiness smoke testing.
- [ ] Verify inventory export pagination, 5,000-row application bound, ordering, and CSV schema.
- [ ] Verify low-stock and reorder-suggestion export schemas and repository-side safety caps.
- [ ] Verify formula-safe serialization and UTC timestamps.
- [ ] Verify download content types and filenames.
- [ ] Run web typecheck and production build.
- [ ] Run Android lint/tests/build and release-networking/security checks.
- [ ] Run browser-companion manifest and unit checks.
- [ ] Run configured CodeQL checks for Go and JavaScript/TypeScript.
- [ ] Complete deployed authorization/CSRF and authentication/session smoke checks.
- [ ] Verify no secrets or credentials are included in exported data or artifacts.
- [ ] Create immutable `v0.1.4` tag on the verified commit.
- [ ] Publish the GitHub Release with the prepared notes and mark it pre-release until the full gate set passes.
- [ ] Perform post-release smoke testing against the published tag/artifacts.

## Publication details

- Version: `v0.1.4`
- Release title: `StockPilot v0.1.4 — Inventory & Reorder CSV Exports`
- Git tag: `v0.1.4`
- Release class: normal pre-1.0 feature release
- Pre-release: yes until every release gate passes
- Release notes: `docs/RELEASE_NOTES_v0.1.4.md`

## Known verification limitation in this workspace

- The repository was updated directly through the connected GitHub integration.
- A local `go test ./...` could not be executed in this workspace because outbound DNS/network access from the local execution environment is unavailable.
- The release must therefore remain pending until GitHub Actions or a local developer environment confirms the full build/test matrix.

## Next exact development tasks after v0.1.4

1. Add lot-inventory and expiry-risk CSV export.
2. Add purchase-order and receiving CSV export contracts.
3. Add export-specific audit coverage for sensitive datasets.
4. Add deterministic streaming behavior for large exports.
5. Add web download controls and accessible export feedback.
6. Add integration coverage around authorization and export headers.
7. Add export job lifecycle/status endpoints for large datasets.
8. Continue toward v0.2.x analytics and operational reporting.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
