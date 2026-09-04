# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased**.

## Unreleased

Future work will continue here after v0.1.6.

## 0.1.6 — 2026-09-04

### Added

- Added bounded purchase-order CSV export at `GET /api/v1/orders/export.csv`.
- Added optional purchase-order status filtering for `draft`, `ordered`, `partially_received`, `received`, and `cancelled`.
- Added line-level ordered, received, remaining, unit-cost, and line-total fields to the export.
- Added a dedicated `PurchaseOrderExportRow` domain model and repository export operation.
- Added a deterministic PostgreSQL export query that selects purchase orders first and then flattens their lines.
- Added focused HTTP coverage for export bounds, status validation, headers, schema, formula-safe values, totals, and timestamps.

### Changed

- Purchase-order exports reuse the shared formula-safe CSV serializer introduced in v0.1.2.
- Purchase-order selection is deterministic by `created_at DESC, id DESC` and line ordering is deterministic by line ID.
- Export timestamps are normalized to UTC RFC 3339 values.
- Browser downloads use the deterministic filename `stockpilot-purchase-orders.csv`.
- Receiving visibility is exposed through the authoritative `received` and `remaining` line quantities already maintained by the existing receiving transaction flow.

### Security and operations

- Purchase-order export is read-only and remains behind the existing HTTP middleware and application authorization/session controls.
- Application-level export bounds cap requests at 5,000 selected purchase orders; repository bounds enforce the same ceiling.
- Export serialization protects spreadsheet formula-like values.
- Export schemas exclude credentials, passwords, session secrets, and payment information.

### Verification

- Added unit coverage for purchase-order export normalization and CSV contracts.
- Full release verification remains required for Go, PostgreSQL, Web, Android, extension, authentication/session, authorization/CSRF, and CodeQL gates where configured.

## 0.1.5 — 2026-09-03

### Added

- Added bounded lot-inventory CSV export at `GET /api/v1/inventory/lots/export.csv`.
- Added optional product, warehouse, location, and lot filters to the export contract.
- Added the inclusive `expiringBy` date filter using `YYYY-MM-DD` input.
- Added focused HTTP coverage for bounds, date parsing, headers, schema, timestamps, expiry filtering, and formula-safe values.

### Changed

- Lot inventory exports reuse the shared formula-safe CSV serializer introduced in v0.1.2.
- Export rows retain deterministic PostgreSQL ordering by expiry date, product name, lot number, and location name.
- Expiry timestamps are normalized to UTC RFC 3339 values.
- Browser downloads use the deterministic filename `stockpilot-lot-inventory.csv`.

### Security and operations

- Lot inventory export is read-only and remains behind the existing HTTP middleware and application authorization/session controls.
- Export requests are bounded at the application layer and the repository retains its own safety cap.
- Export schemas exclude credentials, passwords, session secrets, and payment information.

### Verification

- Added unit coverage for lot export normalization and CSV contracts.
- Full release verification remains required for Go, PostgreSQL, Web, Android, extension, authentication/session, authorization/CSRF, and CodeQL gates where configured.

## 0.1.4 — 2026-09-03

### Added

- Added bounded inventory-balance CSV export at `GET /api/v1/inventory/export.csv`.
- Added low-stock CSV export at `GET /api/v1/inventory/low-stock/export.csv`.
- Added reorder-suggestions CSV export at `GET /api/v1/inventory/reorder-suggestions/export.csv`.
- Added a repository-level bounded `ListBalances` operation for deterministic inventory export pagination.
- Added focused HTTP coverage for export bounds, headers, formula safety, timestamps, low-stock output, and reorder-suggestion output.

### Changed

- Inventory exports reuse the shared formula-safe CSV serializer introduced in v0.1.2.
- Inventory-balance export ordering is deterministic by product, location, and lot.
- Export timestamps are normalized to UTC RFC 3339 values.
- Browser downloads use deterministic CSV filenames and content types.

### Security and operations

- Export endpoints are read-only and remain behind the existing HTTP middleware and application authorization/session controls.
- Inventory-balance export has a hard application limit of 5,000 rows per request.
- Low-stock and reorder-suggestion queries retain their existing repository-side safety bounds.
- Export schemas exclude credentials, passwords, session secrets, and payment information.

### Verification

- Added unit coverage for inventory export normalization and CSV contracts.
- Full release verification remains required for Go, PostgreSQL, Web, Android, extension, authentication/session, authorization/CSRF, and CodeQL gates where configured.

## 0.1.3 — 2026-09-03

### Added

- Added bounded product catalog CSV export at `GET /api/v1/products/export.csv`.
- Reused existing catalog filters for search, category, supplier, active-only selection, limit, and offset.
- Added deterministic CSV column ordering, UTC timestamp serialization, formula-safe serialization, and browser download headers.
- Added focused export-bound tests and dedicated v0.1.3 release notes.

### Security and operations

- Product export remains read-only and uses bounded application-level pagination.
- Export serialization protects spreadsheet formula-like values.
- Existing request-ID, origin, authentication/session, authorization, and security-header controls remain applicable.

### Verification

- Full repository release gates remain required for publication, including Go, PostgreSQL, Web, Android, extension, and CodeQL checks where configured.

## 0.1.2 — 2026-09-03

### Added

- Added the reusable `internal/csvexport` package for deterministic, RFC 4180-compatible CSV serialization.
- Added optional spreadsheet-formula protection for cells beginning with `=`, `+`, `-`, or `@` after leading whitespace.
- Added focused unit coverage for quoting, formula safety, validation, and writer errors.
- Added `docs/CSV_EXPORT_DESIGN.md` describing the planned export contracts, authorization boundaries, resource limits, and security requirements.

### Security and operations

- Export design explicitly keeps authorization, filtering, ordering, and resource limits outside the serializer.
- Export guidance prohibits credentials, session values, secrets, and other sensitive authentication material from downloadable datasets.
- Large exports are planned to use bounded/streaming or asynchronous approaches rather than unbounded in-memory buffering.

### Verification

- CSV serialization tests cover commas, newlines, formula-like values, invalid headers, and nil-writer handling.
- Full repository release gates remain required for publication, including Go, PostgreSQL, Web, Android, extension, and CodeQL checks where configured.

## 0.1.1 — 2026-09-03

### Changed

- Clarified the v0.1.x maintenance-release policy and verification gates.
- Added a dedicated release runbook covering source freeze, automated/manual verification, tagging, artifacts, post-release checks, and rollback.
- Documented the v0.1.1 maintenance-release scope without introducing a large feature payload.

### Security and operations

- Reaffirmed backup/restore verification before release publication.
- Reaffirmed HTTPS enforcement, encrypted Android session storage, scoped browser-companion permissions, security headers, CORS review, backup retention, and request-ID logging as release gates.
- Documented the rule that a published release tag must not be rewritten; corrective defects should ship as a subsequent patch release.

### Verification

- v0.1.1 publication requires green Go quality, PostgreSQL migration smoke-test, Web quality, and CodeQL checks for the exact release commit.
- Manual release gates remain required for backup/restore, authentication/session behavior, responsive/keyboard review, Android device smoke testing, browser-companion installation/handoff, rollback rehearsal, and post-release smoke testing.

## 0.1.0-preview.1 — 2026-09-03

### Added

- Product, category, supplier, warehouse, location, lot, inventory movement, transfer, and purchase-order foundations.
- Session-based authentication, RBAC, CSRF protection, and administrator bootstrap tooling.
- PostgreSQL persistence with ordered migrations and transactional inventory operations.
- Responsive React + TypeScript web dashboard and installable PWA behavior.
- Native Android client with encrypted session storage, authenticated API access, dark mode, and release TLS enforcement.
- Manifest V3 browser companion with scoped optional host permissions and server health/launcher flow.
- CodeQL and Dependabot automation.
- Aggregate reorder recommendations, inventory valuation, exact barcode lookup, reporting coverage, catalog management, guided inventory operations, warehouse/location lifecycle, multi-line purchasing, lot/expiry receiving, append-only business auditability, lot inventory visibility, browser/Android barcode scanning, and companion scan handoff.
- Authentication/session audit-event definitions and regression coverage.
- Companion workflow choices for product lookup and direct inventory-operation handoff.
- Product CSV dry-run validation with field parsing, duplicate detection, existing-SKU checks, and category/supplier reference validation.
- Transactional CSV product import with server-side revalidation, generated IDs for rows without IDs, batch-level audit events, and complete-batch rollback on persistence failure.

### Changed

- Sensitive successful business mutations and authentication/account lifecycle events use the append-only audit stream with request-ID correlation.
- Authentication audit metadata is deliberately coarse for failed login/session events and excludes credential material.
- Scanned companion barcodes can preselect an authenticated product in the web inventory workflow without copying session credentials into the extension.
- Catalog users with write permission can launch the CSV validation panel without bypassing existing product-management permissions.
- CSV product import separates dry-run validation from an explicit write action; the write request reparses and revalidates the file before one atomic database transaction.
- Lot listing now requires an explicit `productId` filter to keep the endpoint bounded to a product-scoped query.

### Security

- HTTP mutation requests retain explicit CSRF confirmation requirements.
- Authorization remains server-side and role-aware.
- Production Android networking requires TLS.
- Browser companion permissions remain scoped to the configured StockPilot origin.
- Audit history remains read-only from the web client and audit storage exposes append/list operations without an update/delete API.
- Authentication audit events never store passwords, raw session tokens, cookie values, or credential-bearing metadata.
- Scanner flows do not copy StockPilot session credentials into scanner state.
- Companion inventory handoff is navigation-only; stock mutations still require the authenticated web application and explicit user submission.
- CSV dry-run validation does not write imported rows to the database, and the write endpoint never trusts a previous dry-run result as proof of current database state.
- Product import relies on database uniqueness and foreign-key constraints as the final integrity boundary and rolls back the complete batch on constraint failure.

### Verification

- Go module tidy verification passed.
- `gofmt` formatting gate passed.
- `go vet ./...` passed.
- Race-enabled Go tests and server build passed.
- PostgreSQL migration/readiness smoke test passed.
- Web typecheck and production build passed.
- CodeQL passed for Go and JavaScript/TypeScript.

## Release discipline

Dependency upgrades and generated artifacts are reviewed independently from product changes. Stable releases require the checks in `docs/RELEASE_CHECKLIST.md` and the restore procedure in `docs/RESTORE_DRILL.md`.
