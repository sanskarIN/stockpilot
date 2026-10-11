# StockPilot v0.1.7 — Audit Log CSV Export

Release date: 2026-09-04

## Highlights

- Added a bounded audit-log CSV export.
- Reused existing audit filters for actor, action, entity type, and entity ID.
- Added deterministic export bounds with a 500-row default and 5,000-row application maximum.
- Reused the shared formula-safe CSV serializer.
- Normalized audit timestamps to UTC RFC3339 values.
- Added deterministic browser download headers and filename.
- Added focused HTTP coverage for bounds, schema, headers, timestamps, and formula-safe metadata.

## New API endpoint

`GET /api/v1/audit/export.csv`

### Query parameters

- `actorId` — optional actor filter.
- `action` — optional audit action filter.
- `entityType` — optional entity type filter.
- `entityId` — optional entity ID filter.
- `limit` — default 500, maximum 5000.
- `offset` — default 0; negative values normalize to 0.

### CSV columns

`id`, `occurredAt`, `actorId`, `action`, `entityType`, `entityId`, `requestId`, `metadata`

## Export behavior

- Responses use `text/csv; charset=utf-8`.
- Browser downloads use `stockpilot-audit-log.csv`.
- Formula-like values are protected by the shared CSV serializer.
- Audit timestamps are normalized to UTC RFC3339.
- Metadata is serialized as compact JSON inside the CSV cell.
- Existing JSON audit listing remains unchanged.

## Security and privacy

- Export is read-only.
- Existing authentication, authorization, origin, CSRF, request-ID, and security-header controls remain applicable.
- Export bounds prevent unbounded application-level result requests.
- Audit metadata is exported as stored; callers should use existing authorization boundaries and operational policies when exposing audit data.
- The export contract does not add passwords, raw session tokens, cookie values, or credentials.

## Compatibility

- No database migration is required.
- Existing audit append/list behavior remains unchanged.
- Existing product, inventory, lot, and purchase-order exports remain available.
- No new third-party CSV dependency was introduced.

## Verification gates

Before treating this release as fully verified:

1. Run `gofmt`, `go vet ./...`, normal tests, and race-enabled Go tests.
2. Verify every audit filter reaches the repository contract.
3. Verify default, negative, valid, and clamped pagination bounds.
4. Verify deterministic CSV headers and filename.
5. Verify UTC timestamp serialization.
6. Verify formula-safe metadata serialization.
7. Verify audit export authorization for privileged and non-privileged users.
8. Verify CSRF/origin/security middleware behavior where applicable.
9. Run PostgreSQL readiness/migration smoke testing.
10. Run Web, Android, browser-companion, and CodeQL checks where configured.
11. Verify no credentials or session secrets are present in exported metadata.

## Known limitations

- Export is bounded and paginated rather than asynchronous.
- The export exposes stored audit metadata and does not independently redact arbitrary application metadata; sensitive metadata must remain excluded at event creation time.
- Dedicated asynchronous export jobs and richer export observability remain future work.

## Upgrade notes

No schema migration or data conversion is required for v0.1.7.
