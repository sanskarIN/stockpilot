# StockPilot v0.2.10 — Reporting Performance & Operational Insights

**Version:** `0.2.10`  
**Tag:** `v0.2.10`  
**Release type:** Stable only after implementation and verification  
**Preparation date:** 2026-09-04

## Status

v0.2.10 is the next planned pre-1.0 milestone after the v0.2.9 operational-readiness track. This document is a release plan, not a claim that v0.2.10 has been implemented or published.

## Goals

The milestone focuses on making StockPilot reporting faster and more useful for day-to-day decisions while preserving authoritative transactional data and backward-compatible APIs.

## Planned highlights

### Reporting performance

- Reuse bounded pagination and continuation primitives established during the large-report scalability work.
- Add deterministic ordering and explicit limits to newly introduced analytical queries.
- Reduce repeated database work for common report views where safe and measurable.
- Preserve cancellation, timeout, and request-bound propagation through report execution.

### Operational insights

- Add additive trend summaries for inventory, purchasing, and replenishment activity.
- Surface period-over-period changes without silently mixing currencies.
- Keep source-of-truth identifiers available for traceability.
- Clearly distinguish observed historical metrics from advisory calculations.

### API quality

- Keep existing endpoints backward compatible.
- Validate date windows, limits, pagination inputs, and unsupported combinations defensively.
- Return stable response shapes with bounded metadata and actionable public errors.
- Preserve authenticated access controls and existing security middleware.

### Export and data safety

- Keep CSV output formula-safe and non-cacheable.
- Bound export rows and memory consumption.
- Never place credentials, session tokens, payment information, or sensitive request payloads in logs or exports.
- Preserve currency separation and authoritative database semantics.

### Web experience

- Add clear report period/bound context.
- Provide loading, empty, error, retry, and no-data states.
- Preserve keyboard accessibility and semantic table alternatives.
- Avoid introducing client-side calculations that diverge from server-side report semantics.

## Engineering constraints

- Reporting remains read-only.
- No automatic purchase-order creation or inventory mutation.
- Avoid schema migrations unless implementation evidence proves they are necessary and compatibility is reviewed.
- Prefer small, focused, reviewable commits over artificial commit-count inflation.

## Verification gates

Before publishing v0.2.10 as stable, verify:

- Go format, vet, unit tests, race tests, and production build;
- PostgreSQL integration and query-bound tests;
- report cancellation, timeout, pagination, ordering, and duplicate-prevention tests;
- CSV safety and bounded-memory export tests;
- Web lint/type-check/production build;
- CodeQL/security checks;
- Android/browser companion checks where configured;
- E2E and accessibility checks where configured;
- logging/privacy review;
- rollback and release-artifact checks where configured;
- blocker and critical-defect review.

## Upgrade notes

v0.2.10 is intended to be additive and backward compatible. Existing inventory, purchasing, and reporting source-of-truth semantics remain authoritative.

## Publication rule

Do not create or announce the stable `v0.2.10` GitHub Release until implementation, verification, tag creation, and release publication have actually completed.
