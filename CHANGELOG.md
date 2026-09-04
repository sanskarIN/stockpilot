# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased** until each release is published.

## Unreleased

### Next: v0.3.1 — Reporting API Integration & Metadata

The v0.3.1 milestone builds on the v0.3.0 reporting foundation and integrates bounded period and pagination metadata into the public reporting layer without changing transactional semantics.

Planned scope:

- integrate shared reporting period validation into report endpoints;
- expose generated-at timestamps and applied reporting bounds;
- expose complete/partial result metadata where applicable;
- standardize bounded limit and offset handling;
- preserve deterministic report ordering;
- propagate request cancellation and timeouts;
- keep report responses backward-compatible and additive;
- maintain formula-safe, non-cacheable CSV exports;
- add regression coverage for metadata, bounds, cancellation, ordering, and malformed parameters;
- keep reporting read-only with no automatic inventory or purchasing mutation;
- avoid database migrations unless implementation evidence requires them.

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
