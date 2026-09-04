# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased**.

## Unreleased

Future work will continue here after v0.1.9.

## 0.1.9 — 2026-09-04

### Added

- Added a shared CSV download-header helper for all export endpoints.
- Added authorization-contract regression coverage for every current CSV export route.
- Added focused tests for privacy-oriented CSV response headers.

### Changed

- All CSV exports now send `Cache-Control: no-store` and `Pragma: no-cache`.
- CSV content type and deterministic attachment filenames are centralized through the shared export-header helper.
- Export permission mapping is explicitly covered for catalog, inventory, lot, receipt, reorder, purchasing, and audit datasets.

### Security and operations

- Export responses are instructed not to be retained by browser or intermediary caches by default.
- Export access continues to use the authenticated `WithAccess` middleware and domain-specific read permissions.
- The release does not introduce a new role or permission model and does not broaden export access.

### Verification

- Added regression coverage for all current export permission mappings and download privacy headers.
- Full release verification remains required for Go, PostgreSQL, Web, Android, extension, authentication/session, authorization/CSRF, and CodeQL gates where configured.

## 0.1.8 — 2026-09-04

### Added

- Added dedicated receipt-history CSV export at `GET /api/v1/inventory/receipts/export.csv`.
- Added product, warehouse, location, lot, actor, reference, and date-range filters.
- Added a dedicated `ReceiptHistoryRow` domain model and `ReceiptHistoryFilter` repository contract.
- Added a PostgreSQL query backed by authoritative `stock_movements` records with `movement_type = 'receive'`.
- Added focused HTTP coverage for pagination, date validation, filtering, headers, schema, formula safety, and timestamps.

### Changed

- Receipt history exports reuse the shared formula-safe CSV serializer introduced in v0.1.2.
- Export pagination defaults to 500 rows and is capped at 5,000 at the application layer.
- Negative offsets normalize to zero.
- Receipt rows are deterministically ordered by `occurred_at DESC, id DESC`.
- Receipt timestamps are normalized to UTC RFC3339 values.
- Browser downloads use the deterministic filename `stockpilot-receipt-history.csv`.
- Date windows use an inclusive `from` boundary and an exclusive `to` boundary.

### Security and operations

- Receipt history export is read-only and remains behind the existing HTTP middleware and application authorization/session controls.
- Export bounds prevent unbounded application-level result generation.
- The export does not add credential, password, raw session-token, cookie, or payment fields.
- Receipt history is sourced from persisted receiving events rather than reconstructed from current purchase-order received counters.

### Verification

- Added unit coverage for receipt-history pagination normalization, date-range validation, filtering, CSV headers, timestamps, and formula-safe values.
- Full release verification remains required for Go, PostgreSQL, Web, Android, extension, authentication/session, authorization/CSRF, and CodeQL gates where configured.

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

## Release discipline

Dependency upgrades and generated artifacts are reviewed independently from product changes. Stable releases require the checks in `docs/RELEASE_CHECKLIST.md` and the restore procedure in `docs/RESTORE_DRILL.md`.
