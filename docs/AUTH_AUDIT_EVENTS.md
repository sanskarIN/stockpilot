# Authentication Audit Events

StockPilot records security-relevant authentication events without storing passwords, session tokens, cookie values, or other authentication secrets.

## Event types

- `auth.login.success` — authenticated login completed; metadata contains the user ID.
- `auth.login.failure` — login rejected; metadata contains a coarse reason only.
- `auth.session.missing` — a protected request arrived without a session cookie.
- `auth.session.invalid` — a supplied session could not be resolved.
- `auth.logout.success` — an authenticated session was explicitly logged out.
- `auth.logout.failure` — logout could not revoke the requested session.
- `auth.user.role_changed` — an administrator changed another user's role.
- `auth.user.activated` / `auth.user.deactivated` — account activation state changed.
- `auth.user.created` — a new account was created.

The audit stream remains append-only. Authentication metadata is intentionally minimal so auditability does not become credential leakage.
