# StockPilot Browser Companion

This directory contains the browser-extension foundation for the next StockPilot client version. It is a Manifest V3 extension designed for Chromium-compatible browsers.

## Current scope

The current companion deliberately has a small permission and security surface:

- configure one StockPilot deployment origin;
- request host access only after the user selects that server;
- check the public `/api/v1/meta` and `/readyz` endpoints;
- display server/version/readiness state;
- open the configured StockPilot web application;
- store only the configured server URL and extension schema version;
- store no passwords, session cookies, or API credentials.

Authenticated inventory actions are intentionally not included in this preparation version. A later version should use an extension-specific authentication design rather than copying browser session cookies into extension storage.

## Local validation

The extension has no package dependencies. Node.js 22 or newer is sufficient for its static checks and tests:

```bash
cd extension
npm run check
npm test
```

## Load unpacked in Chromium

1. Open the browser's extensions management page.
2. Enable developer mode.
3. Choose **Load unpacked**.
4. Select the repository's `extension/` directory.
5. Open the StockPilot toolbar popup.
6. Enter the origin of a running StockPilot deployment, such as `https://stockpilot.example.com`.
7. Approve access to that host when the browser asks.

For local development, `http://localhost:8080` is accepted. Production deployments should use HTTPS.

## Permissions

### `storage`

Used for the configured StockPilot server URL and the extension schema version only.

### `optional_host_permissions`

The manifest declares HTTP and HTTPS hosts as optional. The popup requests access only to the configured scheme and host after an explicit **Save & connect** action. This avoids granting access to every website at installation time.

Chrome match-pattern host permissions are host-scoped rather than port-scoped. Therefore a server such as `http://localhost:8080` requests the `http://localhost/*` host pattern while still connecting to the exact saved origin, including port `8080`.

## Security boundaries

- Extension pages use Manifest V3's self-only script content security policy.
- No remote JavaScript or inline script execution is used.
- The server URL parser rejects embedded usernames/passwords, non-HTTP schemes, paths, query strings, and fragments.
- Public health checks use `credentials: "omit"`.
- No StockPilot authentication token is currently read or persisted by the extension.
- Host permission is optional and requested from a direct user gesture.
- The configured URL is normalized to a deployment origin before persistence.

## Files

- `manifest.json` — Manifest V3 metadata and minimal permissions.
- `src/background.js` — versioned service-worker bootstrap.
- `src/popup.html` — accessible toolbar popup markup.
- `src/popup.css` — compact light/dark popup presentation.
- `src/popup.js` — permission, storage, health-check, and launcher flow.
- `src/url.js` — server-origin validation and host-pattern generation.
- `test/url.test.mjs` — URL and permission-scope unit tests.

## CI

`.github/workflows/extension.yml` validates the manifest and JavaScript syntax and runs the zero-dependency Node test suite. Pull requests that touch `extension/**` run this quality gate before merge.

## Next authenticated version

The recommended next step is a dedicated extension credential flow with independently revocable credentials, explicit scopes, expiry, and server-side auditability. That work should be implemented together in the backend and extension instead of reusing or scraping the normal web session cookie.
