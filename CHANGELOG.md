# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased**.

## Unreleased

Future work will continue here after v0.2.4.

## 0.2.4 — 2026-09-04

### Added

- Added stock movement history aggregation for product/location/lot activity.
- Added configurable movement-history windows from 1 to 365 days, defaulting to 30 days.
- Added movement count, inbound units, outbound units, net units, and average daily outbound metrics.
- Added bounded JSON reporting and formula-safe CSV export.
- Added a Reports & Analytics movement-velocity panel and export action.

### Changed

- Reports & Analytics now includes recent movement velocity alongside valuation and inventory aging.
- Movement reporting uses server-side aggregation from authoritative `stock_movements` records.
- Deterministic ordering prioritizes outbound activity and recent movement timestamps.

### Security and reliability

- Movement reports remain read-only and use existing authenticated reporting access controls.
- CSV exports retain bounded limits and `no-store` / `no-cache` download headers.
- Formula-safe serialization is retained for spreadsheet-facing exports.

### Verification

- Stable publication requires the configured Go, PostgreSQL, Web, CodeQL, Android/browser companion, E2E/accessibility, restore/rollback, and artifact gates to pass.

## 0.2.3 — 2026-09-04

### Fixed

- Corrected audit-export metadata sanitization so formula-like JSON string values are neutralized before CSV delivery.
- Corrected receipt-history export regression coverage to assert the canonical UTC timestamp representation emitted by the exporter.

### Security and reliability

- Preserved formula-safe CSV serialization for authenticated exports.
- Preserved bounded export limits and privacy-oriented download headers.
- Kept audit and receipt-history exports read-only and behind the existing authorization layer.

## 0.2.2 — 2026-09-04

### Added

- Added expiry-risk reporting for positive inventory with persisted lot expiry dates.
- Added configurable expiry-risk windows with expired, critical, warning, and safe classifications.
- Added bounded expiry-risk reporting through the Reports & Analytics workflow.
- Added formula-safe CSV export for expiry-risk results.
- Added regression coverage for expiry classification boundaries, authorization, export limits, and CSV safety.

### Changed

- Reports & Analytics now exposes expiry risk alongside valuation and inventory aging.
- Expiry classification uses server-side report semantics and deterministic ordering.
- Expiry reporting remains read-only and respects existing authentication and reporting permissions.

### Security and operations

- Expiry-risk exports use bounded result generation and privacy-oriented no-store/no-cache response headers.
- CSV values remain protected against spreadsheet formula injection.
- No credentials, session tokens, passwords, or payment data are introduced by the report.

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

## Release discipline

Dependency upgrades and generated artifacts are reviewed independently from product changes. Stable releases require the checks in `docs/RELEASE_CHECKLIST.md` and the restore procedure in `docs/RESTORE_DRILL.md`.
