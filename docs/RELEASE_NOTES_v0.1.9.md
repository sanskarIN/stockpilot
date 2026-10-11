# StockPilot v0.1.9 — Export Privacy & Authorization Hardening

Release date: 2026-09-04

## Highlights

- Centralized CSV download response headers for all StockPilot export endpoints.
- Added `Cache-Control: no-store` and `Pragma: no-cache` to prevent browser/proxy retention of exported operational datasets.
- Preserved deterministic CSV content type and attachment filenames through the shared header helper.
- Added focused authorization-contract coverage for every current CSV export route.
- Confirmed exports continue to resolve to their domain-specific read permissions rather than write permissions.
- Added regression coverage for the privacy-oriented response headers.

## Export security changes

All current CSV exports now use a shared response-header helper that sets:

- `Content-Type: text/csv; charset=utf-8`
- `Content-Disposition: attachment; filename="..."`
- `Cache-Control: no-store`
- `Pragma: no-cache`

The helper is applied to:

- Product catalog export.
- Inventory balance export.
- Low-stock export.
- Reorder-suggestion export.
- Lot-inventory export.
- Receipt-history export.
- Purchase-order export.
- Audit-log export.

## Authorization coverage

The export access contract is covered for:

- Catalog exports → `PermissionCatalogRead`.
- Inventory, lot, receipt, low-stock, and reorder exports → `PermissionInventoryRead`.
- Purchase-order exports → `PermissionOrdersRead`.
- Audit-log exports → `PermissionAuditRead`.

The existing `WithAccess` middleware remains responsible for resolving the authenticated principal and denying requests that lack the required permission.

## Compatibility

- No database migration is required.
- No API route was renamed or removed.
- CSV schemas are unchanged from their preceding releases.
- Existing authentication, CSRF, origin, request-ID, and security-header middleware remains in effect.
- No new third-party dependency was introduced.

## Verification gates

Before treating this release as fully verified:

1. Run `gofmt`.
2. Run `go vet ./...`.
3. Run the normal Go test suite.
4. Run race-enabled Go tests.
5. Verify every export route through the authenticated access middleware.
6. Verify unauthorized principals receive HTTP 403 for exports they cannot read.
7. Verify every export response includes `no-store` and `no-cache` download headers.
8. Verify deterministic content types and filenames remain unchanged.
9. Run PostgreSQL readiness and export smoke tests.
10. Run Web, Android, browser-companion, and CodeQL checks where configured.
11. Verify no exported schema introduces credentials or session secrets.
12. Create immutable `v0.1.9` tag on the verified commit.
13. Publish the GitHub Release with these notes.
14. Perform post-release smoke testing.

## Known limitations

- The repository-backed export APIs still materialize bounded result sets before CSV serialization.
- True cursor/streaming export remains future work for very large datasets.
- This release strengthens response privacy and permission-contract regression coverage; it does not introduce a new role or permission model.

## Upgrade notes

No schema migration or data conversion is required for v0.1.9.
