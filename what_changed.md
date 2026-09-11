# StockPilot — Work Continuity Log

## Current milestone

Phase 52 — v0.5.7 Browser E2E Regression Repair & Release Integrity.

## Repository state

- Default branch: `main`.
- v0.1.x preview releases and the v0.2.x reporting foundation series are retained in repository history.
- v0.3.x reporting foundations and v0.4.x reporting/replenishment hardening remain part of the established project history.
- v0.5.0–v0.5.4 added authenticated workflow and reporting regression hardening.
- v0.5.5 established release metadata consistency validation.
- v0.5.6 established release-readiness and reproducibility validation.
- v0.5.7 is the next patch baseline focused on repairing the concrete browser E2E regression found after v0.5.6.
- `VERSION`, `/api/v1/meta`, and the API metadata regression test are aligned to `0.5.7`.

## v0.5.7 completed scope

### Browser E2E regression repair

- [x] Identify the CI collection failure caused by a top-level `expect(...)` statement in `web/tests/e2e/mutation-workflows.spec.ts`.
- [x] Restore the complete authenticated mutation-workflow E2E suite.
- [x] Restore the Playwright `expect` import and test definitions.
- [x] Preserve deterministic catalog, inventory, purchasing, and reporting fixtures.
- [x] Preserve the strict report currency assertion.
- [x] Preserve the replenishment-readiness synthetic report fixture.

### Release metadata

- [x] Align `VERSION` with `0.5.7`.
- [x] Align `/api/v1/meta` with `0.5.7`.
- [x] Align `internal/httpapi/meta_version_test.go` with `0.5.7`.
- [x] Add dedicated `docs/releases/v0.5.7.md` release verification documentation.
- [x] Add the v0.5.7 changelog entry.

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
- No database migration is introduced by v0.5.7.
- Existing API behavior remains backward compatible unless explicitly documented otherwise.
- Focused, reviewable commits are preferred over artificial commit-count inflation.
- Do not create or move the `v0.5.7` tag until the final applicable CI/security gate is green.

## Next development track

1. Confirm CI evaluates the repaired browser E2E suite and v0.5.7 metadata.
2. Inspect every failed job rather than assuming success.
3. Fix any newly exposed regression before release publication.
4. Re-run the complete applicable verification suite on the final release commit.
5. Review the final v0.5.7 diff against the v0.5.6 release baseline.
6. Publish `v0.5.7` only after the complete applicable gate is green.
