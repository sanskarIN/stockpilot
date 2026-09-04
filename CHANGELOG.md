# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased** until each release is published.

## Unreleased

### Next: v0.2.8 — Large-Report Scalability & Purchasing Trends

Preparation for v0.2.8 focuses on making reporting scale safely as StockPilot's datasets grow while adding purchasing trend visibility.

Planned scope:

- reusable bounded pagination primitives for large report datasets;
- additive cursor/keyset pagination or streaming without silently changing existing endpoint semantics;
- opaque continuation tokens with deterministic ordering;
- bounded memory use, page sizes, date windows, and query work;
- supplier/product purchasing trend series built from authoritative purchase-order and purchase-order-line data;
- ordered, received, open-unit, and purchasing-value metrics with currencies kept separate;
- safe date-window and row bounds for trend endpoints;
- CSV export that remains formula-safe and non-cacheable;
- Reports & Analytics trend views with clear empty/loading/error states and accessible table alternatives;
- regression coverage for pagination, ordering, bounds, continuation, duplicate prevention, and export behavior.

No feature in this preparation should silently change transactional inventory or purchasing records. A schema migration will be avoided unless implementation evidence proves the existing schema insufficient.

### v0.2.7 — Replenishment Readiness

Completed implementation scope:

- added advisory replenishment-readiness reporting derived from existing reorder suggestions and stock-movement history;
- added on-hand, reorder point, reorder quantity, target stock, and suggested quantity visibility;
- added recent outbound units and average daily outbound velocity;
- added estimated days of cover when outbound velocity is positive;
- added deterministic `out_of_stock`, `critical`, `reorder`, `watch`, and `healthy` classifications;
- added bounded authenticated JSON reporting and formula-safe CSV export;
- added the `GET /api/v1/reports/replenishment-readiness` endpoint;
- updated API metadata to `0.2.7`;
- added domain contract and regression coverage;
- kept the report read-only with no purchase-order creation and no inventory mutation;
- introduced no database migration.

Stable publication still requires the repository's applicable verification gates to be confirmed green.

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
