# StockPilot Android

StockPilot includes a native Kotlin Android client under `android/`. The Android app talks directly to the existing StockPilot HTTP API; it is not a WebView wrapper.

## Supported Android versions

- `minSdk`: 26 (Android 8.0)
- `targetSdk`: 36 (Android 16)
- `compileSdk`: 36
- JDK: 17
- Android Gradle Plugin: 9.3.1
- Gradle: 9.5.0

## Build

Open the `android/` directory in Android Studio, or use a Gradle 9.5 installation from the repository root:

```bash
make android-lint
make android-test
make android-build
```

The debug APK is written to `android/app/build/outputs/apk/debug/app-debug.apk`.

## Local development

1. Start StockPilot on the host machine at port `8080`.
2. Launch the Android app in the Android Emulator.
3. Debug builds default to `http://10.0.2.2:8080`, which maps the emulator to the host machine.
4. Sign in with a StockPilot user created through the administrator bootstrap command.

Debug builds permit cleartext traffic so local development servers can be reached. Release builds disable cleartext traffic and require an `https://` server URL.

## Authentication and security

- The app uses the backend's `stockpilot_session` cookie session.
- The cookie value is encrypted at rest using an AES-GCM key generated inside Android Keystore.
- The app sends `X-StockPilot-CSRF: 1` for authenticated mutations, matching the backend's CSRF requirement.
- Invalid or expired sessions are cleared after an HTTP 401 response.
- Android backup is disabled for the app.
- Release networking trusts system certificate authorities and rejects cleartext HTTP.
- The release build has no default server URL; the user must enter the deployed HTTPS StockPilot endpoint.

## Current native screens

The Android application includes:

- server-aware native sign-in;
- secure session resume;
- product count and product overview;
- low-stock overview;
- purchase-order overview;
- account and role information;
- refresh and sign-out actions;
- light and dark system themes;
- edge-to-edge system inset handling;
- accessible labels, live error announcements, and autofill hints.

## CI

`.github/workflows/android.yml` installs Android API 36 and Gradle 9.5, then runs:

```bash
gradle :app:lintDebug :app:testDebugUnitTest :app:assembleDebug --stacktrace
```

A successful CI run uploads the debug APK as the `stockpilot-android-debug` artifact.

## Production release notes

Before publishing a release build:

1. Configure a production signing key outside the repository.
2. Build an Android App Bundle (`:app:bundleRelease`).
3. Connect only to a trusted HTTPS StockPilot deployment.
4. Keep signing credentials and server secrets out of Git and local source files.
5. Validate the signed build on physical devices and current Android emulators.
