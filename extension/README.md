# StockPilot Browser Companion

This directory contains the browser-extension foundation for the next StockPilot client version. It is a Manifest V3 extension designed for Chromium-compatible browsers.

## Current scope

The companion deliberately keeps a small permission and security surface:

- configure one StockPilot deployment origin;
- request host access only after the user selects that server;
- check the public `/api/v1/meta` and `/readyz` endpoints;
- display server/version/readiness state;
- scan supported barcode/QR values with camera or manual fallback;
- open product lookup or a selected inventory workflow with the scanned value;
- store only the configured server URL and extension schema version;
- store no passwords, session cookies, or API credentials.

The extension never performs an authenticated inventory mutation itself. It hands a barcode and workflow choice to the normal StockPilot web application, where the existing signed-in session, RBAC, CSRF protection, validation, confirmation, and server-side inventory transaction remain authoritative.

## Inventory handoff

After a scan, the popup offers:

- **Product lookup** — opens StockPilot with `barcode` only.
- **Stock in** — opens with `barcode` and `inventoryOperation=stock_in`.
- **Stock out** — opens with `barcode` and `inventoryOperation=stock_out`.
- **Adjust stock** — opens with `barcode` and `inventoryOperation=adjustment`.
- **Transfer stock** — opens with `barcode` and `inventoryOperation=transfer`.

The web client resolves the barcode through its authenticated product lookup API and preselects the matching product. The inventory screen still requires the user to choose/confirm quantity and location details and explicitly submit the operation.

Query parameters are consumed and removed from the visible browser URL after the initial handoff, so the barcode and workflow selection are not retained in navigation state.

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
- Inventory handoff is navigation-only; the extension cannot create, edit, receive, transfer, or adjust stock directly.

## Files

- `manifest.json` — Manifest V3 metadata and minimal permissions.
- `src/background.js` — versioned service-worker bootstrap.
- `src/popup.html` — accessible toolbar popup markup and inventory action chooser.
- `src/popup.css` — compact light/dark popup presentation.
- `src/popup.js` — permission, storage, health-check, scanning, and workflow handoff flow.
- `src/scanner.js` — browser-native barcode/QR scanner.
- `src/url.js` — server-origin validation, host-pattern generation, and safe handoff URL construction.
- `test/url.test.mjs` — URL and permission-scope unit tests.

## CI

`.github/workflows/extension.yml` validates the manifest and JavaScript syntax and runs the zero-dependency Node test suite. Pull requests that touch `extension/**` run this quality gate before merge.

## Authentication boundary

The companion does not copy or scrape the normal StockPilot web session cookie. The handoff opens the configured StockPilot origin and lets the web application establish and enforce its own authenticated state.

A future dedicated extension credential flow should use independently revocable credentials, explicit scopes, expiry, and server-side auditability. That work should be implemented together in the backend and extension instead of reusing the normal web session cookie.
