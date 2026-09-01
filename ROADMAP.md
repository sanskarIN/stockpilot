# StockPilot Roadmap

This roadmap reflects the current repository state.

## Completed foundation

- [x] Go domain model for catalog, warehouses, locations, inventory movements, transfers, lots, purchasing, and access control.
- [x] PostgreSQL schema, ordered migrations, and repository implementations.
- [x] Secure HTTP API with health/readiness endpoints, validation, CORS allowlisting, hardening headers, and request IDs.
- [x] Authentication with password verification, opaque sessions, CSRF protection, session-token protection, and RBAC.
- [x] Administrator bootstrap command.
- [x] React + TypeScript responsive dashboard and PWA behavior.
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

## Next — operational workflows

- [ ] Warehouse/location management screens.
- [ ] Purchase-order creation, editing, approval-state controls, and receiving UI.
- [ ] Barcode/QR camera scanning UI for supported clients.
- [ ] Lot and expiry receiving workflows with expiry warnings.
- [ ] Reorder suggestion actions that can seed draft purchase orders without bypassing approval rules.

## Next — auditability and data operations

- [ ] First-class audit event append-only write path for every sensitive mutation.
- [ ] Audit viewer with actor, action, target, request ID, and timestamp filters.
- [ ] CSV product import with dry-run validation and row-level error reports.
- [ ] CSV inventory and report export.
- [ ] Documented automated backup retention deployment examples.
- [ ] Database integration tests covering concurrent inventory mutations and migration compatibility.

## Next — reporting and analytics

- [ ] Inventory aging report.
- [ ] Expiry-risk report by configurable date window.
- [ ] Stock movement history and velocity report.
- [ ] Supplier purchasing totals and lead-time tracking.
- [ ] Warehouse/location valuation breakdown.
- [ ] Replenishment history and recommendation effectiveness metrics.

## Release hardening

- [ ] End-to-end browser tests for authentication, catalog, inventory, purchasing, and reporting.
- [ ] Android instrumentation tests for authentication and critical inventory reads.
- [ ] Accessibility audit across keyboard, screen-reader, contrast, touch-target, and reduced-motion paths.
- [ ] Production restore rehearsal and migration rollback verification.
- [ ] API compatibility/versioning policy wired into CI checks.
- [ ] Reproducible web, server, Android, and extension release artifacts.
- [ ] Resolve every blocker/critical defect before the first stable release.

## Later

- [ ] Notification adapters for low stock and expiring lots.
- [ ] Optional multi-organization tenancy after single-organization workflows are fully hardened.
- [ ] Additional mobile platforms only after Android workflows reach release quality.
