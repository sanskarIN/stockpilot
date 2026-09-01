# StockPilot — Work Continuity Log

## Current milestone

Phase 11 — warehouse/location lifecycle.

## Repository state

- Default branch: `main`.
- Merged mainline includes Android/browser quality gates, reporting/replenishment, catalog management, guided inventory operations, warehouse/location administration, multi-line purchasing and receiving, lot/expiry-aware receiving, atomic new-lot receiving, purchase-order lifecycle controls, append-only auditability with a web viewer, and reorder-to-draft assistance.
- Active branch: `feat/warehouse-location-lifecycle`.
- The continuation intentionally preserves focused commits instead of squashing feature history.

## Completed in this continuation

- [x] Added repository contracts for warehouse and location updates.
- [x] Added transactional PostgreSQL warehouse/location update persistence.
- [x] Added server-side archive guards for warehouses with active locations.
- [x] Added server-side archive guards for locations with non-zero inventory.
- [x] Prevented activation of locations under archived warehouses.
- [x] Added authenticated PUT endpoints for warehouse and location updates.
- [x] Added audit events for warehouse and location updates.
- [x] Added warehouse/location web edit, archive, and reactivation controls.
- [x] Added management visibility for inactive records.
- [x] Added HTTP regression coverage for warehouse editing and location archival.
- [x] Added dedicated responsive lifecycle styling.
- [x] Updated roadmap and changelog for the lifecycle milestone.

## Verification status

The connected GitHub environment does not expose a local project shell, so this continuation does not claim local Go/web/Android/extension command execution. Existing GitHub Actions remain the authoritative validation path.

## Known limitations

- Warehouse edits do not yet expose per-field conflict/version detection beyond row locking.
- Lot inventory views by location/quantity/expiry remain pending.
- Authentication/session audit events are not yet emitted.
- CSV import/export and advanced analytics remain pending.
- Barcode/QR camera scanning UI remains pending.
- Browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks remain pending.

## Next exact tasks

1. Add lot inventory views by location, quantity, and expiry sorting.
2. Extend audit coverage to authentication/session and future import/export mutations.
3. Add CSV product import with dry-run validation and row-level errors.
4. Add CSV inventory/report export.
5. Add inventory aging, expiry-risk, movement-velocity, supplier, and replenishment analytics.
6. Add browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks.
7. Add barcode/QR camera scanning UI for supported clients.

## Commit discipline

Each functional boundary is intentionally represented by an individual commit where practical: domain rules, repository contracts, persistence, HTTP handlers, tests, client API, UI, routing, styling, fixes, and documentation are kept separately reviewable. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution.
