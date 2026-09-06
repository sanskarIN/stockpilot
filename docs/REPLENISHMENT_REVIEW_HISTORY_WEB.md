# Replenishment Review History Workspace

StockPilot v0.4.6 adds a read-only review-history panel to **Reports & Analytics**.

## Purpose

The workspace makes the v0.4.4 replenishment traceability records visible to authorized web users without adding browser-side mutation operations.

Each row is an immutable review snapshot returned by:

`GET /api/v1/replenishment/reviews`

## Displayed fields

- review timestamp;
- product name and SKU snapshot;
- review outcome;
- on-hand quantity at review time;
- suggested quantity at review time;
- unit;
- linked purchase-order ID when present;
- reviewer identity.

## Filtering

The panel supports:

- All outcomes;
- `accepted`;
- `modified`;
- `dismissed`;
- `expired`.

The selected outcome is sent as the API `outcome` query parameter. Results are bounded to the first 50 records in the current workspace; the server remains authoritative for validation and pagination bounds.

## Safety

The workspace is intentionally read-only. It does not create, edit, delete, approve, submit, cancel, or receive purchase orders. Existing purchase-order lifecycle rules remain authoritative.

Authentication failures are handled consistently with the other Reports workspace panels and return the user to the session-expiry flow.

## Compatibility

No database migration is introduced by v0.4.6. The workspace consumes the existing v0.4.4 review records and endpoint.
