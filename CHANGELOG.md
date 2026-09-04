# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased** until each release is published.

## Unreleased

### Next: v0.3.4 — Reporting Export & Completeness

Planned scope:

- make report completeness reflect actual bounded result behavior;
- standardize CSV export metadata and cache semantics;
- add export regression coverage for reporting filters and pagination;
- keep reporting read-only and backward-compatible.

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
