# StockPilot — Work Continuity Log

## Current milestone

Phase 51 — v0.5.6 Release Verification and Browser E2E Stabilization.

## Repository state

- Default branch: `main`.
- v0.1.x preview releases and the v0.2.x reporting foundation series are retained in repository history.
- v0.3.x reporting foundations and v0.4.x reporting/replenishment hardening remain part of the established project history.
- v0.5.0–v0.5.4 added authenticated workflow and reporting regression hardening.
- v0.5.5 established release metadata consistency validation.
- v0.5.6 adds release-readiness validation and reproducibility-focused verification.
- `VERSION` is `0.5.6` and release metadata consistency checks pass on the current v0.5.6 branch state.

## v0.5.6 completed scope

### Release verification

- [x] Add validation-only v0.5.6 release-readiness checking.
- [x] Expose release readiness through `make release-readiness`.
- [x] Keep `make release-check` as the release metadata consistency gate.
- [x] Verify Go module integrity with `go mod verify`.
- [x] Run Go vet, backend tests, and backend build in CI.
- [x] Run web typecheck/build in CI.
- [x] Run PostgreSQL migration/readiness smoke testing in CI.
- [x] Run CodeQL analysis; Go CodeQL analysis has completed successfully on the current fix commit.
- [x] Document the complete release verification sequence.

### Browser E2E stabilization

- [x] Align purchase-order E2E with the rendered `Create draft order` action.
- [x] Wait for the reports refresh state before validating report content.
- [x] Remove the fragile duplicate-prone report currency locator.
- [x] Add the missing `/api/v1/reports/replenishment-readiness` synthetic fixture used by the reports workspace.
- [x] Keep all authenticated E2E data deterministic and offline from production services.

## Latest verified CI findings

The latest CI run for PR #67 tested the actual E2E fix commit. PostgreSQL smoke, Go quality, web quality, and release metadata consistency passed. Browser E2E had 6/7 tests passing and exposed one precise reports assertion failure: `₹1,500.00` matched two visible elements, causing Playwright strict-mode failure. The same run also exposed a missing replenishment-readiness mock through Vite proxy `ECONNREFUSED` warnings. These findings are now addressed by the latest test changes.

The release-readiness workflow independently exposed a documentation-contract failure because the v0.5.6 release document listed the underlying shell script but did not contain the literal `make release-check` command required by the readiness checker. The release document has now been corrected to use the public Make targets.

## Release verification gate

Run all applicable checks before treating v0.5.6 as stable:

```text
make release-check
make release-readiness
make vet
make test
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

- Release checks must remain read-only with respect to application and production data.
- No production credentials or production data are required for verification.
- Existing API behavior remains backward compatible unless explicitly documented otherwise.
- Focused, reviewable commits are preferred over artificial commit-count inflation.
- A release must not be described as fully verified until the actual applicable CI results for the final release commit are successful.
- Do not create or move the `v0.5.6` tag until the complete applicable CI/security/release gate is green.

## Next development track

1. Wait for CI to evaluate the latest Browser E2E and release-readiness fixes.
2. Inspect any remaining failed job rather than assuming success.
3. Merge PR #67 only after all required checks are successful.
4. Re-run the complete release verification against the resulting main commit.
5. Review the final diff against the v0.5.5 baseline.
6. Publish the v0.5.6 GitHub release/tag only after the final verification gate is green.
7. After v0.5.6 is stable, begin the next separately scoped product/engineering milestone.
