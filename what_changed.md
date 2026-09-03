# StockPilot — Work Continuity Log

## Current milestone

Phase 19 — First public preview release preparation (`v0.1.0-preview.1`).

## Repository state

- Default branch: `main`.
- Release branch: `release/v0.1.0-preview.1`.
- Mainline now includes the transactional CSV product-import workflow merged through PR #38.
- The merged milestone has green Go quality, PostgreSQL migration smoke-test, Web quality, and CodeQL checks.
- This continuation intentionally preserves focused commits instead of squashing feature history.

## Completed in the previous milestone

- [x] Added a dedicated `repository.ProductBatchImporter` contract so batch persistence is explicit and cannot silently fall back to sequential non-atomic writes.
- [x] Added PostgreSQL transactional product-batch persistence with repeated domain validation and database uniqueness/foreign-key constraints as the final integrity boundary.
- [x] Added `POST /api/v1/products/import` with server-side reparse/revalidation, generated IDs, and batch-level audit events.
- [x] Added PostgreSQL rollback coverage for duplicate conflicts inside an import batch.
- [x] Added the web API client and explicit CSV import UI workflow.
- [x] Expanded CSV import documentation and project changelog/roadmap tracking.
- [x] Repaired Go formatting and CI diagnostics after the branch exposed formatting drift.
- [x] Removed duplicate HTTP test fake methods and restored their shared behavior in the central test fake.
- [x] Updated Go 1.26-compatible test fixtures and parser fixtures.
- [x] Added a required `productId` guard to the lot-listing endpoint.
- [x] Verified the final branch with green Go quality, PostgreSQL migration smoke-test, Web quality, and CodeQL workflows.
- [x] Merged PR #38 into `main` as commit `679bda9c281738f8dc56e38fb51fcaac059f7607`.

## Release preparation

- [ ] Decide the first public preview tag and create it on the verified `main` release commit.
- [ ] Publish release notes covering CSV import, transactional integrity, audit behavior, migrations, and operational requirements.
- [ ] Attach reproducible server/web/Android/browser artifacts and SHA-256 checksums when available.
- [ ] Complete the non-automated release checklist items: backup/restore drill, authentication smoke test, responsive/keyboard review, Android device smoke test, and browser-companion installation smoke test.
- [ ] Complete the post-release smoke test.

## Verification status

The merged PR #38 was verified by GitHub Actions before merge. The final run passed Go formatting, vet, race-enabled tests, server build, PostgreSQL migration readiness, Web typecheck/build, and both Go and JavaScript/TypeScript CodeQL analyses. The connected environment still cannot provide a local GitHub clone, so GitHub Actions remains the authoritative CI execution environment.

## Known limitations

- GitHub release/tag creation is not exposed by the current connected GitHub write surface, so publication of the actual GitHub Release must be completed from the repository UI/API with the verified release commit.
- CSV inventory/report export remains the next product-development milestone after the preview publication.
- Advanced analytics, broader concurrent-inventory integration coverage, migration compatibility coverage, browser E2E, Android instrumentation, accessibility, restore automation, and release-artifact automation remain pending.

## Next exact tasks

1. Publish the first public preview from verified main (`v0.1.0-preview.1` or the final approved preview tag).
2. Start the v0.2.0 development branch for CSV inventory/report export.
3. Add bounded/streaming CSV exports for inventory, movements, purchasing, and audit/report views with authorization and audit events.
4. Add inventory aging, configurable expiry-risk, movement velocity, supplier, and replenishment analytics.
5. Expand concurrent inventory database integration and migration compatibility tests.
6. Add browser E2E, Android instrumentation, accessibility, restore verification, artifact checks, and backup-retention examples.
7. Prepare the stable-release gates for a later `v1.0.0` release.

## Commit discipline

Each functional boundary is intentionally represented by an individual commit where practical: domain rules, repository contracts, persistence, HTTP handlers, tests, client API, UI, routing, styling, fixes, CI, and documentation are kept separately reviewable. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution.