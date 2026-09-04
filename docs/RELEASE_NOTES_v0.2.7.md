# StockPilot v0.2.7 — Replenishment Readiness

**Version:** `0.2.7`  
**Tag:** `v0.2.7`  
**Release type:** Stable only after verification gates pass  
**Release date:** 2026-09-04

## Overview

StockPilot v0.2.7 is the next reporting milestone after warehouse/location valuation. Its focus is making replenishment decisions more transparent by connecting current reorder recommendations with recent outbound movement behavior.

## Planned release scope

### Replenishment readiness

The release scope is centered on a read-only replenishment insight that can expose, for eligible active products:

- current on-hand quantity;
- configured reorder point;
- configured reorder quantity;
- target stock;
- suggested replenishment quantity;
- recent outbound units;
- average daily outbound velocity;
- estimated days of cover when recent outbound velocity is positive;
- a deterministic readiness/risk classification.

The calculation must remain advisory. It must not automatically create purchase orders or mutate inventory.

### API design

The intended reporting surface is a bounded authenticated read-only endpoint under `/api/v1/reports/` with:

- a configurable historical demand window from 1 to 365 days;
- a bounded result limit from 1 to 5000;
- deterministic ordering;
- optional formula-safe CSV output;
- `no-store` / `no-cache` response headers for CSV downloads.

The implementation should reuse authoritative persisted inventory movements, product configuration, and existing reorder-suggestion logic rather than introducing a second source of truth.

### Web reporting

Reports & Analytics should expose replenishment readiness alongside the existing inventory valuation, aging, movement velocity, supplier performance, and purchasing views.

The UI should clearly distinguish:

- configured thresholds;
- observed demand;
- calculated estimates;
- advisory recommendations.

No destructive or automatic purchasing action should be introduced by this reporting milestone.

## Reliability and security requirements

- Reporting remains read-only.
- Existing authentication and authorization rules remain enforced server-side.
- Historical windows and result sets remain bounded.
- CSV output must use the existing formula-safe serializer.
- CSV responses must retain `no-store` / `no-cache` headers.
- No credentials, session tokens, passwords, or payment information are included in reports.
- No database migration should be introduced unless implementation evidence shows persisted historical recommendation snapshots are required; prefer derived reporting from existing authoritative records.

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

Do not publish the stable GitHub release until the applicable gates are green.

## Upgrade notes

The target implementation is intended to be backward compatible with the v0.2.6 data model. Existing catalog, inventory, movement, and purchasing records remain authoritative.

## Release metadata

- **Version:** `v0.2.7`
- **Title:** `StockPilot v0.2.7 — Replenishment Readiness`
- **Tag:** `v0.2.7`
- **Release type:** Stable after verification
- **Prerelease:** No
- **Latest:** Yes, if this is the current latest stable release
- **Date:** 2026-09-04

## Publication rule

This document describes the v0.2.7 release scope and readiness requirements. The GitHub Release must not be represented as published until the implementation and verification gates are complete.
