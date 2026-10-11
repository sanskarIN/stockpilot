# StockPilot v0.1.6 — Purchasing & Receiving CSV Export

Release date: 2026-09-04

## Highlights

- Added a bounded purchase-order CSV export.
- Added purchase-order status filtering to the export contract.
- Added line-level ordered, received, remaining, unit-cost, and line-total fields.
- Reused the shared formula-safe CSV serializer.
- Normalized export timestamps to UTC RFC 3339 form.
- Added deterministic browser download headers and filename.
- Added focused HTTP tests for bounds, status validation, CSV schema, formula-safe values, totals, and timestamps.
- Added a dedicated repository export query so purchase-order export pagination is handled without loading full order graphs into the HTTP layer.

## New API endpoint

### Purchase orders

`GET /api/v1/orders/export.csv`

Query parameters:

- `status` — optional purchase-order status: `draft`, `ordered`, `partially_received`, `received`, or `cancelled`.
- `limit` — default 500, maximum 5000 purchase orders selected for export.
- `offset` — default 0; negative values are normalized to 0.

The CSV is line-oriented: each purchase-order line is one CSV row. A purchase order with multiple lines therefore produces multiple rows.

Columns:

`orderId`, `orderNumber`, `supplierId`, `warehouseId`, `status`, `currency`, `expectedAt`, `notes`, `createdBy`, `createdAt`, `updatedAt`, `lineId`, `productId`, `quantity`, `received`, `remaining`, `unitCostMinor`, `lineTotalMinor`

## Receiving visibility

The export exposes the current receiving progress already stored on each purchase-order line:

- `quantity` — ordered quantity.
- `received` — quantity received so far.
- `remaining` — calculated as ordered quantity minus received quantity.
- `status` — current purchase-order lifecycle state.

Receiving itself continues to use the existing transactional endpoint:

`POST /api/v1/orders/{orderID}/lines/{lineID}/receive`

This release does not introduce a second receipt ledger; it exports the authoritative purchase-order receiving state already maintained by the existing transaction flow.

## Export behavior

- Responses use `text/csv; charset=utf-8`.
- Browser downloads receive the deterministic filename `stockpilot-purchase-orders.csv`.
- Spreadsheet formula-like values are protected by the shared CSV writer.
- All exported timestamps are normalized to UTC RFC 3339 form.
- Purchase orders are selected deterministically by `created_at DESC, id DESC`.
- Lines within each selected purchase order are ordered deterministically by line ID.
- The HTTP layer validates unsupported status values before querying the repository.

## Security and operational boundaries

- The export is read-only.
- Existing middleware remains responsible for request IDs, origin controls, security headers, authentication/session integration, and CSRF policy where applicable.
- Application and repository bounds prevent unbounded export requests.
- The export contains operational purchasing data but no passwords, session secrets, credentials, or payment information.
- `createdBy` is included because it is part of the existing purchase-order audit context; downstream consumers should treat it as operational data.

## Compatibility

- No database migration is required.
- Existing product and inventory CSV exports remain available.
- Existing purchase-order create, update, status, and receiving APIs remain unchanged.
- No new third-party CSV dependency was introduced.
- XLSX, PDF, asynchronous bulk jobs, object-storage exports, and dedicated receipt-ledger exports remain outside this release.

## Verification gates

Before treating the GitHub release as fully verified:

1. Run `go test ./...`.
2. Run `go vet ./...`.
3. Run race-enabled Go tests.
4. Verify the purchase-order CSV header and line rows.
5. Verify every supported `status` filter.
6. Verify unsupported status values return HTTP 400.
7. Verify negative offsets normalize to zero.
8. Verify limits above 5000 are clamped at the application boundary.
9. Verify formula-safe serialization for order numbers and other spreadsheet-facing values.
10. Verify remaining and line-total calculations.
11. Verify UTC timestamp formatting.
12. Verify the deterministic download filename and content type.
13. Run PostgreSQL migration/readiness checks.
14. Run Web, Android, browser-companion, and CodeQL checks where configured.
15. Perform authentication, authorization, CSRF, and production smoke testing.

## Known limitations

- The export is bounded and paginated rather than an asynchronous bulk-export job.
- The limit applies to selected purchase orders; each selected order can produce multiple line rows.
- Receiving history is represented by current `received` quantities rather than a separate historical receipt-event export.
- Audit-log CSV export remains future work.

## Upgrade notes

No schema migration or data conversion is required for v0.1.6.
