# StockPilot — Work Continuity Log

## Current milestone

Phase 8 — auditability and mutation traceability.

## Repository state

- Default branch: `main`.
- Merged mainline includes Android/browser quality gates, reporting/replenishment, catalog management, guided inventory operations, warehouse/location administration, purchase-order creation/receiving, lot/expiry-aware receiving, atomic new-lot receiving, and purchase-order lifecycle controls.
- Active branch: `feat/audit-write-path`.
- The continuation intentionally preserves focused commits instead of squashing feature history.

## Completed in this continuation

- [x] Extended the audit repository contract with an append-only write operation.
- [x] Added domain validation for audit actions, entity types, entity IDs, and metadata.
- [x] Added PostgreSQL audit inserts into the existing `audit_log` table without update/delete operations.
- [x] Propagated the generated HTTP request ID through request context for audit correlation.
- [x] Added a centralized completed-mutation audit recorder that preserves business-operation success when audit persistence is unavailable.
- [x] Added audit events for product/category/supplier mutations.
- [x] Added audit events for warehouse/location/lot creation and inventory movement/transfer mutations.
- [x] Added audit events for purchase-order creation, lifecycle transitions, and receiving.
- [x] Added web audit API support with actor/action/entity filters and pagination.
- [x] Added a read-only audit history screen with request IDs and JSON metadata visibility.
- [x] Added role-aware audit navigation; operator accounts are denied audit-history access.
- [x] Updated the roadmap and changelog for the audit milestone.

## Verification status

The connected GitHub environment does not expose a local project shell, so this continuation does not claim local Go/web/Android/extension command execution. Existing GitHub Actions remain the authoritative validation path.

## Known limitations

- Authentication/session mutations are not yet emitted as audit events.
- Future CSV import/export mutations need to use the same recorder.
- The audit recorder intentionally does not make an otherwise successful business mutation fail when audit persistence is unavailable; deployment monitoring should surface audit-storage errors.
- Purchase-order UI still exposes one line per newly created order; multi-line editing/receiving remains pending.
- Lot inventory views, reorder-to-draft automation, CSV tooling, and advanced analytics remain pending.
- Browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks remain pending.

## Next exact tasks

1. Add multi-line purchase-order editing and receiving.
2. Add reorder-to-draft actions with approval-aware behavior.
3. Add lot inventory views and expiry-risk analytics.
4. Extend audit coverage to authentication/session and future import/export mutations.
5. Add CSV import/export with dry-run validation and row-level errors.
6. Add inventory aging, movement-velocity, supplier, and replenishment analytics.
7. Add browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks.

## Commit discipline

Each functional boundary is intentionally represented by an individual commit where practical: domain rules, repository contracts, persistence, HTTP handlers, tests, client API, UI, routing, cleanup, and documentation are kept separately reviewable. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution.
