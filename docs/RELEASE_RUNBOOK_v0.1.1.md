# StockPilot v0.1.1 Release Runbook

## Release identity

- Version: `v0.1.1`
- Release class: maintenance release
- Target: `main`
- Recommended Git tag: `v0.1.1`
- Pre-release: no, provided all release gates are complete
- Latest release: yes

## 1. Source freeze

1. Confirm `main` is the intended release branch.
2. Confirm no open PR contains functionality required by v0.1.1.
3. Confirm the working tree represented by the release commit has the expected changelog and release notes.
4. Record the exact release commit SHA before creating the tag.

## 2. Automated verification

The release candidate must pass the repository CI and CodeQL workflows. At minimum verify:

- Go module tidiness.
- Go formatting.
- Go vet.
- Unit and integration tests.
- Server build.
- PostgreSQL migration smoke test.
- Web typecheck/build.
- Go CodeQL analysis.
- JavaScript/TypeScript CodeQL analysis.

Do not publish while a required check is failing or still running.

## 3. Manual verification

Complete `docs/RELEASE_CHECKLIST.md` in a production-like environment. In particular:

- Perform the PostgreSQL backup/restore drill in `docs/RESTORE_DRILL.md`.
- Verify authentication, authorization, CSRF protection, and session expiry.
- Check responsive and keyboard navigation behavior.
- Run Android lint/tests and a supported-device smoke test.
- Validate browser-companion permissions and installation/launcher flow.
- Verify HTTPS enforcement, encrypted Android session storage, security headers, CORS allow-list, backup retention, logging, and request IDs.
- Rehearse rollback before publication.

## 4. Tagging

Create an annotated Git tag named `v0.1.1` pointing at the verified release commit. Do not tag an unverified moving branch head.

## 5. GitHub Release

Create a GitHub Release from tag `v0.1.1` with title `StockPilot v0.1.1 — Maintenance Release` and use `docs/RELEASE_NOTES_v0.1.1.md` as the release description.

Mark it as the latest release only after every required gate is complete.

## 6. Artifacts

If binary artifacts are published, generate them from the exact tagged commit. Publish a SHA-256 checksum manifest beside the artifacts. Do not attach development builds or artifacts produced from a different commit.

## 7. Post-release verification

Immediately after publication:

1. Verify the release page and tag point to the expected commit.
2. Verify source archives and every attached artifact can be downloaded.
3. Start the published build using a clean release configuration.
4. Check `/healthz` and `/readyz`.
5. Smoke-test login and representative catalog, inventory, purchasing, and reporting reads.
6. Verify audit logging and request IDs for a representative operation.
7. Record the outcome and any rollback decision.

## Rollback

If a release-blocking defect is discovered, stop rollout, preserve logs and the failing artifact, and follow the documented rollback procedure. Do not rewrite or move the published `v0.1.1` tag; publish a corrective patch release instead.
