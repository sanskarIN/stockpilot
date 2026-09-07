# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased** until each release is published.

## Unreleased

### v0.4.8 — Browser Accessibility & End-to-End Reliability

- add a Playwright-based Chromium browser test foundation for the web workspace;
- add signed-out login regression coverage with deterministic `/api/v1/auth/me` mocking;
- verify accessible form labels, autocomplete semantics, and disabled submit behavior;
- verify keyboard traversal through the primary sign-in controls;
- run browser regression checks as a dedicated GitHub Actions CI job;
- install only the Chromium browser required by the v0.4.8 browser suite;
- align `VERSION` and `/api/v1/meta` with release version `0.4.8`;
- align the API metadata regression test with `0.4.8`.

The browser suite is intentionally incremental: it establishes reliable coverage for the public signed-out workflow while leaving authenticated end-to-end scenarios for subsequent releases that can provide a safe test account/bootstrap fixture.

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
