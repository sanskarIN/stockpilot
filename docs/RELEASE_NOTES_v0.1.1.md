# StockPilot v0.1.1 — Maintenance Release

## Release summary

StockPilot v0.1.1 is the first maintenance release after the v0.1.0 public preview. It focuses on release-process hardening, documentation clarity, reproducible verification, and preserving the stability of the existing catalog, inventory, purchasing, reporting, audit, browser-companion, and Android workflows.

## Highlights

- Keeps the v0.1.x release line focused on stability rather than introducing a large new feature set.
- Documents the verification expectations for Go, PostgreSQL migrations, the web application, CodeQL, and release smoke checks.
- Preserves transactional CSV product import behavior and its server-side revalidation guarantees.
- Preserves append-only auditability for sensitive business and authentication lifecycle events.
- Keeps browser-companion barcode handoff navigation-only so stock mutations remain inside the authenticated web application.
- Keeps production Android networking under TLS enforcement.

## Verification gates

Before publishing v0.1.1, the release candidate must have green required CI and CodeQL checks. The release owner should additionally complete the manual checks in `docs/RELEASE_CHECKLIST.md`, including backup/restore verification, Android/device smoke testing, browser-companion installation and handoff testing, keyboard/accessibility review, and post-release health/readiness checks.

## Upgrade notes

v0.1.1 is intended as a compatible maintenance release within the v0.1.x line. Apply database migrations in their recorded order and follow `docs/RESTORE_DRILL.md` before performing a production restore.

## Known limitations

- The project remains pre-1.0.
- Advanced inventory analytics, broader export/reporting capabilities, expanded concurrency coverage, browser E2E, Android instrumentation, and additional release automation remain planned for later milestones.

## Release artifacts

Publish source archives through GitHub Releases. When distributable binaries are attached, generate SHA-256 checksums and publish the checksum manifest alongside the artifacts.
