# StockPilot — Work Continuity Log

## Current milestone

Phase 28 — v0.1.9 export privacy and authorization hardening, with release preparation.

## Repository state

- Default branch: `main`.
- v0.1.0-preview.1 release preparation is merged.
- v0.1.1 maintenance-release preparation is merged.
- v0.1.2 CSV serialization foundation is merged.
- The transactional CSV product-import workflow remains in `main` from PR #38.
- v0.1.3 added the first bounded application-level product catalog CSV export.
- v0.1.4 extended the export surface to inventory balances, low-stock data, and reorder suggestions.
- v0.1.5 extended exports to lot inventory and expiry-risk filtering.
- v0.1.6 extended exports to purchase-order lines and current receiving progress.
- v0.1.7 extended exports to append-only audit events.
- v0.1.8 extended exports to authoritative receipt history stored in `stock_movements`.
- v0.1.9 hardens export response privacy and adds regression coverage for export authorization mapping.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.1.9 preparation

- [x] Added `setCSVDownloadHeaders` to centralize CSV download response policy.
- [x] Added `Cache-Control: no-store` to every current CSV export.
- [x] Added `Pragma: no-cache` to every current CSV export.
- [x] Preserved deterministic CSV content types and attachment filenames through the shared helper.
- [x] Applied the shared privacy headers to product, inventory, low-stock, reorder, lot, receipt-history, purchase-order, and audit exports.
- [x] Added focused tests for CSV privacy headers.
- [x] Added regression coverage mapping every current export route to its domain-specific read permission.
- [x] Fixed the export-access test fixture to use the standard URL parser.
- [x] Added `docs/RELEASE_NOTES_v0.1.9.md`.
- [x] Added the v0.1.9 `CHANGELOG.md` entry.

## v0.1.9 release gates

- [ ] Run `gofmt`.
- [ ] Run `go vet ./...`.
- [ ] Run the normal Go test suite.
- [ ] Run race-enabled Go tests.
- [ ] Verify every export route through the authenticated `WithAccess` middleware.
- [ ] Verify unauthorized principals receive HTTP 403 for exports they cannot read.
- [ ] Verify all CSV exports send `Cache-Control: no-store` and `Pragma: no-cache`.
- [ ] Verify deterministic content types and filenames remain unchanged.
- [ ] Verify formula-safe serialization remains unchanged.
- [ ] Run PostgreSQL readiness and export smoke tests.
- [ ] Run Web, Android, browser-companion, and CodeQL checks where configured.
- [ ] Verify no export schema introduces credentials or session secrets.
- [ ] Create immutable `v0.1.9` tag on the verified commit.
- [ ] Publish the GitHub Release with the prepared notes.
- [ ] Perform post-release smoke testing.

## Publication details

- Version: `v0.1.9`
- Release title: `StockPilot v0.1.9 — Export Privacy & Authorization Hardening`
- Git tag: `v0.1.9`
- Release class: normal pre-1.0 maintenance/security-hardening release
- Pre-release: yes until every release gate passes
- Release notes: `docs/RELEASE_NOTES_v0.1.9.md`

## Known verification limitation in this workspace

- The repository was updated directly through the connected GitHub integration.
- The connected workspace does not provide a trustworthy local checkout/dependency environment for claiming `go test ./...` as completed.
- GitHub currently reports no commit-status records/workflow runs that can be used as a full CI verification signal for this milestone.
- Full release verification must therefore be completed by GitHub Actions or a local developer environment before publication is considered verified.
- Large-export cursor/streaming support remains future work.

## Next exact development tasks after v0.1.9

1. Add deterministic cursor/streaming readers for large export datasets.
2. Add web download controls with accessible loading, success, and failure feedback.
3. Add export job lifecycle/status endpoints for asynchronous large exports.
4. Add explicit export audit events and operational observability for sensitive dataset downloads.
5. Add richer expiry-risk classification and operational alerts.
6. Add export retention and operational metrics where appropriate.
7. Begin v0.2.x analytics and operational reporting with aggregate, trend, and exception views.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
