# StockPilot v0.1.1 Release Artifacts

## Artifact policy

Every distributable artifact must be produced from the exact `v0.1.1` tag. Never mix artifacts from different commits in one release.

## Required publication set

### Source

- GitHub-generated source archive for tag `v0.1.1`.
- Repository source remains the canonical source distribution.

### Server

If a server binary is distributed, publish one binary per supported target documented by the project. Record the target tuple, build command, and exact commit SHA in the release description.

### Web

Publish the production web build only when the deployment process supports a reproducible artifact. Record the build command and exact commit SHA.

### Android

If an APK is published, build it from the exact tag with the release configuration. Verify HTTPS enforcement and encrypted session storage before attaching it.

### Browser companion

If the extension is distributed as an archive, build it from the exact tag and verify manifest validation, host-permission scope, and absence of persisted credentials.

## Checksums

Create a UTF-8 text file named `SHA256SUMS-v0.1.1.txt` containing one SHA-256 digest and filename per published binary/archive, using the conventional two-space separator:

```text
<64 lowercase hexadecimal characters>  <filename>
```

Do not include credentials, private deployment configuration, customer data, database dumps, or development secrets in release artifacts.

## Verification record

Before publication record:

- Tag: `v0.1.1`
- Release commit SHA: `<exact verified SHA>`
- Build timestamp: `<UTC timestamp>`
- Go version: `<version>` where applicable
- Node.js version: `<version>` where applicable
- Android/Gradle toolchain: `<versions>` where applicable
- SHA-256 manifest: `SHA256SUMS-v0.1.1.txt`

## Acceptance criteria

- Every artifact maps to the exact release commit.
- Checksums verify after download.
- No debug-only or development-secret configuration is included.
- Post-release smoke tests pass using the published artifacts.
