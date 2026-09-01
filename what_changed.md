# StockPilot — Work Continuity Log

## Current milestone

Phase 4 — purchasing and receiving workflows.

## Repository state

- Default branch: `main`.
- Merged mainline includes Android/browser quality gates, reporting/replenishment, catalog management, guided inventory operations, and warehouse/location administration.
- Superseded feature PRs are closed after their work was rebuilt cleanly on current `main`.
- Active branch: `feat/purchase-order-workflow-v2`.
- The continuation intentionally preserves focused commits instead of squashing feature history.

## Completed in this continuation

- [x] Added detailed purchase-order line and lifecycle types for the web client.
- [x] Added list, create, and receiving API client methods.
- [x] Added purchase-order register with supplier, warehouse, status, and dates.
- [x] Added draft purchase-order creation with validation for number, supplier, warehouse, product, quantity, cost, and currency.
- [x] Added purchase-order detail with ordered versus received quantities.
- [x] Added controlled receiving into a selected inventory location.
- [x] Added dashboard purchasing navigation and application routing.
- [x] Added responsive purchase-order and receiving presentation styles.
- [x] Kept existing backend authenticated-order/receiving HTTP coverage on the mainline.
- [x] Updated roadmap and changelog for the purchasing milestone.

## Verification status

The connected GitHub environment does not expose a local project shell, so this continuation does not claim local Go/web/Android/extension command execution. Existing GitHub Actions remain the authoritative validation path.

## Known limitations

- The initial purchase-order UI creates one line per new draft.
- Explicit ordered/cancelled workflow controls remain pending.
- Lot/expiry-aware receiving remains pending.
- Reorder recommendations do not yet seed draft purchase orders directly.
- Append-only audit writes and audit viewer remain pending.
- CSV import/export and advanced analytics remain pending.
- Browser E2E and Android instrumentation suites remain pending.

## Next exact tasks

1. Add lot and expiry-aware receiving with configurable expiry warnings.
2. Add purchase-order multi-line editing and explicit approval/lifecycle transitions.
3. Add audit event writes for sensitive mutations and an audit viewer.
4. Add CSV import/export with dry-run validation and row-level errors.
5. Add inventory aging, expiry-risk, movement-velocity, supplier, and replenishment analytics.
6. Add browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks.

## Commit discipline

Each functional boundary is intentionally represented by an individual commit where practical: models, API methods, UI, dashboard integration, routing, styling, tests, fixes, and documentation are kept separately reviewable. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution.
