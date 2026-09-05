# StockPilot

StockPilot is an open-source inventory and purchasing platform with a Go/PostgreSQL backend, a React/TypeScript web application, a native Android client, and a Manifest V3 browser companion.

The project is currently in pre-1.0 development. The current development release is `0.2.6`; the Android and browser companion clients remain independently versioned.

## What StockPilot includes

- Product catalog, categories, suppliers, warehouses, and locations.
- Inventory movements, transfers, lot tracking, balances, and low-stock reporting.
- Purchase orders and receiving workflows.
- Session-based authentication and role-based access control.
- Administrator bootstrap tooling and user-management APIs.
- PostgreSQL migrations with automatic migration support.
- Responsive React/TypeScript web dashboard and PWA assets.
- Native Kotlin Android client with encrypted session storage.
- Manifest V3 browser companion with optional per-server host permissions.
- Health/readiness endpoints, structured HTTP hardening, Docker deployment, backups, and CI/CodeQL checks.
- Reports & Analytics for inventory valuation, warehouse/location valuation, inventory aging, expiry risk, movement velocity, and supplier performance.

## Repository layout

```text
.
├── android/                 Native Kotlin Android application
├── cmd/
│   ├── admin/               Administrator bootstrap CLI
│   └── server/              StockPilot HTTP server
├── extension/               Manifest V3 browser companion
├── internal/
│   ├── auth/                Authentication/session service
│   ├── config/              Environment configuration
│   ├── domain/              Core domain models and rules
│   ├── httpapi/             HTTP API, auth middleware, and hardening
│   ├── idgen/               Identifier generation
│   ├── postgres/            PostgreSQL repositories and migrations runner
│   └── repository/          Repository contracts
├── migrations/              Versioned PostgreSQL schema migrations
├── scripts/                 Operational scripts
├── web/                     React + TypeScript + Vite web application
├── docker-compose.yml       Local/production-style container stack
├── Dockerfile               Multi-stage StockPilot image
└── Makefile                 Common development commands
```

## Requirements

For the backend and web application:

- Go 1.26 or newer compatible toolchain.
- Node.js 22 or newer.
- PostgreSQL 18 for the currently tested database configuration.
- Docker with Compose is optional but recommended for the quickest setup.

For Android development:

- JDK 17.
- Android SDK 36.
- Gradle 9.5.

See [`android/README.md`](android/README.md) for Android-specific setup and security notes.

For the browser companion, Node.js 22 is sufficient for validation and tests. See [`extension/README.md`](extension/README.md).

## Quick start with Docker Compose

1. Create a local environment file:

   ```bash
   cp .env.example .env
   ```

2. Replace the placeholder values in `.env`, especially:

   - `POSTGRES_PASSWORD`
   - `DATABASE_URL_DOCKER`
   - `SESSION_SECRET`

   `SESSION_SECRET` must be at least 32 bytes of high-entropy random data. Do not commit real secrets.

3. Start the stack:

   ```bash
   docker compose --env-file .env up --build -d
   ```

4. Check readiness:

   ```bash
   curl http://localhost:8080/readyz
   ```

5. Bootstrap the first administrator from a shell with the same database configuration:

   ```bash
   export STOCKPILOT_ADMIN_EMAIL="admin@example.com"
   export STOCKPILOT_ADMIN_NAME="StockPilot Admin"
   export STOCKPILOT_ADMIN_PASSWORD="use-a-long-unique-password"
   go run ./cmd/admin bootstrap
   ```

6. Open `http://localhost:8080` after the web application has been built into `web/dist`, or use the development workflow below.

## Local development

Start PostgreSQL:

```bash
make db-up
```

Copy `.env.example` to `.env` and export/load the required environment variables for your shell. The default local database URL is:

```text
postgres://stockpilot:stockpilot@localhost:5432/stockpilot?sslmode=disable
```

For a local-only development database, make sure the database credentials match the running PostgreSQL instance.

Build the web application:

```bash
make web-install
make web-build
```

Run the server:

```bash
make dev
```

The server exposes:

- `GET /healthz` — process health.
- `GET /readyz` — database readiness.
- `GET /api/v1/meta` — public StockPilot metadata.
- `/api/v1/auth/*` — authentication endpoints.
- `/api/v1/categories`, `/suppliers`, `/products` — catalog APIs.
- `/api/v1/warehouses`, `/locations`, `/lots`, `/inventory/*` — inventory APIs.
- `/api/v1/orders/*` — purchasing APIs.
- `/api/v1/reports/*` — read-only reporting APIs, including warehouse/location valuation.
- `/api/v1/users/*` — administrator-only user management APIs.

Authenticated mutation requests require the `X-StockPilot-CSRF: 1` confirmation header.

## Common commands

```bash
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
make db-up
make db-down
make migrate
make backup
make clean
```

Android commands use the `GRADLE` variable when a non-default Gradle executable is required:

```bash
make android-build GRADLE=/path/to/gradle
```

## Authentication and roles

StockPilot uses server-side sessions represented by the `stockpilot_session` cookie. The browser session cookie is HttpOnly and SameSite Strict. The native Android client stores the session value encrypted with an AES-GCM key generated in Android Keystore.

The current roles are:

- `admin`
- `manager`
- `operator`
- `viewer`

Permissions are enforced server-side for catalog, inventory, purchase-order, reporting, audit, and user-management operations.

## Security model

Important defaults include:

- request IDs and structured request logging;
- `X-Content-Type-Options: nosniff`;
- restrictive frame, referrer, permissions, and content-security policies;
- explicit CORS origin allow-listing for API requests;
- CSRF confirmation for authenticated mutations;
- HttpOnly/SameSite session cookies;
- secure-cookie support for production deployments;
- release Android cleartext traffic disabled;
- encrypted Android session storage;
- browser companion credentials intentionally not persisted or copied from the web session;
- Docker `no-new-privileges` hardening;
- CodeQL and dependency-update automation.

Production deployments should terminate TLS with a trusted certificate and use HTTPS end to end. Never commit `.env`, database passwords, session secrets, signing keys, or production credentials.

## Web application

The web client is under `web/` and uses React, TypeScript, and Vite.

```bash
cd web
npm install --no-audit --no-fund
npm run build
```

The build performs TypeScript checking before producing `web/dist`.

## Android application

The Android client is a native application rather than a WebView wrapper. It supports Android 8.0+ (`minSdk 26`) and targets Android 16/API 36.

```bash
make android-lint
make android-test
make android-build
```

Debug builds default to `http://10.0.2.2:8080` for Android Emulator development. Release builds do not provide a default server and require HTTPS.

Full details: [`android/README.md`](android/README.md).

## Browser companion

The extension is a Manifest V3 companion for Chromium-compatible browsers. Its current scope is intentionally small: configure one StockPilot origin, request access only to that host, check public server readiness/metadata, and open the StockPilot web application.

```bash
make extension-check
make extension-test
```

It stores no password, session cookie, or API credential. Authenticated extension operations are reserved for a dedicated, independently revocable credential design in a later version.

Full details: [`extension/README.md`](extension/README.md).

## Database migrations

Versioned migrations live in `migrations/` and are applied in filename order by the server migration runner when `AUTO_MIGRATE=true`.

For the Docker development database, migrations can also be applied with:

```bash
make migrate
```

Do not edit a migration that has already been applied to shared environments. Add a new migration instead.

## Backups

The repository includes an operational backup script exposed by:

```bash
make backup
```

Set `BACKUP_DIR` to the intended protected destination and test restore procedures before relying on backups in production.

## Quality gates

GitHub Actions currently checks:

- Go module tidiness, gofmt, vet, race-enabled tests, and server build;
- React/TypeScript typecheck and Vite build;
- PostgreSQL migration/readiness smoke testing;
- Android lint, unit tests, and debug APK assembly;
- browser companion JavaScript/manifest validation and unit tests;
- CodeQL analysis for Go and JavaScript/TypeScript.

Dependabot monitors Go modules, web npm packages, Android Gradle dependencies, and GitHub Actions dependencies on a weekly schedule.

## License

StockPilot is licensed under the MIT License. See [`LICENSE`](LICENSE).

## Development continuity

Detailed implementation progress, verification results, known limitations, and the next exact development tasks are maintained in [`what_changed.md`](what_changed.md).
