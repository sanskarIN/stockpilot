# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased** until each release is published.

## Unreleased

### Next: v0.2.10 — Reporting Performance & Operational Insights

The next planned milestone focuses on faster bounded reporting and additive operational insights while preserving existing transactional and reporting semantics.

Planned scope:

- reuse bounded pagination and continuation primitives from the large-report scalability work;
- deterministic ordering and explicit limits for new analytical queries;
- reduced repeated database work for common report views where safe and measurable;
- preserved timeout and cancellation propagation through report execution;
- additive inventory, purchasing, and replenishment trend summaries;
- period-over-period changes with currencies kept separate;
- source-of-truth identifiers for report traceability;
- clear distinction between observed historical metrics and advisory calculations;
- defensive validation for date windows, limits, pagination inputs, and unsupported combinations;
- stable backward-compatible response shapes and actionable public errors;
- formula-safe, non-cacheable, bounded CSV exports;
- Reports & Analytics period/bound context plus loading, empty, error, retry, and accessible table states;
- regression coverage for pagination, bounds, cancellation, ordering, duplicate prevention, exports, and security.

v0.2.10 remains read-only. No automatic purchasing or inventory mutation is planned, and schema migrations require evidence and compatibility review.

### Next: v0.2.9 — Operational Readiness & Reporting Reliability

v0.2.9 is the operational hardening milestone preceding v0.2.10. Its planned scope is:

- explicit report query budgets and bounded execution behavior where supported;
- standardized report metadata for generated-at timestamps, applied bounds, and complete/partial result state;
- actionable report failures without exposing database internals;
- deterministic ordering for repeatable report consumption;
- structured request/report logging with stable event names and duration measurements;
- safe slow-report diagnostics without credentials, session tokens, or sensitive request bodies;
- timeout and cancellation propagation across read-heavy/report endpoints;
- defensive validation for query parameters, date ranges, limits, and continuation inputs;
- bounded large-export behavior with formula-safe, non-cacheable CSV output;
- resilient Reports & Analytics loading, empty, error, retry, and accessibility states;
- operational documentation for health, readiness, and troubleshooting;
- regression coverage for timeout, cancellation, bounds, logging safety, and export memory behavior.

v0.2.9 remains read-only with no automatic purchasing or inventory mutation. Database migrations require evidence and compatibility review. Stable publication remains verification-gated.

### v0.2.8 — Large-Report Scalability & Purchasing Trends

Implementation is still required before v0.2.8 can be represented as a completed stable release. The release preparation scope remains:

- reusable bounded pagination primitives for large report datasets;
- additive cursor/keyset pagination or streaming without silently changing existing endpoint semantics;
- opaque continuation tokens with deterministic ordering;
- bounded memory use, page sizes, date windows, and query work;
- supplier/product purchasing trend series built from authoritative purchase-order and purchase-order-line data;
- ordered, received, open-unit, and purchasing-value metrics with currencies kept separate;
- safe date-window and row bounds for trend endpoints;
- CSV export that remains formula-safe and non-cacheable;
- Reports & Analytics trend views with clear empty/loading/error states and accessible table alternatives;
- regression coverage for pagination, ordering, bounds, continuation, duplicate prevention, and export behavior.

No feature in this preparation should silently change transactional inventory or purchasing records. A schema migration will be avoided unless implementation evidence proves the existing schema insufficient.

### v0.2.7 — Replenishment Readiness

Completed implementation scope:

- added advisory replenishment-readiness reporting derived from existing reorder suggestions and stock-movement history;
- added on-hand, reorder point, reorder quantity, target stock, and suggested quantity visibility;
- added recent outbound units and average daily outbound velocity;
- added estimated days of cover when outbound velocity is positive;
- added deterministic `out_of_stock`, `critical`, `reorder`, `watch`, and `healthy` classifications;
- added bounded authenticated JSON reporting and formula-safe CSV export;
- added the `GET /api/v1/reports/replenishment-readiness` endpoint;
- updated API metadata to `0.2.7`;
- added domain contract and regression coverage;
- kept the report read-only with no purchase-order creation and no inventory mutation;
- introduced no database migration.

Stable publication still requires the repository's applicable verification gates to be confirmed green.
