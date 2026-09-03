# StockPilot — Work Continuity Log

## Current milestone

Phase 24 — v0.1.5 lot inventory and expiry-risk CSV export and release preparation.

## Repository state

- Default branch: `main`.
- v0.1.0-preview.1 release preparation is merged.
- v0.1.1 maintenance-release preparation is merged.
- v0.1.2 CSV serialization foundation is merged.
- The transactional CSV product-import workflow remains in `main` from PR #38.
- v0.1.3 added the first bounded application-level product catalog CSV export.
- v0.1.4 extended the export surface to inventory balances, low-stock data, and reorder suggestions.
- v0.1.5 now extends exports to lot inventory and expiry-risk filtering.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.1.5 preparation

- [x] Added `GET /api/v1/inventory/lots/export.csv`.
- [x] Reused `repository.LotInventoryFilter` for product, warehouse, location, lot, and expiry filtering.
- [x] Added inclusive `expiringBy` filtering with strict `YYYY-MM-DD` parsing.
- [x] Added default 500-row and maximum 5,000-row application export bounds.
- [x] Added deterministic lot-inventory CSV columns and UTC RFC 3339 expiry timestamps.
- [x] Reused the shared formula-safe CSV serializer.
- [x] Added deterministic browser download filename and CSV content type.
- [x] Registered the new lot-inventory export route without changing the existing JSON endpoint.
- [x] Added focused tests for bounds, date parsing, CSV schema, formula safety, timestamps, and invalid expiry dates.
- [x] Fixed the PostgreSQL `ListBalances` implementation required by the v0.1.4 inventory export contract.
- [x] Added `docs/RELEASE_NOTES_v0.1.5.md`.
- [x] Added the v0.1.5 `CHANGELOG.md` entry.

## v0.1.5 release gates

- [ ] Run `gofmt`, `go vet ./...`, normal tests, and race-enabled Go tests.
- [ ] Verify `repository.Inventory` is implemented by every production and test store.
- [ ] Verify lot export pagination and the 5,000-row application bound.
- [ ] Verify product, warehouse, location, lot, and `expiringBy` filters.
- [ ] Verify invalid expiry dates return HTTP 400.
- [ ] Verify deterministic schema, ordering, formula safety, UTC timestamps, filename, and content type.
- [ ] Run PostgreSQL migration/readiness smoke testing.
- [ ] Run web typecheck and production build.
- [ ] Run Android lint/tests/build and release-networking/security checks.
- [ ] Run browser-companion manifest and unit checks.
- [ ] Run configured CodeQL checks for Go and JavaScript/TypeScript.
- [ ] Complete deployed authorization/CSRF and authentication/session smoke checks.
- [ ] Verify no secrets or credentials are included in exported data or artifacts.
- [ ] Create immutable `v0.1.5` tag on the verified commit.
- [ ] Publish the GitHub Release with the prepared notes and mark it pre-release until every gate passes.
- [ ] Perform post-release smoke testing against the published tag/artifacts.

## Publication details

- Version: `v0.1.5`
- Release title: `StockPilot v0.1.5 — Lot Inventory & Expiry-Risk CSV Export`
- Git tag: `v0.1.5`
- Release class: normal pre-1.0 feature release
- Pre-release: yes until every release gate passes
- Release notes: `docs/RELEASE_NOTES_v0.1.5.md`

## Known verification limitation in this workspace

- The repository was updated directly through the connected GitHub integration.
- A local `go test ./...` cannot be treated as completed unless executed in an environment with the repository and required dependencies available.
- The release must remain pending until GitHub Actions or a local developer environment confirms the full build/test matrix.

## Next exact development tasks after v0.1.5

1. Add purchase-order and receiving CSV export contracts.
2. Add export-specific audit coverage for sensitive datasets.
3. Add deterministic streaming behavior for large exports.
4. Add web download controls and accessible export feedback.
5. Add integration coverage around authorization and export headers.
6. Add export job lifecycle/status endpoints for large datasets.
7. Add richer expiry-risk classification and operational alerts.
8. Continue toward v0.2.x analytics and operational reporting.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
