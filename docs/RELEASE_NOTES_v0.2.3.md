# StockPilot v0.2.3 — Export Hardening

## Overview

StockPilot v0.2.3 is a focused patch release that hardens CSV export behavior and corrects regression coverage discovered during release verification.

## Highlights

### CSV formula safety

- Corrected audit-export metadata sanitization for formula-like JSON string values.
- Preserved the shared formula-safe CSV serializer for authenticated exports.

### Receipt-history export correctness

- Corrected regression expectations for exported timestamps so they match StockPilot's canonical UTC representation.
- Preserved formula-safe handling of receipt-history notes and other exported cells.

### Reliability and security

- No new authentication, authorization, or session behavior is introduced.
- Existing export bounds remain enforced.
- Existing `no-store` / `no-cache` download policy remains in effect.
- Export workflows remain read-only.

## Upgrade notes

v0.2.3 is a patch release over v0.2.2. No new database migration is introduced by this patch. Existing deployments can follow the normal application upgrade procedure.

## Verification

Before treating the release as stable, verify the final `main` commit through GitHub Actions, including:

- Go formatting, vet, race tests, unit tests, and server build.
- PostgreSQL migration smoke/integration checks.
- Web quality and build checks.
- CodeQL/security checks.
- Any configured Android, browser-companion, E2E, accessibility, restore/rollback, and artifact gates.

## Release metadata

- **Version:** `v0.2.3`
- **Title:** `StockPilot v0.2.3 — Export Hardening`
- **Tag:** `v0.2.3`
- **Release type:** Stable
- **Prerelease:** No
- **Latest:** Yes, if this is intended to replace v0.2.2 as the latest stable release
- **Release date:** September 4, 2026

## Contributors

Thanks to everyone contributing through code, testing, documentation, issues, and feedback.
