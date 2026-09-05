# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased** until each release is published.

## Unreleased

### v0.3.9 — Report Count & Pagination Semantics

The v0.3.9 milestone adds an additive total-count capability for bounded supplier-performance and warehouse-valuation reports without changing existing report response bodies.

- added an optional `repository.CountedReports` capability for total result counts;
- added PostgreSQL-backed supplier-performance and warehouse-valuation count queries;
- exposed the complete bounded-result count through the `X-Total-Count` HTTP response header;
- reused the same validated reporting period and pagination bounds for page and count queries;
- preserved existing JSON and CSV response bodies and export formats;
- preserved legacy repository implementations when the optional counted capability is unavailable;
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

### v0.3.7 — Reporting Handler Capability Selection

The v0.3.7 milestone makes HTTP supplier-performance and warehouse-valuation handlers consume the additive bounded repository capability when it is available, while preserving a safe legacy fallback.

- added shared HTTP helpers for parsing report offsets and constructing repository-facing bounded queries;
- wired supplier-performance and warehouse-valuation handlers to `repository.BoundedReports` when implemented;
- exposed `offset` handling on both endpoints with truthful pagination metadata;
- rejected nonzero offsets with HTTP 501 when only a legacy repository implementation is available instead of silently returning the wrong page;
- preserved legacy report methods and response bodies for compatible callers;
- added regression coverage for bounded capability selection, legacy fallback, snapshot periods, offsets, metadata, and invalid offset input;
- aligned `/api/v1/meta` with the release version `0.3.7` and added a regression test to prevent version drift;
- kept reporting read-only with no database migration.

### v0.3.6 — Repository Query Execution

The v0.3.6 milestone adds concrete PostgreSQL support for the additive bounded reporting capability while preserving the existing repository interfaces.

- implemented `SupplierPerformanceQuery` and `WarehouseValuationQuery` on the PostgreSQL store;
- validated repository query bounds before execution;
- carried supplier reporting periods into the existing time-window query;
- applied requested offsets and limits without changing the legacy `Reports` interface;
- retained deterministic ordering from the underlying report queries;
- preserved request cancellation through the repository context;
- added no database migration and kept reporting read-only.

### v0.3.5 — Reporting Query Options & Repository Bounds

The v0.3.5 milestone defines the additive contract needed to carry validated reporting periods and pagination bounds toward repository-backed analytics without breaking existing repository implementations.

- added a reusable reporting `Query` value for optional period and pagination bounds;
- added an additive `repository.BoundedReports` capability rather than changing the existing `Reports` interface;
- added unit coverage for time-based and non-time-based query semantics;
- documented v0.3.5 as an API-compatible foundation for storage-side bounded execution;
- preserved existing JSON and CSV response bodies;
- kept reporting read-only with no database migration.

### v0.3.4 — Reporting Export & Completeness

The v0.3.4 milestone standardizes report export metadata and makes completeness metadata conservative and truthful for bounded HTTP results.

- centralized CSV export headers and cache-control semantics;
- exposed generated-at, period, limit, offset, and completeness metadata on supplier and warehouse reports;
- marked bounded reports complete only when fewer items than the requested limit are returned;
- added regression coverage for shared export metadata and warehouse export behavior;
- preserved existing JSON and CSV response bodies and filenames;
- kept reporting read-only with no database migration.

### v0.3.3 — Reporting Repository Integration

The v0.3.3 milestone connects the validated HTTP reporting contract to repository-backed analytics where the existing repository interfaces can support it.

- pass request context through reporting storage calls;
- preserve validated reporting bounds at the HTTP boundary;
- add repository integration documentation and tests;
- add deterministic ordering guidance for report implementations;
- add cancellation-focused regression coverage where supported;
- preserve existing JSON and CSV response bodies;
- keep reporting read-only with no automatic inventory or purchasing mutation.

### v0.3.2 — Reporting API Contract & Version Alignment

The v0.3.2 milestone makes the v0.3.0 reporting primitives usable at the HTTP boundary while preserving existing JSON response bodies.

- added shared reporting request validation for `from`/`to`, `limit`, and `offset`;
- added default 30-day reporting windows with the existing 1–365 day validation range;
- exposed generated-at, period, bounds, and completeness metadata through response headers;
- rejected incomplete or malformed reporting periods and invalid pagination bounds with HTTP 400;
- aligned the `/api/v1/meta` version with the repository `VERSION` marker;
- added regression tests for default requests, explicit periods, invalid bounds, and metadata headers;
- kept the reporting changes additive and read-only with no database migration.

### v0.3.1 — Reporting Metadata Foundation

The v0.3.1 milestone introduced the reusable report metadata contract and release documentation used by the HTTP integration work in v0.3.2.

- added reporting metadata construction for generated time, period, bounds, and completeness;
- added validation-focused metadata coverage;
- added reporting integration documentation and release planning;
- updated the canonical `VERSION` marker to 0.3.1.

### v0.3.0 — Reporting Foundation & Analytics Readiness

The v0.3.0 milestone established reusable reporting primitives for bounded periods, period-over-period comparisons, and result bounds.

- added `internal/reporting` bounded reporting period primitive;
- added period-over-period change calculations with safe zero-baseline behavior;
- added centralized limit/offset validation helpers;
- added unit coverage for period, trend, and bounds behavior;
- added reporting foundation documentation;
- added canonical `VERSION` marker for 0.3.0;
- kept the milestone read-only with no database migration.

### v0.2.10 — Reporting Performance & Operational Insights

The v0.2.10 track focused on faster bounded reporting and additive operational insights while preserving existing transactional and reporting semantics.
