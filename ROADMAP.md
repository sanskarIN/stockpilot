# StockPilot Roadmap

This roadmap reflects the implemented repository state rather than the original scaffold checklist.

## Completed foundation

- [x] Go domain model for catalog, warehouses, locations, inventory movements, transfers, lots, purchasing, and access control.
- [x] PostgreSQL schema, ordered migrations, and repository implementations.
- [x] Secure HTTP API with health/readiness endpoints, validation, CORS allowlisting, hardening headers, and request IDs.
- [x] Authentication with password verification, opaque sessions, session-token peppering, CSRF protection, and RBAC.
- [x] Administrator bootstrap command.
- [x] React + TypeScript responsive dashboard.
- [x] PWA manifest, service worker, and production registration.
- [x] Native Android application with authenticated dashboard flow and encrypted persisted sessions.
- [x] Manifest V3 browser companion with scoped optional host permissions.
- [x] Backend, frontend, Android, extension, CodeQL, and dependency-update CI baselines.
- [x] Aggregate reorder recommendations including zero-stock products.
- [x] Inventory valuation grouped by currency.
- [x] Exact barcode lookup API for scanner-driven workflows.

## Next — operational workflows

- [ ] Product/category/supplier create and edit screens in the web application.
- [ ] Warehouse/location management screens.
- [ ] Guided stock-in, stock-out, adjustment, and transfer forms with confirmation summaries.
- [ ] Purchase-order creation, editing, approval-state controls, and receiving UI.
- [ ] Barcode/QR camera scanning UI for supported clients, using exact server-side lookup.
- [ ] Lot and expiry receiving workflows with expiry warnings.
- [ ] Reorder suggestion actions that can seed draft purchase orders without bypassing approval rules.

## Next — auditability and data operations

- [ ] First-class audit event schema and append-only audit repository.
- [ ] Audit viewer with actor, action, target, request ID, and timestamp filters.
- [ ] CSV product import with dry-run validation and row-level error reports.
- [ ] CSV inventory and report export.
- [ ] Backup script and documented restore drill with retention hooks.
- [ ] Scheduled-backup deployment examples that keep credentials outside source control.
- [ ] Database integration tests covering concurrent inventory mutations and migrations.

## Next — reporting and analytics

- [ ] Inventory aging report.
- [ ] Expiry-risk report by configurable date window.
- [ ] Stock movement history and velocity report.
- [ ] Supplier purchasing totals and lead-time tracking.
- [ ] Warehouse/location valuation breakdown.
- [ ] Replenishment history and recommendation effectiveness metrics.

## Release hardening

- [ ] End-to-end browser tests for authentication, catalog, inventory, purchasing, and reporting paths.
- [ ] Android instrumentation tests for authentication and critical inventory reads.
- [ ] Accessibility audit across keyboard, screen-reader, contrast, touch-target, and reduced-motion paths.
- [ ] Production restore rehearsal and migration rollback verification.
- [ ] API compatibility/versioning policy and generated endpoint reference.
- [ ] Release checklist with reproducible web, server, Android, and extension artifacts.
- [ ] Resolve every blocker/critical defect before the first stable release.

## Later

- [ ] Notification adapters for low stock and expiring lots.
- [ ] Optional multi-organization tenancy after single-organization workflows are fully hardened.
- [ ] Additional mobile platforms only after the API and Android workflows reach release quality.
