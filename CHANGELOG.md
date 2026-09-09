# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased** until each release is published.

## Unreleased

### v0.5.3 — Authenticated Reporting Workspace Hardening

- add deterministic authenticated Playwright coverage for the Reports & Analytics workspace;
- fixture the complete reporting data surface, including overview, valuation, inventory aging, movement history, supplier performance, replenishment readiness, and replenishment review history;
- assert representative rendered report metrics and report sections through the authenticated UI;
- keep reporting E2E fixtures synthetic and independent of production credentials/data;
- align `VERSION` with `0.5.3`;
- preserve the server-authoritative reporting and read-only browser model;
- document the completed reporting regression milestone.

The release expands browser-level regression protection across the read-only reporting surface without changing production reporting contracts or mutation behavior.

### v0.5.2 — Purchase Order Lifecycle Regression Hardening

- fix the ambiguous Playwright `Unit` locator that caused authenticated product-creation E2E failure under strict mode;
- add deterministic authenticated purchase-order create coverage;
- add authenticated purchase-order submit/status mutation coverage with request-payload assertions;
- add authenticated purchase-order receiving coverage with request-payload assertions;
- keep all purchase-order E2E fixtures synthetic and independent of production credentials/data;
- align `VERSION` and `/api/v1/meta` with `0.5.2`;
- align the API metadata regression test with `0.5.2`;
- document the completed purchase-order lifecycle regression milestone.

The release hardens the existing purchase-order workflow without changing the server-authoritative lifecycle rules or production data paths.

### v0.5.1 — Authenticated Mutation Regression Hardening

- add deterministic Playwright coverage for authenticated catalog product creation;
- add deterministic Playwright coverage for authenticated stock-in inventory mutations;
- assert the browser sends the expected product and inventory mutation payloads;
- keep mutation E2E fixtures entirely synthetic and isolated from production credentials/data;
- align `VERSION` and `/api/v1/meta` with `0.5.1`;
- align the API metadata regression test with `0.5.1`;
- document the mutation-path hardening milestone.

The new browser coverage validates the existing mutation workflows without changing their server-authoritative validation model.

### v0.5.0 — Authenticated Workflow Reliability

- add authenticated Playwright browser coverage for the primary StockPilot workspaces;
- add deterministic navigation coverage for Products, Inventory, Purchase Orders, Warehouses, Lot Inventory, and Audit History;
- add regression coverage confirming each primary workspace can return to the Inventory Overview;
- keep browser fixtures deterministic and synthetic;
- align application release metadata with `0.5.0`;
- update the project roadmap and release documentation.

### v0.4.9 — Authenticated Browser Reliability

- add deterministic Playwright Chromium coverage for the signed-out login contract;
- add an authenticated dashboard-shell browser fixture using an in-browser API test double;
- verify the authenticated session boundary without requiring live credentials or production data;
- verify dashboard navigation renders with deterministic empty catalog, inventory, reorder, order, and valuation responses;
- add a dedicated Browser E2E CI gate alongside Go, Web, and PostgreSQL validation;
- align `VERSION` and `/api/v1/meta` with release version `0.4.9`;
- align the API metadata regression test with `0.4.9`;
- document that deeper authenticated workflow coverage remains the next incremental step.

The browser tests use synthetic fixtures only. They do not authenticate against a production account, persist credentials, or modify real application data.

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
