# StockPilot v0.2.2 — Expiry-Risk Reporting

## Release status

- Version: `v0.2.2`
- Release type: Stable
- Tag: `v0.2.2`
- Release date: 2026-09-04
- Prerelease: No
- Latest release: Yes

## Highlights

StockPilot v0.2.2 extends Reports & Analytics with expiry-risk visibility for inventory that has persisted lot expiry information.

### Expiry-risk reporting

- Reports expired inventory separately from inventory approaching expiry.
- Supports configurable risk windows so operators can choose an appropriate planning horizon.
- Uses deterministic server-side classification and ordering.
- Restricts the report to positive inventory balances with applicable expiry information.

### Export

- Provides bounded CSV export for expiry-risk results.
- Applies formula-safe CSV serialization.
- Uses `Cache-Control: no-store` and `Pragma: no-cache` for privacy-oriented downloads.
- Keeps exports behind existing authenticated reporting authorization.

### Web application

- Adds expiry-risk visibility to Reports & Analytics.
- Keeps reporting workflows read-only.
- Reuses the existing typed reporting API patterns.

### Quality and security

- Adds regression coverage for expiry classification boundaries.
- Adds authorization and export-limit coverage.
- Preserves existing authentication, authorization, CSRF, audit, and export privacy controls.
- Does not introduce credentials, passwords, raw session tokens, cookies, or payment information into reports.

## Upgrade notes

This is a pre-1.0 release. Existing deployments should run the normal migration procedure and complete the repository's restore/rollback checks before production rollout.

## Verification gates

Before declaring the GitHub release stable, confirm:

- Go formatting, vet, unit, race, and integration checks pass.
- PostgreSQL migration and integration checks pass.
- Web quality/build checks pass.
- Android and browser-companion checks pass where configured.
- CodeQL/security checks pass where configured.
- Browser E2E and accessibility checks pass where configured.
- Restore and migration rollback rehearsal is complete.
- Release artifacts are reproducible and reviewable.
- No blocker or critical defect remains open.

## Known semantics

Expiry risk is based on persisted lot expiry information and the configured reporting window. Inventory without applicable expiry information is not silently assigned an expiry date.

## Credits

Maintainer: Sanskar
Commit identity: `sanskarin@outlook.in`
