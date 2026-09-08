# StockPilot Roadmap

This roadmap reflects the current repository state.

## Completed foundation

- [x] Go domain model for catalog, warehouses, locations, inventory movements, transfers, lots, purchasing, and access control.
- [x] PostgreSQL schema, ordered migrations, and repository implementations.
- [x] Secure HTTP API with health/readiness endpoints, validation, CORS allowlisting, hardening headers, and request IDs.
- [x] Authentication with password verification, opaque sessions, CSRF protection, session-token protection, and RBAC.
- [x] Administrator bootstrap command.
- [x] React + TypeScript responsive web dashboard and PWA behavior.
- [x] Native Android authenticated dashboard with encrypted persisted sessions.
- [x] Manifest V3 browser companion with scoped optional host permissions.
- [x] Backend, frontend, Android, extension, CodeQL, and dependency-update CI baselines.
- [x] Aggregate reorder recommendations including zero-stock products.
- [x] Inventory valuation grouped by currency.
- [x] Exact barcode lookup API for scanner-driven workflows.
- [x] Reporting arithmetic, HTTP contract, and PostgreSQL integration coverage.
- [x] Security, contributor, architecture, restore, release, and API compatibility documentation.
- [x] Product catalog create/edit/search web workflow with role-aware controls.
- [x] Guided web inventory operations for stock-in, stock-out, signed adjustments, and transfers.
- [x] Warehouse/location administration workflow with create, edit, archive, reactivation, timezone, and warehouse association support.
- [x] Purchase-order creation, register, multi-line editing, multi-line receiving, lot/expiry handling, and atomic new-lot receiving.
- [x] Purchase-order lifecycle controls for draft submission and cancellation.
- [x] Append-only audit write path and audit viewer for sensitive business mutations.
- [x] Authentication/session audit events for login, session rejection, logout, and account lifecycle mutations.
- [x] Reorder-suggestion actions that seed reviewed draft purchase orders without auto-submitting them.
- [x] Lot inventory view with product, warehouse, location, lot, quantity, and expiry-risk filtering.
- [x] Browser camera barcode/QR scanning for product-form barcode entry, with manual fallback.
- [x] Android barcode/QR scanning with Google Code Scanner and product/lot-inventory handoff.
- [x] Browser companion barcode/QR scanning with safe handoff to the authenticated StockPilot web origin.
- [x] Browser companion workflow selection for product lookup and direct inventory-operation handoff.
- [x] CSV product dry-run validation with row-level errors, duplicate detection, and catalog-reference validation.
- [x] Transactional CSV product write/import with server-side revalidation and audit event.
- [x] Dedicated Reports & Analytics web workspace.
- [x] Inventory valuation CSV export with bounded, formula-safe, privacy-oriented delivery.
- [x] Inventory aging domain contracts and deterministic age-bucket rules.
- [x] Inventory aging PostgreSQL query and bounded HTTP report endpoint.
- [x] Inventory aging formula-safe CSV export through the report endpoint.
- [x] Web Reports workspace inventory-aging panel and export action.
- [x] Expiry-risk reporting by configurable date window.
- [x] Expiry-risk server-side classifications and bounded report/export workflow.
- [x] v0.2.3 CSV export formula-safety and timestamp regression hardening.
- [x] Stock movement history aggregation with configurable reporting window.
- [x] Stock movement velocity metrics and bounded CSV export.
- [x] Web Reports workspace movement-velocity panel and export action.

## Next — operational workflows

- [ ] Extend companion capabilities with a dedicated, independently revocable authentication design if direct extension API access is required.

## Next — data operations

- [ ] CSV inventory and report export with audit events.
- [ ] Documented automated backup retention deployment examples.
- [ ] Database integration tests covering concurrent inventory mutations and migration compatibility.

## Next — reporting and analytics

- [x] Inventory aging HTTP/report implementation and CSV export.
- [x] Expiry-risk report by configurable date window.
- [x] Stock movement history and velocity report.
- [x] Supplier purchasing totals and lead-time tracking.
- [x] Warehouse/location valuation breakdown.
- [x] Replenishment performance history with fill-rate and receipt-timeliness metrics.
- [x] Explicit traceability between reorder suggestions and resulting purchase-order decisions.
- [x] Deterministic cursor pagination for replenishment-readiness JSON reports.
- [x] Dedicated web replenishment review-history workspace.
- [ ] Cursor/streaming support for large report datasets beyond replenishment-readiness.

## Release hardening

- [x] Deterministic Chromium browser E2E foundation.
- [x] Signed-out login regression coverage.
- [x] Authenticated dashboard-shell coverage using synthetic API fixtures.
- [x] Authenticated primary-workspace navigation coverage using synthetic API fixtures.
- [x] Authenticated catalog create workflow coverage with synthetic test data.
- [x] Authenticated inventory movement workflow coverage with synthetic test data.
- [ ] Authenticated purchase-order draft/receive mutation workflow coverage with safe test data.
- [ ] Authenticated reporting workspace coverage with complete report fixtures.
- [ ] Android instrumentation tests for authentication and critical inventory reads.
- [ ] Accessibility audit across keyboard, screen-reader, contrast, touch-target, and reduced-motion paths.
- [ ] Production restore rehearsal and migration rollback verification.
- [ ] API compatibility/versioning policy wired into CI checks.
- [ ] Reproducible web, server, Android, and extension release artifacts.
- [ ] Resolve every blocker/critical defect before the first stable release.

## v0.5.1 milestone

- [x] Align `VERSION` and `/api/v1/meta` with `0.5.1`.
- [x] Align the API metadata regression test with `0.5.1`.
- [x] Add deterministic authenticated catalog-create browser coverage.
- [x] Add deterministic authenticated inventory stock-in browser coverage.
- [x] Assert representative mutation request payloads in browser tests.
- [x] Keep all mutation E2E fixtures synthetic and independent of production credentials/data.
- [ ] Add authenticated purchase-order create/submit/receive browser coverage.
- [ ] Add complete authenticated reporting browser fixtures and assertions.
- [ ] Run and verify the complete release CI gate on the final release commit before publishing the Git tag.
