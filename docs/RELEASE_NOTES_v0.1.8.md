# StockPilot v0.1.8 — Receipt History & Export Reliability

Release date: 2026-09-04

## Highlights

- Added a dedicated receipt-history export backed by the authoritative `stock_movements` receiving records.
- Added product, warehouse, location, lot, actor, reference, and date-range filters.
- Added deterministic pagination bounds with a 500-row default and 5,000-row application maximum.
- Added stable ordering by receipt occurrence time and movement ID.
- Reused the shared formula-safe CSV serializer.
- Normalized receipt timestamps to UTC RFC3339 values.
- Added deterministic browser download headers and filename.
- Added focused HTTP coverage for pagination, date validation, filtering, headers, schema, formula safety, and timestamps.

## New API endpoint

`GET /api/v1/inventory/receipts/export.csv`

### Query parameters

- `productId` — optional product filter.
- `warehouseId` — optional warehouse filter.
- `locationId` — optional location filter.
- `lotId` — optional lot filter.
- `actorId` — optional receiving actor filter.
- `reference` — optional exact receipt reference filter.
- `from` — optional inclusive start date in `YYYY-MM-DD` format.
- `to` — optional exclusive end date in `YYYY-MM-DD` format.
- `limit` — default 500, maximum 5000.
- `offset` — default 0; negative values normalize to 0.

### CSV columns

`movementId`, `productId`, `sku`, `productName`, `locationId`, `location`, `warehouseId`, `warehouse`, `lotId`, `lotNumber`, `quantity`, `reference`, `note`, `actorId`, `occurredAt`, `createdAt`

## Export behavior

- Only authoritative `receive` stock-movement records are returned.
- Responses use `text/csv; charset=utf-8`.
- Browser downloads use `stockpilot-receipt-history.csv`.
- Formula-like cell values are protected by the shared CSV serializer.
- `occurredAt` and `createdAt` are serialized as UTC RFC3339 timestamps.
- Results are deterministically ordered by `occurred_at DESC, id DESC`.
- `from` is inclusive and `to` is exclusive, making adjacent date windows safe to compose.

## Security and privacy

- Export is read-only and does not create or mutate inventory records.
- Existing authentication, authorization, origin, CSRF, request-ID, and security-header controls remain applicable.
- Export bounds prevent unbounded application-level result requests.
- The export returns operational receipt data already stored by the inventory subsystem; callers should apply the same authorization boundaries used for inventory reporting and receiving operations.
- No passwords, raw session tokens, cookie values, or credentials are added to the export contract.

## Compatibility

- No database migration is required; the endpoint reads the existing `stock_movements`, products, locations, warehouses, and lots tables.
- Existing movement, receiving, inventory, lot, purchase-order, and audit endpoints remain unchanged.
- Existing CSV export endpoints remain available.
- No new third-party CSV dependency was introduced.

## Verification gates

Before treating this release as fully verified:
1. Run `gofmt`, `go vet ./...`, normal tests, and race-enabled Go tests.
2. Verify every production and test implementation of `repository.Inventory` implements `ListReceiptHistory`.
3. Verify every supported filter reaches the repository contract and SQL predicate.
4. Verify inclusive `from`, exclusive `to`, invalid date handling, and invalid ranges.
5. Verify default, negative, valid, and clamped pagination bounds.
6. Verify deterministic receipt ordering and CSV schema.
7. Verify formula safety and UTC timestamp serialization.
8. Verify privileged authorization and denial for unauthorized users.
9. Verify origin/CSRF/request-ID/security middleware behavior where applicable.
10. Run PostgreSQL readiness and receiving smoke tests with representative receipt data.
11. Run Web, Android, browser-companion, and CodeQL checks where configured.
12. Verify no credentials or secrets are present in receipt notes, references, or other exported fields.
13. Create immutable `v0.1.8` tag on the verified commit.
14. Publish the GitHub Release with the prepared notes and keep it pre-release until every gate passes.
15. Perform post-release smoke testing against the published tag/artifacts.

## Known limitations

- Export is bounded and paginated rather than asynchronous.
- The current repository API materializes the selected rows before CSV serialization; true streaming for very large datasets remains a follow-up.
- Receipt history is derived from `stock_movements` with `movement_type = 'receive'`; it intentionally does not synthesize history from current purchase-order `received` counters.

## Upgrade notes

No schema migration or data conversion is required for v0.1.8.
