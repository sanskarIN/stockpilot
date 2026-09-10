# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased** until each release is published.

## Unreleased

### v0.5.5 — Release Metadata Consistency Hardening

- add a validation-only release consistency checker for `VERSION`, `/api/v1/meta`, and the API metadata regression test;
- expose the checker through `make release-check` for local verification;
- run the same consistency gate in GitHub Actions so release metadata drift fails CI;
- align `VERSION` and `/api/v1/meta` with `0.5.5`;
- align the API metadata regression test with `0.5.5`;
- document the release metadata contract and verification workflow;
- keep the release free of database migrations and production-data test dependencies.

The release is intentionally focused on preventing version drift between repository metadata, the running API metadata endpoint, and its regression test. The consistency checker is read-only and does not create tags or modify files.

### v0.5.4 — API Metadata & Reporting Route Hardening

- align `/api/v1/meta` with `0.5.4`;
- align the API metadata regression test with `0.5.4`;
- preserve replenishment-readiness route registration independently of optional reporting-insights wiring;
- retain the authenticated reporting regression coverage introduced in v0.5.3;
- keep the release free of database migrations and production-data test dependencies;
- align `VERSION` with `0.5.4`;
- document the release-hardening milestone.

The release is a focused maintenance patch that protects existing reporting behavior while preventing route registration from becoming accidentally coupled to optional report-repository wiring.

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
