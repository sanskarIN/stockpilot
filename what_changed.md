# StockPilot — Work Continuity Log

## Current milestone

Phase 27 — v0.1.8 receipt history and export reliability, with release preparation.

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
- v0.1.8 extends exports to authoritative receipt history stored in `stock_movements`.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.1.8 preparation

- [x] Added `domain.ReceiptHistoryRow` as a dedicated flattened receipt-history export model.
- [x] Added `repository.ReceiptHistoryFilter` and `ListReceiptHistory` to keep receipt retrieval separate from normal inventory mutation APIs.
- [x] Added a PostgreSQL receipt-history query backed by `stock_movements` with `movement_type = 'receive'`.
- [x] Joined receipt rows to products, locations, warehouses, and lots for a useful operational export without reconstructing history from current purchase-order counters.
- [x] Added product, warehouse, location, lot, actor, reference, and date-range filtering.
- [x] Added inclusive `from` and exclusive `to` date semantics.
- [x] Added default 500 and maximum 5,000 export bounds with negative-offset normalization.
- [x] Added deterministic ordering by `occurred_at DESC, id DESC`.
- [x] Added `GET /api/v1/inventory/receipts/export.csv`.
- [x] Added deterministic browser download filename and CSV content type.
- [x] Reused the shared formula-safe CSV serializer.
- [x] Added UTC RFC3339 timestamp serialization.
- [x] Added focused tests for bounds, filters, date validation, headers, schema, formula safety, and timestamps.
- [x] Added `docs/RELEASE_NOTES_v0.1.8.md`.
- [x] Added the v0.1.8 `CHANGELOG.md` entry.

## v0.1.8 release gates

- [ ] Run `gofmt`, `go vet ./...`, normal tests, and race-enabled Go tests.
- [ ] Verify every production and test implementation of `repository.Inventory` implements `ListReceiptHistory`.
- [ ] Verify every receipt-history filter reaches the repository contract and SQL predicate.
- [ ] Verify `from` is inclusive and `to` is exclusive, including adjacent date windows.
- [ ] Verify invalid dates and invalid ranges return HTTP 400.
- [ ] Verify negative offsets normalize to zero and limits above 5,000 clamp safely.
- [ ] Verify deterministic receipt ordering and CSV schema.
- [ ] Verify formula safety and UTC timestamp serialization.
- [ ] Verify privileged authorization and denial for unauthorized users.
- [ ] Verify origin/CSRF/request-ID/security middleware behavior where applicable.
- [ ] Run PostgreSQL readiness and receiving smoke tests with representative receipt data.
- [ ] Run Web, Android, browser-companion, and CodeQL checks where configured.
- [ ] Verify no credentials or secrets are present in receipt notes, references, or other exported fields.
- [ ] Create immutable `v0.1.8` tag on the verified commit.
- [ ] Publish the GitHub Release with the prepared notes and keep it pre-release until every gate passes.
- [ ] Perform post-release smoke testing against the published tag/artifacts.

## Publication details

- Version: `v0.1.8`
- Release title: `StockPilot v0.1.8 — Receipt History & Export Reliability`
- Git tag: `v0.1.8`
- Release class: normal pre-1.0 feature release
- Pre-release: yes until every release gate passes
- Release notes: `docs/RELEASE_NOTES_v0.1.8.md`

## Known verification limitation in this workspace

- The repository was updated directly through the connected GitHub integration.
- The connected workspace does not provide a trustworthy local checkout/dependency environment for claiming `go test ./...` as completed.
- Full release verification must therefore be completed by GitHub Actions or a local developer environment before publication is considered verified.
- The v0.1.8 receipt-history endpoint intentionally remains bounded and materializes selected rows; true large-export streaming is still future work.

## Next exact development tasks after v0.1.8

1. Add authorization-focused integration coverage across every export family.
2. Add deterministic streaming readers for large export datasets.
3. Add web download controls with accessible loading, success, and failure feedback.
4. Add export job lifecycle/status endpoints for asynchronous large exports.
5. Add export audit events/observability for sensitive dataset downloads.
6. Add richer expiry-risk classification and operational alerts.
7. Add export retention and operational metrics where appropriate.
8. Begin v0.2.x analytics and operational reporting with aggregate, trend, and exception views.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
