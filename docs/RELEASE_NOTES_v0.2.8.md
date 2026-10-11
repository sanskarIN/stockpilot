# StockPilot v0.2.8 — Large-Report Scalability & Purchasing Trends

**Version:** `0.2.8`  
**Tag:** `v0.2.8`  
**Release type:** Stable after verification gates pass  
**Preparation date:** 2026-09-04

## Overview

StockPilot v0.2.8 is planned as a reporting scalability and procurement-analytics milestone. It will make large report datasets safer to consume while adding purchasing trend visibility for products and suppliers.

## Planned highlights

### Large-report scalability

- Add cursor/keyset pagination or streaming where current reports can grow beyond comfortable in-memory result sizes.
- Preserve deterministic ordering across pages and streams.
- Keep continuation tokens opaque and bounded.
- Avoid duplicate or skipped rows when consumers continue from a cursor.
- Prefer server-side iteration and indexed queries over unbounded in-memory aggregation.

### Purchasing trends

Add read-only trend reporting derived from existing purchase-order and purchase-order-line records, with metrics such as:

- ordered units;
- received units;
- open units;
- purchasing value in the stored currency;
- supplier/product dimensions supported by the current schema.

Mixed currencies must remain separated; the report must not perform implicit FX conversion.

### API and export safety

Planned reporting surfaces will use:

- bounded date windows;
- bounded result sizes;
- deterministic ordering;
- formula-safe CSV serialization;
- `no-store` / `no-cache` download headers;
- existing authentication and authorization controls.

Backward compatibility should be preserved for existing report endpoints wherever practical. New pagination should not silently change the meaning of an existing endpoint.

### Web reporting

Reports & Analytics should gain purchasing-trend views with:

- loading, empty, and error states;
- accessible table/chart alternatives;
- clear date-window context;
- supplier/product filters where supported;
- CSV export affordances for authorized users.

## Reliability and security requirements

- Reporting remains read-only.
- No automatic purchase-order creation or inventory mutation is part of this milestone.
- No credentials, session tokens, passwords, or payment information may enter report output.
- Existing CORS, CSRF, security headers, authentication, and authorization behavior must remain intact.
- Memory and query-result growth must remain bounded.
- A database migration should be avoided unless implementation evidence shows the existing schema cannot support the required reporting behavior.

## Verification gates

Before publishing v0.2.8 as stable, verify:

- Go formatting, vet, unit tests, race tests, and server build;
- PostgreSQL migration/readiness and integration checks;
- Web typecheck and production build;
- CodeQL/security checks;
- Android lint/tests/build where configured;
- browser companion checks where configured;
- end-to-end and accessibility checks where configured;
- pagination/streaming correctness under realistic dataset sizes;
- restore/rollback checks where applicable;
- release artifact/reproducibility checks where configured;
- blocker and critical-defect review.

## Upgrade notes

v0.2.8 is intended to be backward compatible. Existing transactional data remains authoritative, and any new reporting contracts should be additive unless a compatibility review explicitly approves a change.

## Release metadata

- **Version:** `v0.2.8`
- **Title:** `StockPilot v0.2.8 — Large-Report Scalability & Purchasing Trends`
- **Tag:** `v0.2.8`
- **Release type:** Stable after verification
- **Prerelease:** No
- **Planned date:** TBD after implementation and verification

## Publication rule

These are preparation notes. Do not represent v0.2.8 as published until implementation, verification, and the GitHub Release publication have actually completed.
