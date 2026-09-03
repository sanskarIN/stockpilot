# StockPilot — Work Continuity Log

## Current milestone

Phase 21 — v0.1.2 CSV export foundation and release preparation.

## Repository state

- Default branch: `main`.
- v0.1.0-preview.1 release preparation is merged.
- v0.1.1 maintenance-release preparation is merged.
- The transactional CSV product-import workflow remains in `main` from PR #38.
- v0.1.2 work now establishes the reusable CSV serialization foundation without prematurely coupling export behavior to a repository or endpoint.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.1.2 preparation

- [x] Added `internal/csvexport/csvexport.go` with standard-library CSV serialization, explicit headers, buffered writes, and flush error reporting.
- [x] Added optional spreadsheet-formula protection for cells beginning with `=`, `+`, `-`, or `@` after leading whitespace.
- [x] Added `internal/csvexport/csvexport_test.go` covering quoting, multiline values, formula safety, invalid headers, and nil writers.
- [x] Added `docs/CSV_EXPORT_DESIGN.md` defining export boundaries, planned datasets, security requirements, deterministic ordering, and resource limits.
- [x] Added `docs/RELEASE_NOTES_v0.1.2.md` with scope, security notes, verification requirements, upgrade notes, and known limitations.
- [x] Recorded the v0.1.2 entry in `CHANGELOG.md`.

## v0.1.2 release gates

- [ ] Confirm the exact v0.1.2 release commit on `main` after all intended changes are merged.
- [ ] Run `gofmt`, `go vet ./...`, normal tests, and race-enabled Go tests.
- [ ] Run web typecheck and production build.
- [ ] Run PostgreSQL migration/readiness smoke testing.
- [ ] Run Android lint/tests/build and release-networking/security checks.
- [ ] Run browser-companion manifest and unit checks.
- [ ] Run configured CodeQL checks for Go and JavaScript/TypeScript.
- [ ] Complete backup/restore, authentication/session, authorization/CSRF, responsive/keyboard, Android device, and browser installation smoke checks.
- [ ] Verify no secrets or credentials are included in release artifacts.
- [ ] Create immutable `v0.1.2` tag on the verified commit.
- [ ] Publish the GitHub Release with the prepared notes and mark it pre-release until the full gate set passes.
- [ ] Perform post-release smoke testing against the published tag/artifacts.

## Publication details

- Version: `v0.1.2`
- Release title: `StockPilot v0.1.2 — CSV Export Foundation`
- Git tag: `v0.1.2`
- Release class: normal pre-1.0 feature/foundation release
- Pre-release: yes until every release gate passes
- Release notes: `docs/RELEASE_NOTES_v0.1.2.md`
- CSV design: `docs/CSV_EXPORT_DESIGN.md`

## Next exact development tasks after v0.1.2

1. Add authorized, bounded product-catalog CSV export.
2. Add inventory-balance and low-stock/reorder CSV exports.
3. Add lot-inventory/expiry-risk CSV export.
4. Add purchase-order/receiving export contracts.
5. Add export-specific audit coverage for sensitive datasets.
6. Add deterministic pagination/streaming behavior for large datasets.
7. Add web download controls and accessible export feedback.
8. Expand Android/browser workflows only after the server-side export contracts are stable.
9. Continue toward v0.2.x analytics and operational reporting.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
