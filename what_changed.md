# StockPilot — Work Continuity Log

## Current milestone

Phase 30 — v0.2.0 Reports & Analytics foundation, with release preparation.

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
- v0.1.10 added an explicit audit event for authenticated CSV export requests.
- v0.2.0 adds the dedicated Reports & Analytics workspace and inventory valuation CSV export.
- Focused, reviewable commits are preferred over meaningless commits solely to increase the commit count.

## Completed for v0.2.0 preparation

- [x] Added the inventory valuation CSV export endpoint.
- [x] Bounded valuation export rows to a safe default and maximum.
- [x] Reused formula-safe CSV serialization.
- [x] Reused `no-store`/`no-cache` export response policy.
- [x] Added focused valuation export regression tests.
- [x] Added typed web reporting summary contracts.
- [x] Added web API methods for report overview, inventory summary, purchasing summary, and valuation.
- [x] Added a dedicated Reports & Analytics web workspace.
- [x] Added inventory and purchasing operational metrics to the reports screen.
- [x] Added currency-grouped valuation totals and product valuation breakdown.
- [x] Added dashboard and navigation entry points to Reports.
- [x] Added responsive Reports workspace styling.
- [x] Added `docs/RELEASE_NOTES_v0.2.0.md`.
- [x] Added the v0.2.0 `CHANGELOG.md` entry.
- [x] Updated the roadmap for the post-v0.2.0 reporting backlog.

## v0.2.0 release gates

- [ ] Run `gofmt`.
- [ ] Run `go vet ./...`.
- [ ] Run the normal Go test suite.
- [ ] Run race-enabled Go tests.
- [ ] Verify valuation CSV bounds and serialization.
- [ ] Verify reporting endpoints remain protected by `PermissionReportsRead`.
- [ ] Verify valuation CSV requests generate `export.csv.requested` audit events after authorization.
- [ ] Verify valuation exports retain `Cache-Control: no-store` and `Pragma: no-cache`.
- [ ] Verify formula-safe serialization remains enabled.
- [ ] Run PostgreSQL reporting integration tests.
- [ ] Run Web, Android, browser-companion, and CodeQL checks where configured.
- [ ] Verify Reports workspace loading, error, refresh, and session-expiry paths.
- [ ] Verify responsive and keyboard-accessible report controls.
- [ ] Create immutable `v0.2.0` tag on the verified commit.
- [ ] Publish the GitHub Release with the prepared notes.
- [ ] Perform post-release reporting and export smoke testing.

## Publication details

- Version: `v0.2.0`
- Release title: `StockPilot v0.2.0 — Reports & Analytics`
- Git tag: `v0.2.0`
- Release class: pre-1.0 feature milestone
- Pre-release: yes until every release gate passes
- Release notes: `docs/RELEASE_NOTES_v0.2.0.md`

## Known verification limitation in this workspace

- The repository was updated directly through the connected GitHub integration.
- The connected workspace does not provide a trustworthy local checkout/dependency environment for claiming `go test ./...` as completed.
- GitHub currently does not expose a full CI verification signal for this milestone through the available integration.
- Full release verification must therefore be completed by GitHub Actions or a local developer environment before publication is considered verified.
- Large-report cursor/streaming support remains future work.

## Next exact development tasks after v0.2.0

1. Add inventory aging reports with deterministic age buckets.
2. Add configurable expiry-risk reports and operational alert thresholds.
3. Add movement history and velocity analytics.
4. Add supplier purchasing totals and lead-time measurements.
5. Add warehouse/location valuation breakdowns.
6. Add replenishment effectiveness metrics.
7. Add cursor/streaming report readers and large-export lifecycle endpoints.
8. Expand web end-to-end, Android instrumentation, accessibility, restore, and compatibility gates.

## Commit discipline

Each functional boundary remains separately reviewable where practical. Repository commits continue to use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution. Meaningful focused commits are preferred over artificial commit-count inflation.
