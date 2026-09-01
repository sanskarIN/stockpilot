# StockPilot Architecture

StockPilot is organized as a modular server with thin adapters and separate client applications.

## Backend boundaries

`internal/domain` contains business entities and validation rules. It should not import HTTP or PostgreSQL packages.

`internal/repository` contains interfaces used by application/HTTP layers. Interfaces describe business capabilities rather than database details.

`internal/postgres` implements repository interfaces and owns SQL, transactions, query mapping, and migration integration.

`internal/httpapi` maps HTTP requests to domain/application operations, applies authentication and authorization middleware, validates request input, and emits stable JSON errors.

`internal/auth` owns password and session primitives. Secrets and session-token values must not leak through logs or error responses.

`internal/config` parses and validates deployment configuration.

## Data flow

```text
Web / Android / Extension
          |
          v
      HTTP API
          |
    auth + RBAC + CSRF
          |
          v
 Repository contracts
          |
          v
   PostgreSQL adapters
          |
          v
       PostgreSQL
```

## Mutation safety

Inventory mutations must be transactional. Reads used to calculate stock availability or replenishment must use consistent database semantics appropriate to the operation. Any new mutation should document its concurrency behavior and failure mode.

## Client responsibilities

The web client owns presentation, form validation for user feedback, and API orchestration. The server remains authoritative for validation and permissions.

The Android client stores only the minimum session state required for authenticated requests and uses HTTPS in release builds.

The browser companion intentionally keeps its current capability narrow and does not persist web passwords, session cookies, or long-lived API credentials.

## Change discipline

When adding a capability, update the domain model first, then the repository contract and persistence implementation, then the HTTP contract, then client integrations and tests. Keep cross-layer changes in reviewable commits.
