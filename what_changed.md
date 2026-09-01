# StockPilot — Work Continuity Log

## Current milestone

Phase 12 — lot inventory visibility.

## Repository state

- Default branch: `main`.
- Merged mainline includes Android/browser quality gates, reporting/replenishment, catalog management, guided inventory operations, warehouse/location lifecycle management, multi-line purchasing and receiving, lot/expiry-aware receiving, atomic new-lot receiving, purchase-order lifecycle controls, append-only auditability with a web viewer, and reorder-to-draft assistance.
- Active branch: `feat/lot-inventory-view`.
- The continuation intentionally preserves focused commits instead of squashing feature history.

## Completed in this continuation

- [x] Added a lot-inventory domain read model covering product, lot, location, warehouse, quantity, and expiry.
- [x] Added filtered lot inventory repository support with bounded pagination.
- [x] Added PostgreSQL lot inventory projection from authoritative `inventory_balances` lot rows.
- [x] Added product, warehouse, location, lot, and expiry-cutoff filters.
- [x] Added strict `YYYY-MM-DD` expiry query validation.
- [x] Added web lot inventory API client and typed row model.
- [x] Added dedicated web lot inventory screen with expiry-risk classification.
- [x] Added dashboard/navigation entry for lot inventory.
- [x] Added responsive lot-inventory styling and pagination controls.
- [x] Added HTTP regression coverage for lot inventory retrieval and expiry-date validation.
- [x] Updated roadmap and changelog for the lot inventory milestone.

## Verification status

The connected GitHub environment does not expose a local project shell, so this continuation does not claim local Go/web/Android/extension command execution. Existing GitHub Actions remain the authoritative validation path.

## Known limitations

- The lot inventory endpoint currently returns positive on-hand lot balances only.
- Expiry-risk labels are computed in the web client from the current browser date.
- Authentication/session audit events are not yet emitted.
- CSV import/export and advanced analytics remain pending.
- Browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks remain pending.
- Barcode/QR camera scanning UI remains pending.

## Next exact tasks

1. Extend audit coverage to authentication/session and future import/export mutations.
2. Add CSV product import with dry-run validation and row-level errors.
3. Add CSV inventory/report export.
4. Add inventory aging, configurable expiry-risk, movement-velocity, supplier, and replenishment analytics.
5. Add barcode/QR camera scanning UI for supported clients.
6. Add browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks.
7. Add automated backup retention examples and concurrent database integration coverage.

## Commit discipline

Each functional boundary is intentionally represented by an individual commit where practical: domain rules, repository contracts, persistence, HTTP handlers, tests, client API, UI, routing, styling, fixes, and documentation are kept separately reviewable. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution.
