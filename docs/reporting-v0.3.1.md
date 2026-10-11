# Reporting v0.3.1 Integration Contract

StockPilot reporting endpoints may include additive execution metadata based on the shared `internal/reporting` primitives.

## Metadata

`generatedAt` identifies when the report response was assembled. `from` and `to` identify the applied inclusive period. `limit` and `offset` identify the applied result bounds. `complete` indicates whether the response represents a complete result set for the requested bounded operation.

## Compatibility

Metadata is additive and must not change the meaning of existing report fields. Existing consumers that ignore unknown JSON fields remain compatible.

## Bounds

Server-side bounds remain authoritative. Invalid limits, negative offsets, and unsupported reporting windows should return actionable public validation errors without database details.

## Cancellation

Report handlers should pass request contexts through to repository operations so client cancellation and server deadlines stop unnecessary work.

## Export behavior

CSV exports remain bounded, formula-safe, and non-cacheable. They should not disclose internal query text, credentials, session values, or database errors.
