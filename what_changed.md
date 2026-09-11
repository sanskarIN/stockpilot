# StockPilot — Work Continuity Log

## Current milestone

Phase 53 — v0.5.8 Release Gate Generalization & Verification Baseline.

## Repository state

- Default branch: `main`.
- v0.1.x preview releases and the v0.2.x reporting foundation series are retained in repository history.
- v0.3.x reporting foundations and v0.4.x reporting/replenishment hardening remain part of the established project history.
- v0.5.0–v0.5.4 added authenticated workflow and reporting regression hardening.
- v0.5.5 established release metadata consistency validation.
- v0.5.6 established release-readiness and reproducibility validation.
- v0.5.7 repaired the concrete browser E2E regression found after v0.5.6 and hardened its Reports workspace selectors.
- v0.5.8 removes a release-process regression where the release-readiness script remained pinned to v0.5.6.

## v0.5.8 completed scope

### Release-readiness generalization

- [x] Replace the hard-coded v0.5.6 expected version with the current repository `VERSION`.
- [x] Derive the release documentation path from the active semantic version.
- [x] Keep strict checks for `VERSION`, changelog entry, release document, release-consistency validation, and module verification.
- [x] Keep the checker validation-only with no tag creation or production-data mutation.

### Version alignment

- [x] Align `VERSION` with `0.5.8`.
- [x] Align `/api/v1/meta` with `0.5.8`.
- [x] Align `internal/httpapi/meta_version_test.go` with `0.5.8`.
- [x] Add the v0.5.8 changelog entry.
- [x] Add `docs/releases/v0.5.8.md`.
- [x] Make the Makefile release-readiness help text version-neutral.

## Verification gate

The release must not be called fully verified until the final main commit passes the actual applicable checks:

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
- No production credentials or production data are required for verification.
- No database migration is introduced by v0.5.8.
- Existing API behavior remains backward compatible unless explicitly documented otherwise.
- Focused, reviewable commits are preferred over artificial commit-count inflation.
- Do not create or move the `v0.5.8` tag until the final applicable CI/security gate is green.
- Do not claim zero bugs merely because a release metadata gate passes; inspect all applicable CI/security results.

## Next development track

1. Let GitHub Actions evaluate the v0.5.8 changes.
2. Inspect every failed job rather than assuming success.
3. Fix any newly exposed regression before release publication.
4. Re-run the complete applicable verification suite on the final release commit.
5. Review the final v0.5.8 diff against the v0.5.7 release baseline.
6. Publish `v0.5.8` only after the complete applicable gate is green.
