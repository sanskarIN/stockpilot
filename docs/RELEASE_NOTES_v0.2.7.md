# StockPilot v0.2.7 — Replenishment Readiness

**Version:** `0.2.7`  
**Tag:** `v0.2.7`  
**Release type:** Stable after verification gates pass  
**Release date:** 2026-09-04

## Overview

StockPilot v0.2.7 adds a read-only replenishment-readiness report that combines the existing reorder-suggestion model with recent stock-movement velocity. It makes the operational reason behind a replenishment review easier to see without introducing a second inventory source of truth or automatic purchasing.

## Highlights

### Replenishment readiness report

New endpoint:

`GET /api/v1/reports/replenishment-readiness`

Supported query parameters:

- `days` — historical demand window from 1 to 365 days; default 30.
- `limit` — maximum result rows from 1 to 5000; default 1000.
- `format=csv` — bounded formula-safe CSV export.

Each item exposes:

- product identity and SKU;
- supplier and unit metadata;
- on-hand quantity;
- configured reorder point and reorder quantity;
- target stock and suggested quantity;
- recent outbound units;
- average daily outbound velocity;
- estimated days of cover when velocity is positive;
- deterministic advisory risk classification.

### Risk classification

The report uses deterministic, explainable categories:

- `out_of_stock` — no on-hand units;
- `critical` — positive stock with fewer than 7 days of cover when outbound velocity is available;
- `reorder` — on-hand quantity is at or below the configured reorder point;
- `watch` — fewer than 14 days of cover when velocity is available;
- `healthy` — none of the above conditions apply.

The classification is advisory and does not change inventory or create purchase orders.

### Data integrity

The report reuses existing reorder suggestions and persisted stock-movement history. No new database migration or recommendation snapshot table is introduced.

### CSV and API safety

- Result windows and row counts are bounded server-side.
- CSV output uses the existing formula-safe serializer.
- CSV downloads use `no-store` / `no-cache` headers.
- Existing authentication, authorization, CORS, CSRF, and security-header behavior remains unchanged.

## Verification gates

Before publishing v0.2.7 as stable, verify:

- Go formatting, vet, unit tests, race tests, and server build;
- PostgreSQL migration/readiness and integration checks;
- Web typecheck and production build;
- CodeQL/security checks;
- Android lint/tests/build where configured;
- browser companion checks where configured;
- end-to-end and accessibility checks where configured;
- restore/rollback checks where applicable;
- release artifact/reproducibility checks where configured;
- blocker and critical-defect review.

## Upgrade notes

No schema migration is required for the replenishment-readiness implementation. Existing catalog, inventory, movement, and reorder configuration remain authoritative.

## Release metadata

- **Version:** `v0.2.7`
- **Title:** `StockPilot v0.2.7 — Replenishment Readiness`
- **Tag:** `v0.2.7`
- **Release type:** Stable after verification
- **Prerelease:** No
- **Date:** 2026-09-04

## Publication rule

The implementation is now present on `main`, but the GitHub stable release must only be published after the applicable verification gates are confirmed green.
