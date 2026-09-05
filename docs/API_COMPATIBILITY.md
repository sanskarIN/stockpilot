# API Compatibility Policy

StockPilot is pre-1.0, but API changes should still be deliberate and reviewable.

## Compatibility rules

- Existing JSON field names should not be renamed without a migration/deprecation plan.
- New response fields and headers should be additive whenever possible.
- New request fields should have safe defaults or be explicitly optional.
- Permission changes must be treated as security-sensitive changes.
- Database migrations must preserve compatibility with the server version deployed during a rolling update when practical.
- Breaking API changes should be introduced under an explicit API-versioning plan rather than silently changing `/api/v1` semantics.

## Reporting pagination

The supplier-performance and warehouse-valuation report endpoints support bounded pagination with `limit` and `offset`.

- `X-Report-Limit` and `X-Report-Offset` describe the applied page bounds.
- `X-Report-Complete` is conservative and does not imply completeness from a full page alone.
- `X-Total-Count` is an additive response header emitted when the configured repository implements the optional counted-report capability.
- `X-Total-Count` represents the complete number of report rows, independent of the current page offset.
- Existing JSON and CSV response bodies do not require clients to consume the new header.
- Cursor/streaming pagination is intentionally not part of the current `/api/v1` contract.

Approved browser origins can read reporting pagination headers through the API CORS exposure policy.

## HTTP behavior

Use stable status classes for clients:

- `2xx` for successful operations.
- `400` for malformed or invalid client input.
- `401` for missing/expired authentication.
- `403` for authenticated callers without permission or required CSRF confirmation.
- `404` for resources that do not exist or are intentionally not disclosed.
- `409` for business-rule conflicts.
- `422` for semantically invalid input when the endpoint already follows that convention.
- `5xx` only for server-side failures.

Error responses should remain machine-readable and include the request ID when available.

## Change checklist

Before changing an endpoint, check server handlers, repository contracts, domain validation, web client types, Android models, browser-companion behavior, tests, and documentation for the affected contract.
