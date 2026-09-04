# StockPilot v0.2.1 — Inventory Aging

## Overview

StockPilot v0.2.1 adds the first end-to-end **inventory aging report**. It exposes deterministic age buckets over positive inventory balances, backed by PostgreSQL movement history and surfaced in the web Reports & Analytics workspace.

## Added

- Inventory aging domain contracts and deterministic bucket rules.
- PostgreSQL-backed inventory aging query.
- Bounded inventory aging HTTP report endpoint:
  - `GET /api/v1/reports/inventory-aging`
- Formula-safe CSV export through the report endpoint:
  - `GET /api/v1/reports/inventory-aging?format=csv`
- Web Reports & Analytics inventory-aging panel.
- Web export action for aging data.
- HTTP coverage for aging export limit normalization.

## Aging model

Aging is calculated as the elapsed whole-day age since the most recent persisted stock movement for the product/location/lot balance. If no movement row is available, the current balance timestamp is used as the fallback reference.

Buckets are deterministic:

| Age | Bucket |
| --- | --- |
| 0–30 days | `0-30` |
| 31–60 days | `31-60` |
| 61–90 days | `61-90` |
| 91–180 days | `91-180` |
| 181+ days | `181+` |

Only positive inventory balances are reported.

## API behavior

- Default limit: 1,000 rows.
- Maximum limit: 5,000 rows.
- Non-positive limits fall back to the default.
- Values above the maximum are capped.
- Results are ordered from oldest movement timestamp first, then product/location/lot for deterministic output.
- CSV responses use the shared privacy-oriented download headers and formula-safe serializer.

## Web experience

The Reports & Analytics workspace now loads inventory aging alongside overview, purchasing, and valuation reports.

Users with export access can export the aging dataset as CSV. The report remains read-only.

## Security and privacy

- The report uses the existing authenticated HTTP stack.
- CSV output uses formula-safe serialization.
- Export delivery uses the existing `no-store`/`no-cache` policy.
- No credentials, session tokens, cookies, or payment data are introduced into the report schema.

## Verification status

Focused domain and HTTP tests have been added to the repository. Full stable-release verification must still be executed through the authoritative CI/local environment, including Go, PostgreSQL integration, Web, Android, extension, CodeQL, end-to-end, accessibility, restore/rollback, and reproducible-artifact gates.

Do not treat unexecuted checks as passed.

## Release metadata

- Version: `v0.2.1`
- Suggested title: `StockPilot v0.2.1 — Inventory Aging`
- Release channel: **Stable**, after all release gates pass.
- Tag: `v0.2.1`
