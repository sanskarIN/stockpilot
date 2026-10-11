# StockPilot v0.1.5 — Lot Inventory & Expiry-Risk CSV Export

Release date: 2026-09-03

## Highlights

- Added a bounded lot-inventory CSV export.
- Added product, warehouse, location, lot, and expiry-date filtering to the export contract.
- Added an expiry-risk filter through `expiringBy`.
- Reused the shared formula-safe CSV serializer.
- Preserved deterministic ordering through the existing PostgreSQL lot-inventory query.
- Added focused HTTP tests for bounds, date parsing, headers, schema, timestamps, expiry filtering, and formula-safe values.

## New API endpoint

### Lot inventory

`GET /api/v1/inventory/lots/export.csv`

Query parameters:

- `productId` — optional product filter.
- `warehouseId` — optional warehouse filter.
- `locationId` — optional location filter.
- `lotId` — optional lot filter.
- `expiringBy` — optional inclusive expiry date in `YYYY-MM-DD` format.
- `limit` — default 500, maximum 5000.
- `offset` — default 0; negative values are normalized to 0.

Columns:

`productId`, `sku`, `productName`, `lotId`, `lotNumber`, `locationId`, `location`, `warehouseId`, `warehouse`, `onHand`, `expiresAt`, `active`

## Export behavior

- Responses use `text/csv; charset=utf-8`.
- Browser downloads receive the deterministic filename `stockpilot-lot-inventory.csv`.
- Spreadsheet formula-like values are protected by the shared CSV writer.
- Expiry timestamps are normalized to UTC RFC 3339 form.
- The existing PostgreSQL query orders rows by expiry date, product name, lot number, and location name.
- Existing JSON lot-inventory behavior remains unchanged.

## Security and operational boundaries

- The export is read-only.
- Existing middleware remains responsible for request IDs, origin controls, security headers, authentication/session integration, and CSRF policy where applicable.
- Application-level export bounds prevent oversized requests; the repository retains its own safety cap.
- No credentials, passwords, session secrets, or payment information are included in the export schema.

## Compatibility

- No database migration is required.
- Existing v0.1.3 product and v0.1.4 inventory/reorder CSV exports remain available.
- No new third-party CSV dependency was introduced.
- XLSX, PDF, asynchronous jobs, and object-storage-backed exports remain outside this release.

## Verification gates

Before publishing the GitHub release, verify:

1. `go test ./...` passes.
2. Lot inventory export returns the expected header and rows.
3. Negative offsets normalize safely.
4. Requests above 5000 rows are clamped.
5. `expiringBy` accepts only `YYYY-MM-DD`.
6. Product, warehouse, location, and lot filters reach the repository contract.
7. Formula-safe serialization remains enabled.
8. UTC timestamp formatting is stable.
9. Download filename and content type are correct.
10. Existing authentication, origin, CSRF, and security-header tests remain green.
11. Android, extension, CodeQL, and full repository CI gates pass where configured.

## Known limitations

- The export is bounded and paginated rather than an asynchronous bulk-export job.
- Expiry risk is currently represented by the `expiringBy` filter; richer risk classification and alerting remain future work.
- Purchase-order, receiving, and audit CSV exports remain planned.

## Upgrade notes

No schema migration or data conversion is required for v0.1.5.
