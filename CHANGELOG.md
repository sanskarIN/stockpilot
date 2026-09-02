# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased**.

## Unreleased

### Added

- Product, category, supplier, warehouse, location, lot, inventory movement, transfer, and purchase-order foundations.
- Session-based authentication, RBAC, CSRF protection, and administrator bootstrap tooling.
- PostgreSQL persistence with ordered migrations and transactional inventory operations.
- Responsive React + TypeScript web dashboard and installable PWA behavior.
- Native Android client with encrypted session storage, authenticated API access, dark mode, and release TLS enforcement.
- Manifest V3 browser companion with scoped optional host permissions and server health/launcher flow.
- CodeQL and Dependabot automation.
- Aggregate reorder recommendations, inventory valuation, exact barcode lookup, reporting coverage, catalog management, guided inventory operations, warehouse/location lifecycle, multi-line purchasing, lot/expiry receiving, append-only business auditability, lot inventory visibility, browser/Android barcode scanning, and companion scan handoff.
- Authentication/session audit-event definitions and regression coverage.

### Changed

- Sensitive successful business mutations and authentication/account lifecycle events now use the append-only audit stream with request-ID correlation.
- Authentication audit metadata is deliberately coarse for failed login/session events and excludes credential material.

### Security

- HTTP mutation requests retain explicit CSRF confirmation requirements.
- Authorization remains server-side and role-aware.
- Production Android networking requires TLS.
- Browser companion permissions remain scoped to the configured StockPilot origin.
- Audit history remains read-only from the web client and audit storage exposes append/list operations without an update/delete API.
- Authentication audit events never store passwords, raw session tokens, cookie values, or credential-bearing metadata.
- Scanner flows do not copy StockPilot session credentials into scanner state.

## Release discipline

Dependency upgrades and generated artifacts are reviewed independently from product changes. Stable releases require the checks in `docs/RELEASE_CHECKLIST.md` and the restore procedure in `docs/RESTORE_DRILL.md`.
