## Why

API Key authentication is cumbersome for users and requires managing sensitive secrets in plaintext config files. OAuth login provides a more secure, user-friendly, and standardized way to authenticate agents with providers like Google and GitHub.

## What Changes

- Add `lango login [provider]` CLI command for OAuth authentication.
- Implement local loopback server to handle OAuth callbacks.
- Store OAuth tokens securely in `~/.lango/tokens/`.
- Update Supervisor to prioritize OAuth tokens over API keys when initializing providers.
- Support token refresh for continuous agent operation.
- Update `lango.json` configuration schema to include OAuth client details.

## Capabilities

### New Capabilities
- `oauth-login`: CLI command and local server to handle OAuth login flow and token management.

### Modified Capabilities
- `agent-provider-config`: Enhance provider configuration to support OAuth client credentials and token-based initialization alongside API keys.

## Impact

- **CLI**: New `auth` package and commands.
- **Config**: `lango.json` schema updates for `providers`.
- **Supervisor**: Logic update to load and refresh OAuth tokens.
- **Dependencies**: Added `golang.org/x/oauth2` for OAuth flow handling.
