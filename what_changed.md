# StockPilot — Work Continuity Log

## Current milestone

Phase 55 — v0.6.0 Release Integrity & Verification Baseline.

## Repository state

- Default branch: `main`.
- v0.1.x preview releases and the v0.2.x reporting foundation series are retained in repository history.
- v0.3.x reporting foundations and v0.4.x reporting/replenishment hardening remain part of the established project history.
- v0.5.0–v0.5.4 added authenticated workflow and reporting regression hardening.
- v0.5.5 established release metadata consistency validation.
- v0.5.6 established release-readiness and reproducibility validation.
- v0.5.7 repaired the concrete browser E2E regression found after v0.5.6 and hardened its Reports workspace selectors.
- v0.5.8 removed release-process drift by making release-readiness follow the active `VERSION`.
- v0.5.9 removed stale current-version documentation and aligned the repository/API/test release metadata with `0.5.9`.
- v0.6.0 continues the release-integrity track by aligning the active release metadata with `0.6.0` and preserving the verification contract.

## v0.6.0 completed scope

### Version alignment

- [x] Align `VERSION` with `0.6.0`.
- [x] Align `/api/v1/meta` with `0.6.0`.
- [x] Align `internal/httpapi/meta_version_test.go` with `0.6.0`.
- [x] Align the README current-development version with `0.6.0`.
- [x] Add the v0.6.0 changelog entry.
- [x] Add `docs/releases/v0.6.0.md`.

### Release-process continuity

- [x] Preserve the version-aware `scripts/check-release-readiness.sh` introduced in v0.5.8.
- [x] Preserve `scripts/check-release-consistency.sh` as the release metadata contract.
- [x] Preserve validation-only release checks with no production-data mutation.
- [x] Keep release documentation explicit about the required verification gate.

## Verification gate

The release must not be called fully verified until the final main/release commit passes the actual applicable checks:

```text
make release-check
make release-readiness
make vet
make test
make build
make web-build
```

Also require successful applicable Browser E2E, PostgreSQL migration/readiness smoke testing, CodeQL/security checks, Android checks, and extension checks where configured.

## Engineering rules

- Release checks remain read-only with respect to application and production data.
- No production credentials or production data are required for release metadata validation.
- No database migration is introduced by v0.6.0.
- Existing API behavior remains backward compatible unless explicitly documented otherwise.
- Focused, reviewable commits are preferred over artificial commit-count inflation.
- Do not create or move the `v0.6.0` tag until the final applicable CI/security gate is green.
- Do not claim zero bugs merely because release metadata checks pass; inspect all applicable CI/security results.

## Verification status

The v0.6.0 release-preparation changes have been applied to the dedicated `release/v0.6.0` branch. The GitHub repository tools available in this conversation can inspect and modify repository contents, but they do not execute the project's local Make, Go, Node, Android, PostgreSQL, or CodeQL commands directly. Therefore the release is **prepared but not yet claimed fully verified** until GitHub Actions or an equivalent execution environment reports the applicable checks as successful.

## Next development track

1. Let GitHub Actions evaluate the v0.6.0 release branch/PR.
2. Inspect every failed job rather than assuming success.
3. Fix any newly exposed regression before release publication.
4. Re-run the complete applicable verification suite on the final release commit.
5. Review the final v0.6.0 diff for unintended changes.
6. Merge the verified release changes into `main`.
7. Publish `v0.6.0` only after the complete applicable gate is green.

Made by the Sanskar
