# StockPilot — Work Continuity Log

## Current milestone

Unreleased operational-insights milestone — move beyond the core MVP into reliable replenishment, valuation, scanner lookup, secure web entrypoint, and release-quality project documentation.

## Repository state reviewed

- Default branch: `main`.
- Continuation branch: `feat/replenishment-reporting`.
- Pull request: `#10` — `feat: add replenishment and inventory reporting`.
- Base commit for this continuation: `09e90918e87f9c26a71acc4d103aa6a1f2c6d655`.
- No open GitHub issues were present when this continuation started.
- The original checklist in this file was stale: the repository already contained the Go/PostgreSQL core, web/PWA, authentication/RBAC, Android app, Manifest V3 extension, tests, CI, and security work.
- The repository code, CI results, and commit history are now treated as the source of truth for future continuations.

## Implemented foundation already present before this continuation

- [x] Catalog domain: categories, suppliers, products, SKU/barcode metadata, costing, reorder settings.
- [x] Warehouse/location model.
- [x] Stock balances, stock movements, transfers, and negative-stock prevention.
- [x] Lot and expiry rules.
- [x] Purchase orders and receiving.
- [x] PostgreSQL schema, constraints, indexes, and migration runner.
- [x] Go HTTP server, health/readiness endpoints, request IDs, validation, and security middleware.
- [x] Authentication, opaque sessions, password hashing, session-token peppering, CSRF protection, and RBAC.
- [x] Operator-only administrator bootstrap command.
- [x] React + TypeScript responsive dashboard.
- [x] PWA manifest/service worker integration.
- [x] Native Android client with encrypted session persistence and release TLS enforcement.
- [x] Manifest V3 browser companion with scoped optional host permissions.
- [x] Backend/web/Android/extension CI, CodeQL, and Dependabot baselines.

## Work completed in this continuation

### Replenishment

- [x] Added `domain.ReorderSuggestion`.
- [x] Added aggregate product-level reorder queries instead of relying only on individual location/lot balances.
- [x] Zero-stock active products are now eligible for replenishment recommendations even when no inventory-balance row exists.
- [x] Suggested order quantity targets `reorder point + reorder quantity` and subtracts current aggregate on-hand stock.
- [x] Added checked overflow handling for reorder-target arithmetic.
- [x] Added `GET /api/v1/inventory/reorder-suggestions`.
- [x] Added unit tests for replenishment calculations.

### Inventory valuation

- [x] Added valuation item, currency-total, and report domain models.
- [x] Added product-level inventory valuation from aggregate on-hand stock and unit cost.
- [x] Added complete totals grouped by currency so incompatible currencies are never silently combined.
- [x] Valuation multiplication uses PostgreSQL numeric arithmetic before checked conversion to application `int64` minor units.
- [x] Added `GET /api/v1/reports/inventory-valuation`.
- [x] Added API tests for report serialization.

### Barcode workflow foundation

- [x] Confirmed the existing schema already enforces unique non-empty product barcodes.
- [x] Added repository support for exact barcode lookup.
- [x] Added `GET /api/v1/products/by-barcode/{barcode}` for scanner-driven clients.
- [x] Added input validation and API coverage for exact barcode lookup.

### Web dashboard

- [x] Replaced the misleading per-balance low-stock dashboard signal with aggregate reorder recommendations.
- [x] Added suggested-order quantities and target stock to the replenishment table.
- [x] Added inventory valuation display grouped by currency.
- [x] Preserved the existing RBAC boundary: operators do not receive `reports:read`; the dashboard skips valuation for that role instead of broadening privileges or failing the entire load.
- [x] Added a typed web client method for exact barcode lookup.

### CI and web entrypoint stabilization

- [x] Added a verified `go.sum` so Go module downloads, vetting, tests, builds, and CodeQL can resolve declared dependencies reproducibly.
- [x] Added Vite client type declarations through `web/src/vite-env.d.ts`, fixing `ImportMeta.env` TypeScript resolution in CI.
- [x] Upgraded CodeQL workflow actions from v3 to v4.
- [x] Rewired `web/src/main.tsx` to use the existing secure login/session flow and the authenticated `Dashboard` component instead of the legacy unauthenticated dashboard implementation.
- [x] Added initial session verification, sign-in, sign-out, and session-expiry handling at the web application entrypoint.
- [x] Added the pgx transitive build dependencies to `go.mod` after `go vet` correctly reported that the module graph required them.
- [x] Added `golang.org/x/text v0.24.0` to the indirect Go dependency graph after CI correctly identified the `pgx/pgconn` SCRAM dependency.
- [x] Refreshed CI actions to current Node 24-compatible major versions: `actions/checkout@v6`, `actions/setup-go@v7`, and `actions/setup-node@v6`.
- [x] Refreshed the CodeQL workflow to the same current checkout/setup action majors while keeping CodeQL Action v4.

### Reliability and operations

- [x] Added a real PostgreSQL integration test for barcode lookup, reorder suggestions, and valuation against the migrated schema.
- [x] Extended PostgreSQL CI to run the reporting integration test.
- [x] Repaired the pre-existing broken `make backup` command by adding `scripts/backup.sh` and invoking it portably with `sh`.
- [x] Backup creation now uses restrictive file permissions, custom PostgreSQL dump format, partial-file cleanup, and empty-output protection.
- [x] CI validates backup-script shell syntax.
- [x] Improved the Go formatting gate so a failure prints the exact files requiring `gofmt`.
- [x] The PR formatting gate exposed existing drift in `internal/domain/catalog_test.go`, `internal/domain/purchasing_test.go`, `internal/httpapi/access.go`, and `internal/httpapi/api_test.go`; all four files were normalized with `gofmt` without behavior changes.

### Project documentation

- [x] Added root `README.md` with architecture, setup, quality commands, security baseline, API insight endpoints, clients, support/contact, credit, and license information.
- [x] Added `CHANGELOG.md`.
- [x] Added `ROADMAP.md` based on actual repository state.
- [x] Replaced the stale Phase 1 continuity checklist in this file.

## Verification

### Static/unit/build gates configured

- Go formatting check: `gofmt -l cmd internal`, with exact non-formatted file reporting and a failing exit status when drift exists.
- Go vet: `go vet ./...`.
- Go race/unit tests: `go test -race -coverprofile=coverage.out ./...`.
- Go server build: `go build -trimpath -o bin/stockpilot ./cmd/server`.
- Web typecheck/build: `npm run build` in `web`.
- PostgreSQL migration readiness smoke test.
- PostgreSQL reporting integration test: `go test ./internal/postgres -run TestReportingIntegration -count=1`.
- Backup script syntax: `sh -n scripts/backup.sh`.
- Existing Android and extension workflows remain separate quality gates.
- CodeQL remains an independent pull-request security gate, using CodeQL Action v4 with current checkout/setup runtimes.

### PR validation status

- PR `#10` remains the validation surface for the replenishment/reporting milestone.
- The first validation pass correctly failed the Go formatting gate, which revealed previously hidden formatting drift.
- The formatting gate was made diagnostic and all reported files were then normalized.
- A later pass exposed missing Vite `ImportMeta.env` types and missing Go checksums; both were addressed.
- `go vet` then identified the pgx transitive build dependencies; those were declared in `go.mod`.
- The following pass identified the `golang.org/x/text/secure/precis` dependency used by pgx SCRAM; `golang.org/x/text v0.24.0` was then declared as an indirect dependency.
- The PostgreSQL smoke test failed for the same missing `x/text` dependency before the module-graph fix; the database container itself initialized and became healthy, so the failure was dependency resolution rather than PostgreSQL startup.
- CI and CodeQL were refreshed to current Node 24-compatible action majors, eliminating the earlier forced-Node-24 deprecation warnings from checkout/setup actions.
- The latest PR head has fresh CI and CodeQL runs queued; final merge must occur only after the latest required checks complete successfully.

### Environment note

The connected GitHub environment does not provide a local checkout with dependency/network access for an independent full build. The executable GitHub Actions gates therefore remain the authoritative validation environment for repository-wide CI, while the code and commit history provide the implementation source of truth.

## Known limitations / remaining product work

- Product/category/supplier/warehouse/location management is still primarily API-backed; full web catalog CRUD is the next client milestone.
- Stock movement, transfer, purchase-order creation, and receiving need guided first-class web workflows.
- Exact barcode lookup is now available, but camera-based barcode/QR scanning UI is not yet implemented.
- Audit permissions exist, but a first-class append-only audit event store/viewer still remains.
- CSV import/export remains.
- Backup creation now works; restore drill tooling, retention, and scheduled deployment examples remain.
- Full browser E2E and Android instrumentation coverage remain release-hardening tasks.
- Accessibility and production restore/migration-rollback audits remain before a stable release.

## Next exact tasks

1. Complete PR `#10` CI validation and merge the replenishment/reporting milestone.
2. Build web product/category/supplier management screens with validation, loading/error states, and role-aware write controls.
3. Add guided stock movement and transfer workflows.
4. Add purchase-order creation/editing/receiving UI and an action to seed draft orders from reorder recommendations without bypassing approvals.
5. Add camera barcode/QR scanning on supported clients backed by the exact barcode endpoint.
6. Add append-only audit events and an audit viewer.
7. Add CSV import/export with dry-run validation and row-level errors.
8. Add restore tooling, retention hooks, and scheduled-backup deployment examples.
9. Add end-to-end browser tests and continue stable-release acceptance work.

## Migration notes

No schema migration is required for this continuation. The existing `products_barcode_uq` partial unique index already guarantees uniqueness for non-empty barcodes. Replenishment and valuation are computed from the current `products` and `inventory_balances` schema.

## Release notes draft

Unreleased: add aggregate reorder recommendations that include stockouts, currency-safe inventory valuation, exact barcode lookup, dashboard reporting integration, real PostgreSQL reporting integration coverage, a repaired database-backup command, stronger CI diagnostics, Go formatting cleanup, reproducible Go dependency checksums, CodeQL v4, a secure authenticated web entrypoint, and current Node 24-compatible GitHub Actions runtimes.

## Continuation commits

- `f6a8801` — `feat(replenishment): add reorder and valuation models`
- `314a399` — `feat(replenishment): extend inventory reporting contract`
- `69e4503` — `feat(reporting): implement reorder and valuation queries`
- `7c8ab14` — `test(reporting): cover replenishment calculations`
- `3064355` — `feat(api): add replenishment and valuation handlers`
- `5142836` — `feat(api): expose replenishment reporting routes`
- `b1e734a` — `feat(catalog): add barcode lookup contract`
- `f4cd13a` — `feat(catalog): implement exact barcode lookup`
- `59e4259` — `feat(api): add scanner-friendly barcode lookup`
- `9526133` — `feat(api): route exact barcode lookup`
- `d2308a1` — `test(api): cover barcode and inventory reports`
- `3235ca4` — `feat(web): type replenishment and valuation data`
- `fcd27fb` — `feat(web): consume replenishment reporting APIs`
- `9123151` — `feat(web): surface reorder and valuation insights`
- `986e758` — `docs: add StockPilot project guide`
- `1cfbb30` — `docs: add delivery roadmap`
- `3763d77` — `docs: add project changelog`
- `aec5d0f` — `fix(ops): repair database backup command`
- `21c6a27` — `fix(ops): invoke portable backup script`
- `6255d7d` — `test(postgres): verify reporting against real schema`
- `21b5537` — `ci: verify reporting queries and backup script`
- `6921c89` — `fix(web): preserve operator dashboard permissions`
- `2b8206b` — `fix(web): render valuation by role permission`
- `fa96b7f` — `docs: refresh StockPilot continuity log`
- `d5330a8` — `docs: mark backup creation implemented`
- `7218dcb` — `ci: report files that need gofmt`
- `8a089f9` — `style(go): format catalog tests`
- `1c8c770` — `style(go): format purchasing tests`
- `ebd7742` — `style(go): format access handler`
- `6515956` — `style(go): format API tests`
- `c198260` — `fix(deps): add verified Go module checksums`
- `5bfe376` — `fix(web): load Vite import-meta environment types`
- `7f1a057` — `ci(security): upgrade CodeQL actions to v4`
- `6c8ab0b` — `fix(web): wire authenticated dashboard entrypoint`
- `94659eb` — `docs: record CI and web entrypoint fixes`
- `e5c17d4` — `fix(deps): declare pgx transitive build dependencies`
- `cf1f409` — `fix(deps): declare pgx text dependency`
- `9879f1c` — `ci: refresh GitHub Actions runtimes`
- `ef24f05` — `ci(security): refresh action runtimes`
- `docs-current` — `docs: refresh operational validation status`
