# StockPilot — Work Continuity Log

## Current milestone

Phase 50 — v0.5.5 Release Metadata Consistency Hardening.

## Repository state

- Default branch: `main`.
- v0.1.x preview releases and the v0.2.x reporting foundation series are retained in repository history.
- v0.3.x reporting foundations and v0.4.x reporting/replenishment hardening remain part of the established project history.
- v0.5.0–v0.5.4 added authenticated workflow and reporting regression hardening.
- v0.5.5 is based directly on the v0.5.4 baseline.
- Release metadata is synchronized across `VERSION`, `/api/v1/meta`, and its regression test.
- `scripts/check-release-consistency.sh` provides a read-only local consistency check.
- GitHub Actions now runs the same release consistency check as a dedicated CI job.
- Stable GitHub publication still requires the applicable CI and verification gates to pass.

## v0.5.5 completed scope

### Release metadata consistency

- [x] Add a validation-only release consistency checker.
- [x] Verify `VERSION` against the expected semantic version.
- [x] Verify `/api/v1/meta` exposes the same release version.
- [x] Verify the API metadata regression test asserts the same release version.
- [x] Expose the check through `make release-check`.
- [x] Run the check in GitHub Actions.
- [x] Align release metadata with `0.5.5`.
- [x] Document the release consistency contract.
- [x] Add dedicated v0.5.5 release notes.
- [x] Preserve existing changelog history while adding the v0.5.5 entry.

## v0.5.5 verification gate

Run all applicable repository checks before treating the release as stable:

```text
make release-check
make fmt
make vet
make test
make test-unit
make build
make web-build
make android-lint
make android-test
make android-build
make extension-check
make extension-test
```

Also verify CI/CodeQL, PostgreSQL migration/readiness behavior, browser E2E, accessibility, export safety, artifact integrity, and backup/restore procedures where applicable.

## Engineering rules

- The release consistency checker must remain read-only.
- No database migration is introduced by v0.5.5.
- No production credentials or production data are required for verification.
- Existing API behavior remains backward compatible.
- No authentication, authorization, inventory, purchasing, reporting, Android, extension, or web workflow contract is intentionally changed by this release.
- Focused, reviewable commits are preferred over artificial commit-count inflation.
- A release must not be described as fully verified until the actual applicable CI results are successful.

## Next development track

1. Verify the complete v0.5.5 CI suite on the final release commit.
2. Review the final v0.5.5 diff against v0.5.4 for unintended changes.
3. Publish the `v0.5.5` tag/release only after verification succeeds.
4. Begin v0.5.6 from the verified v0.5.5 baseline with a separately scoped improvement.
