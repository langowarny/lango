## Why

AI providers (Google, GitHub) may ban accounts that use OAuth for programmatic API access. Inspired by NanoClaw's philosophy of simplicity and environment-based credential management, we remove OAuth provider login and strengthen API key-based authentication security.

## What Changes

- **BREAKING**: Remove `lango login [provider]` CLI command and OAuth authentication flow
- **BREAKING**: Remove `clientId`, `clientSecret`, `scopes` fields from `ProviderConfig` struct
- Remove OAuth token storage/refresh from Supervisor
- Remove OAuth-related imports (`golang.org/x/oauth2` endpoints) from Supervisor
- Add `APIKeySecurityCheck` to `lango doctor` to detect plaintext API keys
- Update OpenSpec specs to reflect OAuth removal and API key best practices

## Capabilities

### New Capabilities
- `apikey-security-check`: Doctor diagnostic check that detects plaintext API keys and recommends environment variable references or encrypted profiles

### Modified Capabilities
- `oauth-login`: Marked as REMOVED — OAuth provider login is no longer supported
- `agent-provider-config`: OAuth scenarios removed, API key is now the sole authentication method for providers

## Impact

- **CLI**: `lango login` command no longer exists (breaking for users relying on OAuth flow)
- **Config**: `ProviderConfig` struct loses 3 fields; existing configs with these fields will silently ignore them (Go JSON unmarshal behavior)
- **Supervisor**: OAuth fallback path removed; providers without API keys get a warning log
- **Dependencies**: `golang.org/x/oauth2` remains in `go.mod` (used by Gateway OIDC, separate system)
- **Doctor**: New "API Key Security" check added to diagnostics
