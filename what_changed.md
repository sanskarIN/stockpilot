# StockPilot — Work Continuity Log

## Current milestone

Phase 7 — purchase-order lifecycle controls.

## Repository state

- Default branch: `main`.
- Merged mainline includes Android/browser quality gates, reporting/replenishment, catalog management, guided inventory operations, warehouse/location administration, purchase-order creation/receiving, lot/expiry-aware receiving, and atomic new-lot receiving.
- Active branch: `feat/purchase-order-lifecycle`.
- The continuation intentionally preserves focused commits instead of squashing feature history.

## Completed in this continuation

- [x] Added a domain-level purchase-order lifecycle state machine.
- [x] Allowed deliberate `draft → ordered` submission and `draft/ordered → cancelled` transitions.
- [x] Kept `partially_received` and `received` states controlled by receipt processing instead of manual status edits.
- [x] Added row-locked PostgreSQL status persistence.
- [x] Added authenticated `PATCH /api/v1/orders/{id}/status` HTTP endpoint.
- [x] Added domain and HTTP regression coverage for lifecycle rules.
- [x] Added web lifecycle API client support.
- [x] Added a unified purchasing screen covering create, submit/cancel, lot-aware receive, and receipt progress.
- [x] Routed purchasing through the unified workflow and removed the superseded receiving-only screen.
- [x] Updated roadmap and changelog for the lifecycle milestone.

## Verification status

The connected GitHub environment does not expose a local project shell, so this continuation does not claim local Go/web/Android/extension command execution. Existing GitHub Actions remain the authoritative validation path.

## Known limitations

- Purchase-order UI still exposes one line per newly created order.
- Multi-line order editing and receiving remain pending.
- Reorder recommendations do not yet seed draft purchase orders directly.
- Lot inventory views by location/quantity/expiry remain pending.
- Append-only audit writes and audit viewer remain pending.
- CSV import/export and advanced analytics remain pending.
- Browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks remain pending.

## Next exact tasks

1. Add multi-line purchase-order editing and receiving.
2. Add reorder-to-draft actions with approval-aware behavior.
3. Add lot inventory views and expiry-risk analytics.
4. Add append-only audit writes and an audit viewer.
5. Add CSV import/export with dry-run validation and row-level errors.
6. Add inventory aging, movement-velocity, supplier, and replenishment analytics.
7. Add browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks.

## Commit discipline

Each functional boundary is intentionally represented by an individual commit where practical: domain rules, repository contracts, persistence, HTTP handlers, tests, client API, UI, routing, cleanup, and documentation are kept separately reviewable. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution.
