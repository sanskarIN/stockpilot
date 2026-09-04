# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased**.

## Unreleased

Future work will continue here after v0.2.6.

## 0.2.6 — 2026-09-04

### Added

- Added warehouse/location inventory valuation reporting.
- Added location-level on-hand units, valuation, currency, and distinct-product counts.
- Added warehouse-level valuation totals grouped by currency.
- Added bounded JSON reporting and formula-safe CSV export.

### Changed

- Reports & Analytics can now trace inventory valuation from products into physical warehouse locations.
- Valuation remains currency-separated and does not perform implicit FX conversion.
- Warehouse valuation uses authoritative positive inventory balances and active product costs.

### Security and reliability

- Warehouse valuation is read-only and uses existing authenticated reporting access controls.
- Result rows are bounded to 1–5000 with a safe default.
- CSV exports use no-store/no-cache headers and formula-safe text serialization.
- No database migration is required for this release.

### Verification

- Stable publication requires the configured Go, PostgreSQL, Web, CodeQL, Android/browser companion, E2E/accessibility, restore/rollback, and artifact gates to pass.

## 0.2.5 — 2026-09-04

### Added

- Added supplier performance reporting for purchasing activity.
- Added supplier-level order, ordered-unit, received-unit, open-unit, and purchasing-value metrics.
- Added observed receipt lead-time measurement from purchase-order creation to the first recorded receive movement.
- Added completed-order and on-time-order counts using the purchase order expected date when available.
- Added bounded JSON reporting and formula-safe CSV export for supplier performance.
- Added a Reports & Analytics supplier reliability panel.

### Changed

- Reports & Analytics now combines supplier reliability with purchasing pipeline, valuation, inventory aging, and movement velocity.
- Supplier reporting uses server-side aggregation from authoritative purchase-order, purchase-order-line, supplier, and stock-movement records.
- Supplier results use deterministic ordering by order activity, supplier name, and supplier ID.

### Security and reliability

- Supplier reporting is read-only and uses the existing authenticated reporting access controls.
- Report windows are bounded to 1–365 days and result sets to 1–5000 rows, with safe defaults.
- CSV exports use no-store/no-cache headers and formula-safe text serialization.
- No database migration is required for this release.
