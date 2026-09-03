# StockPilot v0.1.2 — Release Notes

## Overview

StockPilot v0.1.2 is a small feature-and-hardening release following the v0.1.1 maintenance cycle. It establishes the reusable CSV export foundation and documents the security and operational contract required before application-level CSV download endpoints are enabled.

## Highlights

- Added a reusable `internal/csvexport` package built only on Go's standard library.
- Added RFC 4180-compatible quoting and buffered error-aware flushing.
- Added opt-in spreadsheet-formula protection for downloadable CSV data.
- Added a dedicated CSV export architecture and security design document.
- Kept export serialization independent from repositories so authorization, filtering, ordering, and resource limits remain application responsibilities.

## Security

- Formula-safe mode can prevent spreadsheet applications from treating values beginning with `=`, `+`, `-`, or `@` as formulas.
- Export design requires server-side authorization, bounded result sizes, explicit column contracts, deterministic ordering, and exclusion of secrets.
- Sensitive export actions should be audited when they expose administrator-only or otherwise sensitive datasets.

## Verification

The new CSV package includes tests covering:

- CSV quoting for commas and newlines.
- Formula-safe escaping.
- Empty-header validation.
- Nil-writer error handling.

Before publishing the release, run the complete repository quality gates from `docs/RELEASE_CHECKLIST.md`, including Go tests, race-enabled tests, web checks, PostgreSQL smoke testing, Android checks, extension checks, and CodeQL where configured.

## Upgrade notes

No database migration is required by the CSV foundation in this release. Existing JSON API contracts remain unchanged.

## Known limitations

Application-level CSV download endpoints are intentionally not part of this foundation commit. The next implementation stage should add bounded, authorized exports for catalog, inventory, low-stock/reorder, lot inventory, purchasing, and approved audit/report datasets.

XLSX and PDF generation remain out of scope until the CSV contracts and resource-management model are stable.

## Release metadata

- Version: `v0.1.2`
- Suggested title: `StockPilot v0.1.2 — CSV Export Foundation`
- Suggested tag: `v0.1.2`
- Release type: normal pre-1.0 release
- Pre-release flag: **recommended: yes** until the full manual release gates are completed.
