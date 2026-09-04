# StockPilot v0.2.6 — Warehouse & Location Valuation

**Version:** `0.2.6`  
**Tag:** `v0.2.6`  
**Release type:** Stable after all required verification gates pass  
**Release date:** 2026-09-04

## Highlights

StockPilot v0.2.6 extends inventory valuation from product-level totals to an operational warehouse/location breakdown. This makes on-hand capital easier to trace to the physical storage structure without introducing currency conversion assumptions.

### Warehouse and location valuation report

New endpoint:

`GET /api/v1/reports/warehouse-valuation`

Supported query parameters:

- `limit`: maximum location rows from 1 to 5000; default 1000.
- `format=csv`: returns a bounded CSV export.

Each breakdown row includes:

- warehouse identity and code;
- location identity and code;
- currency;
- on-hand units;
- valuation in minor currency units;
- number of distinct active products represented.

The response also includes warehouse-level totals grouped by currency.

### Currency safety

Valuation remains grouped by the product's stored currency. StockPilot does not silently convert currencies, so mixed-currency inventory cannot produce a misleading single total.

### CSV safety

The export is read-only and uses the existing formula-safe CSV cell serialization and privacy-oriented download headers.

## Reliability and security

- No database migration is required.
- The report uses authoritative inventory balances, products, locations, and warehouses.
- Only positive balances for active products contribute to valuation.
- Result generation is bounded to protect the reporting endpoint.
- CSV exports use `no-store` / `no-cache` headers.
- Existing authentication, authorization, CORS, CSRF, and security-header behavior remains unchanged.

## Upgrade notes

No schema migration is required for v0.2.6. Existing inventory balances, product costs/currencies, locations, and warehouses are sufficient.

## Verification gates

Before declaring the release stable, verify:

- Go formatting, vet, race tests, unit tests, and production build;
- PostgreSQL migration and integration tests;
- Web lint/type-check/build;
- CodeQL;
- Android client checks;
- browser companion checks;
- end-to-end and accessibility checks;
- restore/rollback checks;
- release artifact checks.

Do not mark the GitHub release as the latest stable release until the required gates are green.
