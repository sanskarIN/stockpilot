# StockPilot v0.2.5 — Supplier Performance

**Version:** `0.2.5`  
**Tag:** `v0.2.5`  
**Release type:** Stable after all required verification gates pass  
**Release date:** 2026-09-04

## Highlights

StockPilot v0.2.5 adds supplier purchasing performance reporting to make procurement reliability measurable from existing transactional data.

### Supplier performance report

New endpoint:

`GET /api/v1/reports/supplier-performance`

Supported query parameters:

- `days`: reporting window from 1 to 365 days; default 30.
- `limit`: maximum supplier rows from 1 to 5000; default 1000.
- `format=csv`: returns a bounded CSV export.

The report includes:

- supplier identity and code;
- purchase-order count;
- ordered, received, and open units;
- ordered and received purchasing value in minor currency units;
- average observed receipt lead time in days;
- completed-order count;
- on-time completed-order count.

Lead time is measured from purchase-order creation to the first recorded `receive` stock movement referencing that purchase order. On-time completion uses `expected_at` when it is present.

### Web reporting

Reports & Analytics now includes a **Supplier performance** panel with:

- supplier ranking by order activity;
- open-unit visibility;
- average lead-time display;
- completed/on-time order counts;
- CSV export for authorized users.

## Reliability and security

- Reporting is read-only.
- No database migration is required.
- Server-side aggregation keeps report calculation close to authoritative transactional data.
- Report windows and row counts are bounded.
- CSV exports use `no-store` / `no-cache` headers.
- Text fields are passed through the existing formula-safe CSV serialization path.
- Existing authentication, authorization, CORS, CSRF, and security-header behavior remains unchanged.

## Upgrade notes

No schema migration is required for v0.2.5. Existing purchase orders, purchase-order lines, suppliers, and receive movements are sufficient to populate the report.

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
