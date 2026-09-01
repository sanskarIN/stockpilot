# StockPilot — Work Continuity Log

## Current milestone

Phase 2 — reporting/replenishment integration and release-readiness hardening.

## Repository state

- Default branch: `main`.
- `ci: verify Android and browser companion quality gates` was merged as PR #9.
- PR #10 remains open because its original branch predates PR #9 and now conflicts with the updated `main`; its functionality is being reconciled on `feat/replenishment-reporting-v2` without rewriting the published branch.
- The current work intentionally uses small, reviewable commits rather than squashing feature history.

## Completed in the current continuation

- [x] Added a security policy covering reporting, secrets, TLS, authorization, and dependency review.
- [x] Added contributor workflow and focused-commit guidance.
- [x] Added pull-request and issue templates.
- [x] Added API compatibility policy.
- [x] Added layered architecture documentation.
- [x] Added reproducible release checklist.
- [x] Added PostgreSQL restore-drill procedure.
- [x] Added replenishment and valuation domain read models.
- [x] Added repository contracts for barcode lookup, replenishment suggestions, and inventory valuation.
- [x] Added exact barcode persistence lookup and HTTP handler.
- [x] Added replenishment and valuation PostgreSQL queries with checked numeric conversion.
- [x] Added HTTP routes for barcode lookup, reorder suggestions, and inventory valuation.
- [x] Added unit and PostgreSQL integration tests for reporting calculations and data flow.
- [x] Added web API support for barcode lookup and reporting.
- [x] Added typed web models for replenishment and valuation data.
- [x] Added dashboard presentation of reorder recommendations and valuation totals with role-aware reporting behavior.
- [x] Added current `CHANGELOG.md` and `ROADMAP.md`.

## Verification status

The GitHub connector does not expose a local shell, so this continuation did not claim to execute the Go, web, Android, or extension test commands locally. The branch contains the requested automated tests and is intended to be validated by GitHub Actions when the pull request is opened.

## Known limitations

- PR #10 still needs to be reconciled/closed after its functionality is superseded by the new integration branch.
- Web product-management UI remains in PR #11 and depends on the replenishment branch lineage.
- Operational UI workflows (inventory mutations, purchasing, warehouse/location administration, lot/expiry flows) remain on the roadmap.
- Stable-release end-to-end browser and Android instrumentation coverage is still pending.

## Next exact tasks

1. Review and merge the replenishment/reporting integration PR.
2. Reconcile the catalog-management UI on top of the updated `main`.
3. Add guided inventory and purchase-order workflows to the web client.
4. Add append-only audit writes for sensitive mutations and expose the audit viewer.
5. Add CSV import/export and backup-retention deployment examples.
6. Add aging, expiry-risk, movement-velocity, and supplier analytics.
7. Add browser E2E and Android instrumentation suites.
8. Run the full release checklist and restore drill before first stable tagging.

## Commit discipline

Each functional boundary is intentionally represented by an individual commit where practical: documentation, repository contracts, persistence, HTTP handlers/routes, tests, and client integration are kept separately reviewable. GitHub's available connector does not provide an author-email override, so commits use the authenticated GitHub connection identity rather than forcing a specific email address.
