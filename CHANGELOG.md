# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased** until each release is published.

## Unreleased

### v0.4.7 — Web Review History Reliability & Accessibility

- prevent stale asynchronous review-history responses from overwriting a newer outcome filter;
- keep the last valid table visible while a refresh is in progress instead of rendering duplicate loading rows;
- add explicit accessible labeling for the outcome filter and refresh action;
- expose loading/result state through a polite status announcement;
- add semantic table captions and column scopes for assistive technology;
- align `VERSION` and `/api/v1/meta` with release version `0.4.7`.

The v0.4.7 milestone is intentionally incremental: it hardens the v0.4.6 review-history workspace while broader browser end-to-end and accessibility coverage remains on the release-hardening roadmap.

### v0.4.6 — Replenishment Review History Workspace

- add a dedicated web review-history panel to the Reports & Analytics workspace;
- expose immutable replenishment review snapshots through the web API client;
- filter review history by `accepted`, `modified`, `dismissed`, `expired`, or all outcomes;
- display review time, product snapshot, outcome, on-hand quantity, suggested quantity, purchase-order linkage, and reviewer;
- preserve the existing read-only reporting model and session-expiry handling;
- align `/api/v1/meta` and its regression test with release version `0.4.6`.

The web history view consumes the existing `GET /api/v1/replenishment/reviews` endpoint and does not add mutation capabilities to the browser.

### v0.4.4 — Replenishment Review Traceability

- advance the repository release metadata from `0.5.9` to `0.6.0`;
- align `/api/v1/meta` and its regression test with `0.6.0`;
- align the README current-development version with `0.6.0`;
- add dedicated v0.6.0 release verification documentation;
- retain the version-aware release-readiness and release-consistency gates;
- keep the release free of database migrations and production-data dependencies;
- require the complete applicable CI/security gate before publication.

The release establishes the next pre-1.0 release-integrity baseline after the v0.5.x verification hardening track. No intentional breaking `/api/v1` contract change is introduced.

### v0.5.9 — Release Metadata & Documentation Consistency

- advance the repository release metadata from `0.5.8` to `0.5.9`;
- align `/api/v1/meta` and its regression test with `0.5.9`;
- correct the README's stale current-development version reference;
- add dedicated v0.5.9 release verification documentation;
- retain the version-aware release-readiness and release-consistency gates introduced in v0.5.8;
- keep release publication blocked until the final v0.5.9 commit passes the complete applicable CI/security gate.

The release continues the project's release-integrity track by eliminating stale user-facing version metadata while preserving the existing validation-only release process.

### v0.5.8 — Release Gate Generalization & Verification Baseline

- generalize the release-readiness checker so it follows the repository `VERSION` instead of remaining pinned to an older patch release;
- align `VERSION`, `/api/v1/meta`, and the API metadata regression test with `0.5.8`;
- make the Makefile release-readiness help text version-neutral so future patch releases do not inherit stale version references;
- add dedicated v0.5.8 release verification documentation;
- retain validation-only release checks with no production credentials or production-data dependency;
- keep release publication blocked until the final v0.5.8 commit passes the complete applicable CI/security gate.

The release is focused on removing release-process drift exposed during v0.5.7 verification and establishing a reusable, version-aware verification baseline.

### v0.5.7 — Browser E2E Regression Repair & Release Integrity

- restore the complete authenticated browser E2E suite after a malformed test-file refactor caused the CI runner to execute an out-of-scope top-level assertion;
- preserve deterministic synthetic fixtures for catalog, inventory, purchasing, and reporting workflows;
- retain strict report currency assertions without introducing ambiguous locator matches;
- align `VERSION`, `/api/v1/meta`, and the API metadata regression test with `0.5.7`;
- add dedicated v0.5.7 release verification documentation;
- keep release verification independent of production credentials and production data;
- require the complete applicable CI/security gate to pass before publication.

The release is intentionally focused on correcting the concrete v0.5.6 follow-up regression and strengthening release integrity. Publication remains blocked until the final v0.5.7 commit has successful applicable CI and security checks.

### v0.5.6 — Release Verification & Reproducibility Hardening

- add a validation-only release-readiness gate for the v0.5.6 repository state;
- validate `VERSION`, changelog/release documentation, and the existing release consistency contract;
- add `make release-readiness` and `make release-verify` developer targets;
- add a dedicated GitHub Actions release-readiness workflow;
- verify Go module integrity with `go mod verify`;
- run `go vet ./...`, race-enabled backend tests, and a reproducible backend build in the release-readiness workflow;
- keep the release free of database migrations and production-data dependencies.

The release is intentionally focused on preventing inconsistent or unreproducible release state. Release publication remains blocked until the final release commit passes the applicable CI and security checks.

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
