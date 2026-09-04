# StockPilot — Work Continuity Log

## Current milestone

Phase 42 — v0.3.0 Reporting Foundations & Analytics Safety.

## Repository state

- Default branch: `main`.
- v0.1.x preview releases and the v0.2.0 reporting foundation are merged.
- v0.2.1 inventory aging, v0.2.2 expiry risk, v0.2.3 export hardening, v0.2.4 movement velocity, v0.2.5 supplier performance, and v0.2.6 warehouse/location valuation are implemented in the reporting series.
- v0.2.7 replenishment readiness, v0.2.8 scalability/purchasing-trends work, v0.2.9 reliability preparation, and v0.2.10 performance/operational-insights planning were carried forward in the repository history.
- The v0.2.8 release metadata previously contained a v0.2.9 body and should remain an explicit historical cleanup item.
- `VERSION` now records `0.3.0` as the current release marker.
- `docs/releases/v0.3.0.md` contains the release checklist and notes.
- `internal/reporting/` now provides bounded periods, previous-period calculation, period-over-period changes, and bounded query parameters with unit tests.
- Stable GitHub publication still requires the applicable CI and verification gates to pass.

## v0.3.0 completed scope

### Reporting foundations

- [x] Add bounded inclusive reporting periods with a 1–365 day safety envelope.
- [x] Add deterministic previous-period calculation using the same number of days.
- [x] Add period-over-period delta and percentage calculation.
- [x] Represent zero-baseline percentage changes as undefined rather than inventing a value.
- [x] Add centralized non-negative limit/offset validation.
- [x] Add safe parsing for non-negative integer query values.
- [x] Add focused unit tests for period, trend, and bounds behavior.
- [x] Document the reporting package's read-only contract.
- [x] Add a canonical `VERSION` marker and v0.3.0 release documentation.

## v0.3.0 verification gate

Run all applicable repository checks before treating the release as stable:

```text
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

Also verify CI/CodeQL, PostgreSQL migration/readiness behavior, accessibility, export safety, artifact integrity, and backup/restore procedures where applicable.

## Engineering rules

- Reporting remains read-only.
- No automatic purchase-order creation or inventory mutation is introduced by v0.3.0.
- No database migration is introduced by the reporting foundation work.
- Existing API contracts must remain backward compatible.
- Bounds and period validation must fail closed and deterministically.
- Focused, reviewable commits are preferred over artificial commit-count inflation.
- Repository commits use `sanskarin@outlook.in` where the connected GitHub identity supports author attribution.

## Next development track

1. Wire the shared reporting primitives into HTTP query parsing without changing existing response contracts.
2. Add server-side report metadata for period, bounds, generated-at, and completeness state.
3. Add bounded analytical trend endpoints backed by authoritative PostgreSQL aggregation.
4. Add regression tests for cancellation, timeout propagation, deterministic ordering, and duplicate prevention.
5. Integrate report-period context and retry/empty/error states into the web dashboard.
6. Re-run the complete cross-platform and security verification suite before the next release.
