# StockPilot — Work Continuity Log

## Current milestone

Phase 18 — Transactional CSV product import and mainline CI repair.

## Repository state

- Default branch: `main`.
- Active branch: `feat/csv-product-import`.
- Merged mainline includes Android/browser quality gates, reporting/replenishment, catalog management, guided inventory operations, warehouse/location lifecycle management, multi-line purchasing and receiving, lot/expiry-aware receiving, atomic new-lot receiving, purchase-order lifecycle controls, append-only business and authentication auditability with a web viewer, reorder-to-draft assistance, lot inventory visibility, browser camera scanning, Android scanner handoff, browser companion inventory workflow handoff, and CSV product dry-run validation.
- This continuation intentionally preserves focused commits instead of squashing feature history.

## Completed in this continuation

- [x] Added a dedicated `repository.ProductBatchImporter` contract so batch persistence is explicit and cannot silently fall back to sequential non-atomic writes.
- [x] Added PostgreSQL transactional product-batch persistence with repeated domain validation and database uniqueness/foreign-key constraints as the final integrity boundary.
- [x] Added `POST /api/v1/products/import` and registered it in the HTTP router beside the existing dry-run endpoint.
- [x] Added server-side reparse/revalidation on every write request, closing the validation-to-write time-of-check/time-of-use gap.
- [x] Added server-generated product IDs for CSV rows that omit `id`.
- [x] Added a successful `products.imported` audit event containing request ID and batch count without storing CSV contents.
- [x] Added a PostgreSQL integration test proving a duplicate inside a batch rolls the entire transaction back.
- [x] Added the web API client for the write endpoint.
- [x] Reworked the web import panel so validation and writing are separate explicit actions and successful imports refresh the catalog.
- [x] Expanded `docs/CSV_PRODUCT_IMPORT.md` with endpoint, transactional, concurrency, audit, and failure semantics.
- [x] Corrected Go module metadata and checksums to satisfy `go mod tidy` under the repository's Go 1.26 CI environment.
- [x] Removed duplicate HTTP barcode and lot handler declarations exposed by the migration smoke test.
- [x] Restored the reporting and audit HTTP handlers required by the existing optional routes and repository contracts.
- [x] Restored the atomic product batch-import repository contract in the repository package.
- [x] Formatted the new HTTP and PostgreSQL import code/tests to satisfy the repository's `gofmt` quality gate.

## Verification status

The connected environment cannot clone GitHub repositories locally because outbound GitHub DNS/network access is unavailable. GitHub Actions is therefore the authoritative execution environment. The previous CI run exposed formatting plus a migration-smoke compile failure. The compile failure identified duplicate handlers and missing reporting/audit HTTP implementations; those issues have now been repaired on the active branch. Fresh CI and CodeQL runs are queued against the repaired head and must finish successfully before merge.

## Known limitations

- CSV inventory/report export remains pending.
- Advanced analytics remain pending.
- Concurrent inventory database integration and migration compatibility coverage remain pending beyond the product-import rollback case.
- Browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks remain pending.
- A future extension-specific credential model is still required if direct extension API mutations are ever introduced.

## Next exact tasks

1. Finish PR #38 only after fresh CI and CodeQL are green, then merge it to `main`.
2. Add CSV inventory/report export with audit events and bounded streaming output.
3. Add inventory aging, configurable expiry-risk, movement-velocity, supplier, and replenishment analytics.
4. Expand concurrent inventory database integration coverage and migration compatibility tests.
5. Add browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks.
6. Add automated backup retention examples and first stable-release gates.

## Commit discipline

Each functional boundary is intentionally represented by an individual commit where practical: domain rules, repository contracts, persistence, HTTP handlers, tests, client API, UI, routing, styling, fixes, and documentation are kept separately reviewable. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution.
