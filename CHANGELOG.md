# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased**.

## Unreleased

Future work will continue here after v0.1.2.

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
- Full repository release gates remain required before publication, including Go, PostgreSQL, Web, Android, extension, and CodeQL checks where configured.

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
