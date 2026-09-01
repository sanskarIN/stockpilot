# Contributing to StockPilot

Thanks for contributing to StockPilot.

## Development workflow

1. Create a focused branch from `main` using a descriptive name such as `feat/inventory-receiving` or `fix/expiry-validation`.
2. Keep changes small and cohesive. Prefer separate commits for domain rules, persistence, HTTP/API changes, UI changes, tests, and documentation when they can be reviewed independently.
3. Add or update tests with every behavior change.
4. Run the relevant quality gates before opening a pull request.

## Local checks

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
```

Database-backed changes should also be tested against PostgreSQL using the repository's Compose workflow.

## Commit messages

Use imperative, scoped Conventional Commit-style messages where practical:

```text
feat(inventory): add guided stock receiving
fix(auth): reject expired session tokens
 test(web): cover catalog validation
 docs: document restore drill
 chore(ci): tighten dependency checks
```

Avoid mixing unrelated refactors into feature commits.

## Pull requests

A pull request should explain what changed, why it changed, how it was tested, and any migration or security implications. Changes that affect an API contract should include compatibility notes.

Never include real credentials, private keys, customer data, or production dumps in commits, issues, or pull-request artifacts.
