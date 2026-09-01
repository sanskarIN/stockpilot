# StockPilot — Work Continuity Log

## Current milestone

Phase 10 — reorder-to-purchase-order workflow.

## Repository state

- Default branch: `main`.
- Merged mainline includes Android/browser quality gates, reporting/replenishment, catalog management, guided inventory operations, warehouse/location administration, multi-line purchasing and receiving, lot/expiry-aware receiving, atomic new-lot receiving, purchase-order lifecycle controls, and append-only auditability with a web viewer.
- Active branch: `feat/reorder-to-purchase-order`.
- The continuation intentionally preserves focused commits instead of squashing feature history.

## Completed in this continuation

- [x] Added a Create draft action to dashboard reorder suggestions.
- [x] Passed the selected reorder suggestion into the purchasing workflow without persisting it immediately.
- [x] Prefilled suggested product, supplier, quantity, currency, target-stock context, and a review note.
- [x] Kept the normal draft creation, editing, submission, cancellation, and receiving controls in place.
- [x] Added a one-time seed-consumption guard so closing a seeded draft does not reopen it automatically.
- [x] Preserved server-authoritative approval/lifecycle behavior: the reorder shortcut creates only a client-side draft proposal until the user explicitly saves it.

## Verification status

The connected GitHub environment does not expose a local project shell, so this continuation does not claim local Go/web/Android/extension command execution. Existing GitHub Actions remain the authoritative validation path.

## Known limitations

- Reorder seeding is currently a web workflow convenience; the server still receives an ordinary draft purchase-order create request after user confirmation.
- Suggested supplier falls back to the first active supplier when the recommendation has no supplier association.
- Warehouse selection uses the normal purchasing default rather than a reorder-specific destination.
- Warehouse/location edit and archive workflows remain pending.
- Lot inventory views by location/quantity/expiry remain pending.
- Authentication/session audit events are not yet emitted.
- CSV import/export and advanced analytics remain pending.
- Browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks remain pending.

## Next exact tasks

1. Add warehouse/location edit and archive workflows.
2. Add lot inventory views by location, quantity, and expiry sorting.
3. Extend audit coverage to authentication/session and future import/export mutations.
4. Add CSV import/export with dry-run validation and row-level errors.
5. Add inventory aging, expiry-risk, movement-velocity, supplier, and replenishment analytics.
6. Add browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks.

## Commit discipline

Each functional boundary is intentionally represented by an individual commit where practical: domain rules, repository contracts, persistence, HTTP handlers, tests, client API, UI, routing, styling, fixes, and documentation are kept separately reviewable. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution.
