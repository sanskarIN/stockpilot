# StockPilot — Work Continuity Log

## Current milestone

Phase 13 — browser camera barcode/QR scanning.

## Repository state

- Default branch: `main`.
- Merged mainline includes Android/browser quality gates, reporting/replenishment, catalog management, guided inventory operations, warehouse/location lifecycle management, multi-line purchasing and receiving, lot/expiry-aware receiving, atomic new-lot receiving, purchase-order lifecycle controls, append-only auditability with a web viewer, reorder-to-draft assistance, and lot inventory visibility.
- Active branch: `feat/barcode-camera-scanner`.
- The continuation intentionally preserves focused commits instead of squashing feature history.

## Completed in this continuation

- [x] Added a reusable web barcode/QR scanner using `BarcodeDetector` when the browser provides it.
- [x] Requested audio-free, rear-facing camera video where supported.
- [x] Added deterministic scanner cleanup that stops camera tracks when the scanner closes or unmounts.
- [x] Added graceful fallback for browsers without camera detection support or without granted camera access.
- [x] Added manual barcode entry so product management remains usable without camera APIs.
- [x] Connected scanning to the existing product barcode field without auto-saving.
- [x] Preserved the normal product validation and save flow after scanning.
- [x] Added responsive scanner presentation and compact mobile controls.
- [x] Updated roadmap and changelog for the browser scanning milestone.

## Verification status

The connected GitHub environment does not expose a local project shell, so this continuation does not claim local Go/web/Android/extension command execution. Existing GitHub Actions remain the authoritative validation path.

## Known limitations

- Camera scanning is currently implemented in the web product form; Android inventory scanning remains pending.
- Browser support depends on `BarcodeDetector` availability and camera permissions.
- Scanned values populate an unsaved draft; users must explicitly save a product.
- Browser companion scan-to-product/inventory actions remain pending.
- Authentication/session audit events are not yet emitted.
- CSV import/export and advanced analytics remain pending.
- Browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks remain pending.

## Next exact tasks

1. Add Android camera barcode/QR scanning and inventory/product handoff.
2. Add browser companion scan-to-product/inventory action.
3. Extend audit coverage to authentication/session and future import/export mutations.
4. Add CSV product import with dry-run validation and row-level errors.
5. Add CSV inventory/report export.
6. Add inventory aging, expiry-risk, movement-velocity, supplier, and replenishment analytics.
7. Add browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks.

## Commit discipline

Each functional boundary is intentionally represented by an individual commit where practical: domain rules, repository contracts, persistence, HTTP handlers, tests, client API, UI, routing, styling, fixes, and documentation are kept separately reviewable. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution.
