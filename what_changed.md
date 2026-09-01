# StockPilot — Work Continuity Log

## Current milestone

Phase 5 — lot and expiry-aware receiving.

## Repository state

- Default branch: `main`.
- Merged mainline includes Android/browser quality gates, reporting/replenishment, catalog management, guided inventory operations, warehouse/location administration, and purchase-order creation/receiving.
- Superseded feature PRs were closed after their work was rebuilt cleanly on current `main`.
- Active branch: `feat/lot-expiry-receiving`.
- The continuation intentionally preserves focused commits instead of squashing feature history.

## Completed in this continuation

- [x] Added typed web lot model with manufacturing and expiry metadata.
- [x] Added web lot creation and reusable lot-listing API methods.
- [x] Added repository contract and PostgreSQL implementation for product lot listing.
- [x] Added authenticated `GET /api/v1/lots?productId=...` endpoint and regression coverage.
- [x] Added lot-aware purchase receiving for products configured with lot tracking.
- [x] Added existing-lot reuse so partial receipts do not create duplicate batches.
- [x] Added new-lot creation flow with lot number and optional manufacturing date.
- [x] Added mandatory expiry date for expiry-tracked products.
- [x] Added near-expiry and already-expired visual warnings.
- [x] Kept product/lot ownership and lot-tracking rules enforced server-side.
- [x] Updated roadmap and changelog to reflect completed lot/expiry receiving.

## Verification status

The connected GitHub environment does not expose a local project shell, so this continuation does not claim local Go/web/Android/extension command execution. Existing GitHub Actions remain the authoritative validation path.

## Known limitations

- Lot creation and receipt are currently two HTTP operations, so an unused lot can remain if receipt fails after creation; an atomic receive-and-create operation is a future hardening task.
- Purchase-order receiving still focuses on the first line in the current UI; full multi-line receiving remains pending.
- Explicit purchase-order lifecycle controls remain pending.
- Reorder recommendations do not yet seed draft purchase orders directly.
- Append-only audit writes and audit viewer remain pending.
- CSV import/export and advanced analytics remain pending.
- Browser E2E and Android instrumentation suites remain pending.

## Next exact tasks

1. Add an atomic server-side receive-and-create-lot operation to eliminate orphaned lot risk.
2. Add purchase-order multi-line editing and explicit ordered/cancelled lifecycle controls.
3. Add audit event writes for sensitive mutations and an audit viewer.
4. Add CSV import/export with dry-run validation and row-level errors.
5. Add inventory aging, expiry-risk, movement-velocity, supplier, and replenishment analytics.
6. Add browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks.

## Commit discipline

Each functional boundary is intentionally represented by an individual commit where practical: contracts, persistence, HTTP handlers, client API, UI, styling, tests, fixes, and documentation are kept separately reviewable. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution.
