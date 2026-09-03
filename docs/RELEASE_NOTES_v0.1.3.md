# StockPilot v0.1.3 — Product Catalog CSV Export

## Release summary

v0.1.3 moves the CSV export foundation from v0.1.2 into its first application-level operational export: a bounded, read-only product catalog download.

## Added

- `GET /api/v1/products/export.csv` for product catalog export.
- Query filters aligned with the existing product listing contract: `q`, `categoryId`, `supplierId`, `activeOnly`, `limit`, and `offset`.
- A hard server-side maximum of 5,000 exported rows per request.
- Deterministic CSV column order and UTC timestamps.
- Formula-safe serialization for user-downloadable catalog data.
- `Content-Type` and `Content-Disposition` response headers suitable for browser downloads.

## Security and reliability

- Export remains read-only and uses the existing catalog repository rather than a second data-access path.
- Export size is bounded to reduce memory and database pressure.
- Formula-like cell values are protected by the shared CSV serialization package.
- No authentication credentials, sessions, cookies, CSRF secrets, or other secret material are included in the export schema.
- Existing origin, request-ID, security-header, and server-side authorization boundaries remain in force at the surrounding application layer.

## Compatibility

- No database migration is required.
- Existing JSON product endpoints are unchanged.
- Existing product filters are reused rather than introducing a parallel filter vocabulary.
- CSV remains the first portable export format; XLSX/PDF are still out of scope.

## Verification

Before publication, run the complete release checklist plus focused verification of:

- catalog export route registration;
- query/filter propagation;
- 5,000-row export bound;
- CSV header and field ordering;
- formula-safe output;
- content-disposition filename;
- error handling when catalog persistence fails;
- authorization behavior through the deployed application stack.

## Known limitations

v0.1.3 exposes the first operational export only. Inventory, reorder, lot/expiry, purchasing, and audit exports remain planned follow-up work.

Large multi-dataset exports and asynchronous export jobs are not part of this release.
