# Security Policy

## Supported development line

StockPilot is pre-1.0 software. Security fixes are prioritized on the `main` branch and the latest unreleased development line.

## Reporting a vulnerability

Please do not open a public issue for a suspected vulnerability. Report it privately through the repository's GitHub Security Advisories interface when available, or contact the project security contact listed in the repository profile.

Include enough information to reproduce the issue safely, the affected component, expected impact, and any suggested mitigation.

## Secrets

Never commit passwords, session secrets, database credentials, private keys, Android signing material, or production environment files. Use local `.env` files and CI/hosting secret stores instead.

## Security baseline

StockPilot expects:

- HTTPS for production traffic.
- Server-side authorization for every protected operation.
- CSRF confirmation on authenticated browser mutations.
- HttpOnly/SameSite session cookies.
- Encrypted Android session persistence.
- Least-privilege browser-extension host permissions.
- Parameterized PostgreSQL queries.
- Automated dependency and CodeQL checks.

## Dependency updates

Dependency-update pull requests must still be reviewed for compatibility, license changes, and security impact before merge. Do not merge an automated update solely because the dependency is newer.
