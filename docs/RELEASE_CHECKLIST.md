# StockPilot Release Checklist

## Source and dependencies

- [ ] `main` contains the intended release commit set.
- [ ] No unresolved pull requests contain required release functionality.
- [ ] Go modules are tidy.
- [ ] Web lockfile is present and dependency audit is reviewed.
- [ ] Android and GitHub Actions dependency updates are reviewed for compatibility and licensing.

## Backend and database

- [ ] `make fmt` produces no changes.
- [ ] `make vet` passes.
- [ ] `make test` passes.
- [ ] `make test-unit` passes.
- [ ] PostgreSQL integration tests pass against the supported PostgreSQL version.
- [ ] Migrations apply cleanly to an empty database.
- [ ] A backup is created and a restore drill succeeds.

## Web

- [ ] `make web-build` passes.
- [ ] Authentication, authorization, CSRF, and session expiry are manually smoke-tested.
- [ ] Responsive layouts are checked at mobile and desktop widths.
- [ ] Keyboard navigation and visible focus are checked.
- [ ] Empty, loading, error, and unauthorized states are present.

## Android

- [ ] `make android-lint` passes.
- [ ] `make android-test` passes.
- [ ] Debug APK assembles reproducibly.
- [ ] Release configuration enforces HTTPS.
- [ ] Encrypted session storage is verified.
- [ ] Critical read paths are smoke-tested on a supported emulator/device.

## Browser companion

- [ ] Manifest validation passes.
- [ ] Extension tests pass.
- [ ] Host permission scope matches the configured origin.
- [ ] No password/session/API credential is persisted.
- [ ] Installation and launcher flow are smoke-tested.

## Security and operations

- [ ] Production secrets exist only in deployment secret stores.
- [ ] TLS certificates and trusted origins are configured.
- [ ] CORS allow-list is reviewed.
- [ ] Security headers are present.
- [ ] Database backup retention is configured.
- [ ] Logging and request IDs are available for incident diagnosis.
- [ ] Rollback procedure has been documented and rehearsed.

## Publication

- [ ] Changelog updated.
- [ ] Version identifiers updated consistently.
- [ ] Release notes include migrations and operational changes.
- [ ] Git tag and release artifacts are reproducible.
- [ ] Post-release smoke test is completed.
