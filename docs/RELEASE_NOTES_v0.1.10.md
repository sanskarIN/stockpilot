# StockPilot v0.1.10 — Export Audit Trail

Release date: 2026-09-04

## Highlights

- Added an audit event for authenticated CSV export requests.
- Export audit events capture the authenticated actor and request ID already established by StockPilot's access layer.
- Export events identify the requested CSV route without storing query parameters or downloaded dataset contents.
- Added regression coverage for CSV-export recognition and audit-event identity.
- Preserved the v0.1.9 privacy hardening: CSV responses continue to use no-store/no-cache download headers.

## Audit behavior

After an authenticated request passes the existing permission check, GET requests whose path ends in `.csv` are recorded as:

- Action: `export.csv.requested`
- Entity type: `export`
- Entity ID: requested route path
- Actor ID: authenticated principal
- Request ID: existing request correlation ID when available
- Metadata: HTTP method only

The audit record deliberately excludes URL query parameters, credentials, session tokens, and exported row contents. This keeps the audit trail useful for accountability without duplicating sensitive export data.

## Compatibility

- No database migration is required.
- No API route was renamed or removed.
- Existing CSV schemas and download filenames remain unchanged.
- Existing RBAC permissions continue to control export access.
- No new third-party dependency was introduced.

## Verification

Before treating this release as fully verified:

1. Run `gofmt`.
2. Run `go vet ./...`.
3. Run the normal Go test suite.
4. Run race-enabled Go tests.
5. Verify authenticated CSV requests create `export.csv.requested` audit events.
6. Verify unauthenticated and unauthorized requests remain rejected before export handling.
7. Verify audit records contain actor and request IDs but no query parameters, credentials, or session secrets.
8. Verify all CSV exports retain `Cache-Control: no-store` and `Pragma: no-cache`.
9. Run PostgreSQL readiness and export smoke tests.
10. Run Web, Android, browser-companion, and CodeQL checks where configured.
11. Create the immutable `v0.1.10` tag on the verified commit.
12. Publish the GitHub Release.
13. Perform post-release export and audit smoke testing.

## Known limitations

- The audit event records the export request, not the exact number of rows successfully written to the client.
- Current exports remain bounded and materialize result sets before CSV serialization; cursor-based streaming is future work.

## Upgrade notes

No schema migration or data conversion is required for v0.1.10.
