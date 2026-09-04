# StockPilot v0.2.9 — Operational Readiness & Reporting Reliability

**Version:** `0.2.9`  
**Tag:** `v0.2.9`  
**Release type:** Stable after verification gates pass  
**Preparation date:** 2026-09-04

## Overview

StockPilot v0.2.9 is the next planned pre-1.0 hardening milestone after v0.2.8. The goal is to make reporting and day-to-day operations more observable, predictable, and resilient without changing the transactional source of truth.

## Planned highlights

### Reporting reliability

- Add explicit report query budgets and bounded execution behavior where supported.
- Standardize report response metadata for generated-at timestamps, applied bounds, and partial/complete result state.
- Make expensive report failures actionable without exposing database internals.
- Preserve deterministic ordering for repeatable report consumption.

### Operational observability

- Improve structured request/report logging with stable event names and duration measurements.
- Add safe diagnostics for slow report paths without recording credentials, session tokens, or sensitive request bodies.
- Document health, readiness, and operational troubleshooting expectations.

### API resilience

- Review timeout and cancellation propagation across report and read-heavy endpoints.
- Add defensive validation for query parameters, date ranges, limits, and continuation inputs.
- Keep existing API contracts backward compatible unless an explicitly reviewed additive change is required.

### Data and export safety

- Continue formula-safe CSV serialization and non-cacheable download behavior.
- Verify that large exports remain bounded and do not accidentally materialize unbounded datasets in memory.
- Preserve currency separation and source-of-truth semantics for purchasing analytics.

### Client experience

- Improve Reports & Analytics loading, empty, error, and retry states.
- Surface useful report-bound context to users.
- Keep accessibility and keyboard navigation as release gates for new reporting UI.

## Reliability and security requirements

- Reporting remains read-only.
- No automatic purchasing or inventory mutation is part of this milestone.
- No credentials, passwords, session tokens, payment information, or sensitive request payloads may enter logs or report exports.
- Existing authentication, authorization, CORS, CSRF, and security-header behavior must remain intact.
- Database schema changes require evidence and compatibility review; avoid migrations where existing structures are sufficient.

## Verification gates

Before publishing v0.2.9 as stable, verify:

- Go formatting, vet, unit tests, race tests, and server build;
- PostgreSQL integration/readiness checks;
- report cancellation, timeout, bound, and export regression tests;
- Web typecheck and production build;
- CodeQL/security checks;
- Android lint/tests/build where configured;
- browser companion checks where configured;
- E2E/accessibility checks where configured;
- operational logging review for sensitive-data leakage;
- restore/rollback and release-artifact checks where configured;
- blocker and critical-defect review.

## Upgrade notes

v0.2.9 is intended to remain backward compatible. Existing transactional data and established reporting semantics remain authoritative. New reliability metadata and operational safeguards should be additive.

## Release metadata

- **Version:** `v0.2.9`
- **Title:** `StockPilot v0.2.9 — Operational Readiness & Reporting Reliability`
- **Tag:** `v0.2.9`
- **Release type:** Stable after verification
- **Prerelease:** No
- **Planned date:** TBD after implementation and verification
