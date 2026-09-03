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
- [x] Reformatted `internal/httpapi/api.go` after the CI formatting gate identified the router/middleware file as non-gofmt-compliant.
- [x] Improved the CI formatting step so a future formatting failure prints the exact non-gofmt-formatted files instead of failing silently.
- [x] Preserved the existing PostgreSQL Go module cache configuration and existing Node action version while making the CI diagnostic change.
- [x] Re-ran the failed Go-quality workflow job to validate the updated CI diagnostics.

## Verification status

The connected environment cannot clone GitHub repositories locally because outbound GitHub DNS/network access is unavailable. GitHub Actions is therefore the authoritative execution environment. The previous CI run for PR #38 failed at the Go formatting gate while the PostgreSQL migration smoke test and Web quality job passed, and CodeQL passed. The CI formatting step has now been changed to print the exact failing file list, and a workflow retry has been initiated. The new head must receive a fresh green CI result before PR #38 is merged.

## Known limitations

- CSV inventory/report export remains pending.
- Advanced analytics remain pending.
- Concurrent inventory database integration and migration compatibility coverage remain pending beyond the product-import rollback case.
- Browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks remain pending.
- A future extension-specific credential model is still required if direct extension API mutations are ever introduced.

## Next exact tasks

1. Verify fresh PR #38 CI and CodeQL against the current head; merge only when the required checks are green.
2. If formatting still fails, use the diagnostic file list and make focused formatting-only commits.
3. Add CSV inventory/report export with audit events and bounded streaming output.
4. Add inventory aging, configurable expiry-risk, movement-velocity, supplier, and replenishment analytics.
5. Expand concurrent inventory database integration coverage and migration compatibility tests.
6. Add browser E2E, Android instrumentation, accessibility, restore, and release-artifact checks.
7. Add automated backup retention examples and first stable-release gates.

## Commit discipline

Each functional boundary is intentionally represented by an individual commit where practical: domain rules, repository contracts, persistence, HTTP handlers, tests, client API, UI, routing, styling, fixes, CI, and documentation are kept separately reviewable. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution.
