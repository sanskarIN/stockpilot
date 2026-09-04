# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased**.

## Unreleased

Future work will continue here after v0.1.7.

## 0.1.7 — 2026-09-04

### Added

- Added bounded audit-log CSV export at `GET /api/v1/audit/export.csv`.
- Added actor, action, entity-type, and entity-ID filters to the export contract.
- Added focused HTTP coverage for export bounds, headers, schema, timestamps, and formula-safe metadata.

### Changed

- Audit exports reuse the shared formula-safe CSV serializer introduced in v0.1.2.
- Export pagination defaults to 500 rows and is capped at 5,000 at the application layer.
- Negative offsets normalize to zero.
- Audit timestamps are normalized to UTC RFC3339 values.
- Browser downloads use the deterministic filename `stockpilot-audit-log.csv`.
- Audit metadata is serialized as compact JSON within its CSV cell.

### Security and operations

- Audit export is read-only and remains behind the existing HTTP middleware and application authorization/session controls.
- The export does not add credential, password, raw session-token, cookie, or payment fields.
- Existing audit-event creation remains responsible for preventing sensitive credential material from entering metadata.
- Export bounds prevent unbounded application-level result generation.

### Verification

- Added unit coverage for pagination normalization and the CSV contract.
- Full release verification remains required for Go, PostgreSQL, Web, Android, extension, authentication/session, authorization/CSRF, and CodeQL gates where configured.

## 0.1.6 — 2026-09-04

### Added

- Added bounded purchase-order CSV export at `GET /api/v1/orders/export.csv`.
- Added optional purchase-order status filtering for `draft`, `ordered`, `partially_received`, `received`, and `cancelled`.
- Added line-level ordered, received, remaining, unit-cost, and line-total fields to the export.
- Added a dedicated `PurchaseOrderExportRow` domain model and repository export operation.
- Added a deterministic PostgreSQL export query that selects purchase orders first and then flattens their lines.
- Added focused HTTP coverage for export bounds, status validation, headers, schema, formula-safe values, totals, and timestamps.

### Changed

- Purchase-order exports reuse the shared formula-safe CSV serializer introduced in v0.1.2.
- Purchase-order selection is deterministic by `created_at DESC, id DESC` and line ordering is deterministic by line ID.
- Export timestamps are normalized to UTC RFC 3339 values.
- Browser downloads use the deterministic filename `stockpilot-purchase-orders.csv`.
- Receiving visibility is exposed through the authoritative `received` and `remaining` line quantities already maintained by the existing receiving transaction flow.

### Security and operations

- Purchase-order export is read-only and remains behind the existing HTTP middleware and application authorization/session controls.
- Application-level export bounds cap requests at 5,000 selected purchase orders; repository bounds enforce the same ceiling.
- Export serialization protects spreadsheet formula-like values.
- Export schemas exclude credentials, passwords, session secrets, and payment information.

### Verification

- Added unit coverage for purchase-order export normalization and CSV contracts.
- Full release verification remains required for Go, PostgreSQL, Web, Android, extension, authentication/session, authorization/CSRF, and CodeQL gates where configured.

## 0.1.5 — 2026-09-03

### Added

- Added bounded lot-inventory CSV export at `GET /api/v1/inventory/lots/export.csv`.
- Added optional product, warehouse, location, and lot filters to the export contract.
- Added the inclusive `expiringBy` date filter using `YYYY-MM-DD` input.
- Added focused HTTP coverage for bounds, date parsing, headers, schema, timestamps, expiry filtering, and formula-safe values.

### Changed

- Lot inventory exports reuse the shared formula-safe CSV serializer introduced in v0.1.2.
- Export rows retain deterministic PostgreSQL ordering by expiry date, product name, lot number, and location name.
- Expiry timestamps are normalized to UTC RFC 3339 values.
- Browser downloads use the deterministic filename `stockpilot-lot-inventory.csv`.

### Security and operations

- Lot inventory export is read-only and remains behind the existing HTTP middleware and application authorization/session controls.
- Export requests are bounded at the application layer and the repository retains its own safety cap.
- Export schemas exclude credentials, passwords, session secrets, and payment information.

### Verification

- Added unit coverage for lot export normalization and CSV contracts.
- Full release verification remains required for Go, PostgreSQL, Web, Android, extension, authentication/session, authorization/CSRF, and CodeQL gates where configured.

## Release discipline

Dependency upgrades and generated artifacts are reviewed independently from product changes. Stable releases require the checks in `docs/RELEASE_CHECKLIST.md` and the restore procedure in `docs/RESTORE_DRILL.md`.
