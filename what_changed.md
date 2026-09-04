# StockPilot — Work Continuity Log

## Current milestone

Phase 29 — v0.1.10 export audit trail, with release preparation.

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
- v0.1.9 hardened export response privacy and added regression coverage for export authorization mapping.
- v0.1.10 adds an explicit audit event for authenticated CSV export requests.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.1.10 preparation

- [x] Added authenticated CSV export request detection in the access layer.
- [x] Added `export.csv.requested` audit events after the existing permission check succeeds.
- [x] Captured the authenticated actor ID on export audit events.
- [x] Preserved the existing request ID correlation on export audit events.
- [x] Identified exports by route path without recording query parameters or CSV contents.
- [x] Added regression coverage for GET/HEAD/POST CSV request classification.
- [x] Added regression coverage for actor, entity, request-ID, and method fields on export audit events.
- [x] Added `docs/RELEASE_NOTES_v0.1.10.md`.
- [x] Added the v0.1.10 `CHANGELOG.md` entry.

## v0.1.10 release gates

- [ ] Run `gofmt`.
- [ ] Run `go vet ./...`.
- [ ] Run the normal Go test suite.
- [ ] Run race-enabled Go tests.
- [ ] Verify authenticated CSV requests create `export.csv.requested` audit events.
- [ ] Verify unauthenticated and unauthorized export requests remain rejected before export handling.
- [ ] Verify actor and request IDs are correlated without storing query strings, credentials, session secrets, or dataset contents.
- [ ] Verify all CSV exports retain `Cache-Control: no-store` and `Pragma: no-cache`.
- [ ] Verify deterministic content types and filenames remain unchanged.
- [ ] Verify formula-safe serialization remains unchanged.
- [ ] Run PostgreSQL readiness and export smoke tests.
- [ ] Run Web, Android, browser-companion, and CodeQL checks where configured.
- [ ] Create immutable `v0.1.10` tag on the verified commit.
- [ ] Publish the GitHub Release with the prepared notes.
- [ ] Perform post-release export/audit smoke testing.

## Publication details

- Version: `v0.1.10`
- Release title: `StockPilot v0.1.10 — Export Audit Trail`
- Git tag: `v0.1.10`
- Release class: normal pre-1.0 maintenance/security-observability release
- Pre-release: yes until every release gate passes
- Release notes: `docs/RELEASE_NOTES_v0.1.10.md`

## Known verification limitation in this workspace

- The repository was updated directly through the connected GitHub integration.
- The connected workspace does not provide a trustworthy local checkout/dependency environment for claiming `go test ./...` as completed.
- GitHub currently reports no commit-status records/workflow runs that can be used as a full CI verification signal for this milestone.
- Full release verification must therefore be completed by GitHub Actions or a local developer environment before publication is considered verified.
- Large-export cursor/streaming support remains future work.

## Next exact development tasks after v0.1.10

1. Add deterministic cursor/streaming readers for large export datasets.
2. Add web download controls with accessible loading, success, and failure feedback.
3. Add export job lifecycle/status endpoints for asynchronous large exports.
4. Add richer expiry-risk classification and operational alerts.
5. Add export retention and operational metrics where appropriate.
6. Complete reporting/analytics foundations for inventory aging, expiry risk, movement velocity, supplier totals, and valuation breakdowns.
7. Begin the broader v0.2.x analytics and operational reporting milestone after the export foundation is fully hardened.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
