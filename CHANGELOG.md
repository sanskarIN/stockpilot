# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased** until each release is published.

## Unreleased

### Next: v0.4.1 — Replenishment Recommendation Traceability

Planned scope:

- add explicit linkage between reviewed reorder suggestions and resulting purchase orders;
- expose recommendation outcome state without inferring causality from timing alone;
- preserve existing report contracts and additive repository capabilities.

### v0.4.0 — Replenishment Performance & Operational Analytics

The v0.4.0 milestone expands reporting with supplier-level replenishment execution metrics based on historical purchase orders and recorded receipts.

- added `repository.ReplenishmentReports` as an additive capability;
- added PostgreSQL-backed replenishment-performance analytics;
- honored explicit reporting `from`/`to` periods and standard pagination bounds;
- added ordered, received, outstanding, fill-rate, on-time, late, and average-lead-time metrics;
- added `GET /api/v1/reports/replenishment-performance`;
- preserved existing JSON/CSV report contracts and legacy repository interfaces;
- added HTTP and PostgreSQL regression coverage for bounds, receipts, fill rate, timeliness, and period handling;
- documented the report's non-causal interpretation limits;
- added no database migration and kept reporting read-only.

### v0.3.9 — Report Count & Pagination Semantics

The v0.3.9 milestone adds an additive total-count capability for bounded supplier-performance and warehouse-valuation reports without changing existing report response bodies.

- added an optional `repository.CountedReports` capability for total result counts;
- added PostgreSQL-backed supplier-performance and warehouse-valuation count queries;
- exposed the complete bounded-result count through the `X-Total-Count` HTTP response header;
- reused the same validated reporting period and pagination bounds for page and count queries;
- preserved existing JSON and CSV response bodies and export formats;
- preserved legacy repository implementations when the optional counted capability is unavailable;
- exposed pagination response headers to approved browser origins through CORS;
- aligned `/api/v1/meta` with the release version `0.3.9`;
- added HTTP regression coverage for count-header behavior and pagination bounds;
- added PostgreSQL integration coverage for report count queries;
- corrected the roadmap to reflect completed supplier and warehouse reporting capabilities;
- kept cursor/streaming pagination as a later optimization rather than changing default offset behavior;
- added no database migration and kept reporting read-only.

### v0.3.8 — SQL-Native Pagination & Query Optimization

The v0.3.8 milestone moves bounded report pagination into PostgreSQL rather than fetching an expanded result set and slicing it in Go.

- applied SQL `LIMIT` and `OFFSET` directly to bounded supplier-performance queries;
- applied SQL `LIMIT` and `OFFSET` directly to bounded warehouse-valuation queries;
- retained deterministic ordering for stable page boundaries;
- preserved the additive `repository.BoundedReports` contract;
- preserved legacy `repository.Reports` behavior and existing JSON/CSV response bodies;
- retained request-context cancellation through PostgreSQL query execution;
- kept warehouse valuation totals independent of page slicing;
- added no database migration;
- kept reporting read-only.
