# StockPilot v0.1.4 — Inventory & Reorder CSV Exports

Release date: 2026-09-03

## Highlights

- Added a bounded inventory-balance CSV export.
- Added a low-stock CSV export.
- Added a reorder-suggestions CSV export.
- Reused the shared formula-safe CSV serializer introduced in v0.1.2.
- Added a repository-level bounded inventory balance listing for deterministic pagination.
- Added focused HTTP tests for export bounds, headers, CSV columns, timestamps, and formula-safe values.

## New API endpoints

### Inventory balances

`GET /api/v1/inventory/export.csv`

Query parameters:

- `limit` — default 1000, maximum 5000.
- `offset` — default 0; negative values are normalized to 0.

Columns:

`productId`, `locationId`, `lotId`, `quantity`, `updatedAt`

### Low stock

`GET /api/v1/inventory/low-stock/export.csv`

Query parameters:

- `limit` — bounded before reaching the repository.

Columns:

`productId`, `locationId`, `lotId`, `quantity`, `updatedAt`

### Reorder suggestions

`GET /api/v1/inventory/reorder-suggestions/export.csv`

Query parameters:

- `limit` — bounded before reaching the repository.

Columns:

`productId`, `sku`, `name`, `supplierId`, `unit`, `onHand`, `reorderPoint`, `reorderQuantity`, `targetStock`, `suggestedQuantity`

## Export behavior

- Responses use `text/csv; charset=utf-8`.
- Browser downloads receive deterministic filenames.
- Spreadsheet formula-like values are protected by the shared CSV writer.
- Export timestamps are normalized to UTC RFC 3339 form.
- Inventory balance exports use stable product/location/lot ordering in the PostgreSQL repository.
- Existing JSON endpoints remain unchanged.

## Security and operational boundaries

- The exports are read-only.
- Existing HTTP middleware remains responsible for request IDs, origin controls, security headers, authentication/session integration, and CSRF policy where applicable.
- Inventory balance export has a hard application limit of 5000 rows per request.
- Low-stock and reorder suggestion repository queries retain their existing repository-side safety bounds.
- No credentials, passwords, session secrets, or payment information are included in export schemas.

## Compatibility

- No database migration is required.
- Existing catalog CSV export from v0.1.3 remains available.
- No new third-party CSV dependency was introduced.
- XLSX, PDF, asynchronous jobs, and object-storage-backed exports remain outside this release.

## Verification gates

Before publishing the GitHub release, verify:

1. `go test ./...` passes.
2. Inventory export returns the expected header and rows.
3. Negative offsets normalize safely.
4. Inventory export clamps requests above 5000 rows.
5. Low-stock and reorder exports return deterministic schemas.
6. Formula-safe serialization remains enabled.
7. UTC timestamp formatting is stable.
8. Download filenames and content types are correct.
9. Existing authentication, origin, CSRF, and security-header tests remain green.
10. Android, extension, CodeQL, and full repository CI gates pass where configured.

## Known limitations

- Inventory balance export is paginated and bounded; it is not an asynchronous bulk-export job.
- Low-stock and reorder suggestion exports inherit the existing repository-side query caps.
- Lot inventory, expiry-risk, purchasing, receiving, and audit CSV exports are still planned for subsequent milestones.
- The changelog file may require a separate conflict-safe update if another commit has changed it concurrently.

## Upgrade notes

No schema migration or data conversion is required for v0.1.4.
