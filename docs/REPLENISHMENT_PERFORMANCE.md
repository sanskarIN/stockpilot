# Replenishment Performance

The replenishment-performance report measures historical purchase-order execution by supplier.

## Endpoint

`GET /api/v1/reports/replenishment-performance`

The endpoint accepts the standard reporting bounds:

- `from` and `to` as RFC3339 timestamps;
- `limit` for page size;
- `offset` for page offset.

The response keeps StockPilot's established JSON report envelope and report metadata headers.

## Metrics

- `orderCount`: non-cancelled purchase orders in the reporting period.
- `orderedUnits`: sum of quantities ordered on those purchase-order lines.
- `receivedUnits`: sum of quantities recorded as received.
- `outstandingUnits`: ordered units that remain outstanding.
- `fillRate`: received units divided by ordered units.
- `onTimeOrderCount`: orders with a first receipt at or before the expected receipt time; orders without an expected time are counted as on-time when a receipt exists.
- `lateOrderCount`: orders with a first receipt after the expected receipt time.
- `averageLeadDays`: average elapsed days from purchase-order creation to first recorded receipt.

The report describes historical purchasing performance. It does not claim that a reorder recommendation caused a purchase order unless the application has an explicit causal linkage.

## Compatibility

The capability is additive. Existing `repository.Reports` consumers and existing report response bodies are unchanged. The endpoint requires a repository implementing `repository.ReplenishmentReports`.
