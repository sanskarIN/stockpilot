# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased**.

## Unreleased

Future work will continue here after v0.2.0.

## 0.2.0 — 2026-09-04

### Added

- Added a dedicated Reports & Analytics workspace in the web application.
- Added typed report-overview, inventory-summary, purchasing-summary, and valuation API clients.
- Added an inventory valuation CSV export with formula-safe serialization and bounded row counts.
- Added a direct dashboard entry point to the Reports workspace.
- Added regression coverage for valuation-export bounds, headers, privacy policy, and formula-safe cells.

### Changed

- Reporting navigation is now an actual application workflow rather than a dashboard placeholder link.
- Valuation exports reuse the same privacy-oriented `no-store`/`no-cache` download policy as other CSV exports.
- Reports remain read-only and continue to use server-side reporting permissions.

### Security and operations

- Valuation CSV downloads remain behind authenticated reporting authorization.
- Export audit instrumentation introduced in v0.1.10 automatically records authenticated CSV requests without storing query strings or dataset contents.
- Formula-safe CSV serialization is enabled for the new valuation export.

### Verification

- Added focused Go tests for the valuation export contract.
- Repository-local execution is not available in the connected workspace; GitHub Actions or a local checkout remains authoritative for full Go, PostgreSQL, Web, Android, extension, and CodeQL verification.

## 0.1.10 — 2026-09-04

### Added

- Added an `export.csv.requested` audit event for authenticated CSV export requests.
- Added regression coverage for CSV export request recognition and audit-event identity.
- Added release documentation for the export audit trail.

### Changed

- The access layer now records the authenticated actor and existing request ID for GET requests ending in `.csv` after the normal permission check passes.
- Export audit metadata is intentionally limited to the HTTP method; query parameters and dataset contents are not duplicated into the audit log.

### Security and operations

- CSV export requests now have an explicit accountability trail while retaining the v0.1.9 `no-store`/`no-cache` response policy.
- Export audit events do not store passwords, session tokens, query strings, or downloaded rows.
- Existing domain-specific read permissions and authenticated session controls remain unchanged.

### Verification

- Added regression coverage for export recognition and actor/request correlation.
- Full release verification remains required for Go, PostgreSQL, Web, Android, extension, authentication/session, authorization/CSRF, and CodeQL gates where configured.

## 0.1.9 — 2026-09-04

### Added

- Added a shared CSV download-header helper for all export endpoints.
- Added authorization-contract regression coverage for every current CSV export route.
- Added focused tests for privacy-oriented CSV response headers.
