# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased** until each release is published.

## Unreleased

### v0.4.6 — Replenishment Review History Workspace

- add a dedicated web review-history panel to the Reports & Analytics workspace;
- expose immutable replenishment review snapshots through the web API client;
- filter review history by `accepted`, `modified`, `dismissed`, `expired`, or all outcomes;
- display review time, product snapshot, outcome, on-hand quantity, suggested quantity, purchase-order linkage, and reviewer;
- preserve the existing read-only reporting model and session-expiry handling;
- align `/api/v1/meta` and its regression test with release version `0.4.6`.

The web history view consumes the existing `GET /api/v1/replenishment/reviews` endpoint and does not add mutation capabilities to the browser.

### v0.4.5 — StockPilot GitHub Bot

- add a repository-native `StockPilot Bot` GitHub Actions workflow at `.github/workflows/bot.yml`;
- respond to `/stockpilot help`, `/stockpilot status`, and `/stockpilot triage` commands in issue conversations;
- automatically provide a safe bug-triage checklist when `/stockpilot` is used in a newly opened issue;
- keep bot permissions least-privileged with read-only repository contents and write access limited to issue comments;
- avoid checking out or executing untrusted issue/PR content;
- align `VERSION` and `/api/v1/meta` with `0.4.5`;
- add documentation for enabling and using the bot.

The bot is intentionally informational: it does not modify application data, merge pull requests, deploy StockPilot, or bypass the normal review and CI process.

### v0.4.4 — Replenishment Review Traceability

- persist immutable snapshots of reorder recommendations at review time;
- record explicit `accepted`, `modified`, `dismissed`, or `expired` outcomes;
- link accepted/modified reviews to an existing purchase order without auto-submitting it;
- validate that linked purchase orders contain the reviewed product;
- add read/list HTTP endpoints for replenishment review records;
- add an additive PostgreSQL migration with a tested rollback script;
- add domain and PostgreSQL integration coverage;
- align `/api/v1/meta` with the release version `0.4.4`.

The release deliberately keeps purchase-order lifecycle authority unchanged and does not infer recommendation decisions from timestamps or free-form notes.

### v0.4.3 — Report Cursor Pagination & Scalability

Development focus:

- add opaque cursor pagination to the replenishment-readiness JSON report;
- keep cursor ordering aligned with the report's deterministic risk, quantity, SKU, and product ordering;
- reject malformed or incomplete cursors with a bounded `400` response;
- preserve the existing JSON fields and CSV export contract;
- expose optional cursor continuation metadata to the web client;
- add regression coverage for cursor round trips, validation, and ordering semantics.

The v0.4.3 branch will not be tagged until GitHub Actions validates the complete change set on the final release commit.

### v0.4.2 — Replenishment Recommendation Traceability

The v0.4.2 release metadata documented the intended replenishment recommendation traceability milestone. The actual implementation remained a follow-up engineering item and was not treated as complete by v0.4.3.

- documented explicit recommendation-to-purchase-order linkage;
- documented recommendation snapshots and explicit outcomes;
- preserved the known-good v0.4.1 baseline while the larger implementation was validated independently.

### v0.4.1 — Release Hardening

The v0.4.1 milestone is a clean maintenance release focused on release metadata and regression-test alignment after the v0.4.0 reporting work.

- aligned the repository release version with `0.4.1`;
- aligned the `/api/v1/meta` regression test with the release version;
- kept the v0.4.0 reporting contracts and database schema unchanged;
- deferred replenishment recommendation traceability to v0.4.2 so it can be implemented and validated independently;
- kept the release branch free of the previously failing traceability implementation.

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
