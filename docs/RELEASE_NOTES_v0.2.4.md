# StockPilot v0.2.4 — Stock Movement Velocity

## Overview

StockPilot v0.2.4 extends Reports & Analytics with a movement-history and velocity view built directly from persisted stock movement records.

## Highlights

### Stock movement history

- Added aggregated movement activity by product, location, and lot.
- Added configurable reporting windows from 1 to 365 days.
- The default reporting window is 30 days.
- Added movement count, inbound units, outbound units, and net units.
- Added the latest movement timestamp for each grouped result.

### Velocity reporting

- Added average daily outbound units for the selected reporting window.
- Results are deterministically ordered with the highest outbound activity first.
- Reporting is based on authoritative persisted `stock_movements` records rather than client-side estimates.

### API and export

New endpoint:

`GET /api/v1/reports/stock-movement-history`

Query parameters:

- `days`: reporting window, default 30, maximum 365.
- `limit`: result limit, default 1,000, maximum 5,000.
- `format=csv`: formula-safe CSV export.

CSV downloads retain StockPilot's privacy-oriented `no-store` / `no-cache` policy.

### Web application

- Added movement velocity to the Reports & Analytics workspace.
- Added recent movement activity table.
- Added outbound and average-daily-outbound visibility.
- Added CSV export action.

## Security and reliability

- The report is read-only.
- Existing authenticated reporting authorization remains in effect.
- Result generation is bounded server-side.
- Spreadsheet formula injection protection remains enabled for CSV output.
- No credentials, session tokens, passwords, or payment information are included in the report.

## Upgrade notes

v0.2.4 does not require a database migration. The feature reads the existing `stock_movements` table introduced by the core schema.

## Verification

Before stable publication, verify the final `main` commit through the configured repository gates:

- Go formatting, vet, unit tests, race tests, and server build.
- PostgreSQL migration and integration checks.
- Web quality and production build.
- CodeQL/security verification.
- Android and browser-companion checks where configured.
- E2E and accessibility checks where configured.
- Restore/rollback verification.
- Reproducible release-artifact checks.
- Blocker/critical-defect review.

## Release metadata

- **Version:** `v0.2.4`
- **Title:** `StockPilot v0.2.4 — Stock Movement Velocity`
- **Tag:** `v0.2.4`
- **Release type:** Stable, after all applicable gates pass
- **Prerelease:** No
- **Latest:** Yes, if replacing v0.2.3 as latest stable
- **Release date:** September 4, 2026

Thanks to everyone contributing through code, testing, documentation, issues, and feedback.
