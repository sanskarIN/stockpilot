# StockPilot

StockPilot is an open-source inventory and purchasing control system built around a Go API, PostgreSQL, and a React + TypeScript web application. The repository also contains a native Android client and a Manifest V3 browser companion.

> Made by the Sanskar

## Current capabilities

- Product catalog with SKU, category, supplier, barcode, unit cost, currency, reorder point, and reorder quantity.
- Warehouses and storage locations.
- Stock-in, stock-out, adjustment, receiving, and transactional transfer workflows.
- Lot and expiry-aware inventory rules.
- Aggregate reorder recommendations that include products with zero on-hand stock.
- Inventory valuation grouped safely by currency.
- Exact barcode lookup for scanner-driven clients.
- Purchase orders and line receiving.
- Password authentication, opaque sessions, CSRF protection, and role-based access control.
- Responsive React dashboard with PWA support and offline-safe service worker behavior.
- Native Android application with encrypted session persistence and release TLS enforcement.
- Manifest V3 browser companion with scoped optional host permissions.
- Docker deployment, PostgreSQL migrations, CI, CodeQL, and Dependabot.

## Architecture

```text
cmd/server             Go HTTP server entrypoint
cmd/admin              Operator-only administration bootstrap
internal/domain        Inventory, purchasing, catalog, and access rules
internal/repository    Persistence contracts
internal/postgres      PostgreSQL repositories and reporting queries
internal/httpapi       HTTP routes, validation, security middleware
migrations             Ordered PostgreSQL migrations
web                    React + TypeScript PWA
android                Native Android client
extension              Manifest V3 browser companion
.github/workflows      Backend, web, Android, extension, and security CI
```

The backend is the system of record. Inventory mutations are persisted transactionally in PostgreSQL and stock cannot become negative through supported movement workflows.

## Quick start with Docker

1. Copy the environment template and replace every development secret before exposing StockPilot outside your machine.

   ```sh
   cp .env.example .env
   ```

2. Start the stack.

   ```sh
   docker compose up --build
   ```

3. Apply migrations when required by your deployment workflow.

   ```sh
   make migrate
   ```

4. Create the initial administrator using the operator-only command documented by `cmd/admin` and your configured environment.

Never commit `.env`, production passwords, session peppers, database credentials, or TLS private keys.

## Development

### Backend

```sh
make fmt
make test
make vet
make build
```

### Web

```sh
make web-install
make web-build
```

### Android

A compatible Gradle installation is required by the Makefile unless you invoke the project wrapper available in your environment.

```sh
make android-lint
make android-test
make android-build
```

See [`android/README.md`](android/README.md) for local server, TLS, build, and release-security details.

### Browser extension

```sh
make extension-check
make extension-test
```

See [`extension/README.md`](extension/README.md) for installation and the companion security model.

## Inventory insights API

Authenticated users with the corresponding read permissions can use:

- `GET /api/v1/inventory/reorder-suggestions?limit=100` — product-level on-hand totals and suggested replenishment quantity.
- `GET /api/v1/reports/inventory-valuation?limit=100` — item valuation plus complete totals grouped by currency.
- `GET /api/v1/products/by-barcode/{barcode}` — exact lookup against the unique product barcode index.

The legacy `GET /api/v1/inventory/low-stock` endpoint remains available for balance-level diagnostics.

## Security baseline

- Secrets are supplied through environment configuration rather than source code.
- Session tokens are stored as peppered hashes server-side.
- Browser sessions use `HttpOnly` cookies and explicit CSRF confirmation for mutations.
- API access is enforced by role and permission.
- CORS origins are allowlisted.
- Common HTTP hardening headers are applied at the server boundary.
- Android release networking requires TLS; cleartext is limited to debug development configuration.
- The browser companion requests only the configured origin as an optional host permission.
- CI runs quality gates and CodeQL security analysis.

Review deployment configuration for your threat model before production use.

## Project status

StockPilot is under active development. The core end-to-end inventory foundation, authentication, web/PWA, Android client, and browser companion are implemented. The next milestones are tracked in [`ROADMAP.md`](ROADMAP.md); notable changes are recorded in [`CHANGELOG.md`](CHANGELOG.md) and [`what_changed.md`](what_changed.md).

## Contributing

Keep changes focused, tested, and reversible. For backend changes, run formatting, tests, vet, and a build. For client changes, run the relevant TypeScript/Android/extension quality gates. Never weaken inventory invariants or authentication checks to make a test pass.

## Support and contact

- GitHub: https://github.com/sanskarIN
- Buy Me a Coffee: https://buymeacoffee.com/sanskarIN
- Contact: sanskarin@outlook.in

## License

StockPilot is licensed under the MIT License. See [`LICENSE`](LICENSE).
