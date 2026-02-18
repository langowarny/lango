## Why

The current configuration structure (`lango.json`) duplicates credential management: API keys can be set in both `agent.apiKey` and `providers.<name>.apiKey`. This ambiguity causes confusion and inconsistency, especially with the introduction of OAuth credentials. Centralizing credentials in the `providers` map simplifies configuration and security management.

## What Changes

- **BREAKING**: Remove `agent.apiKey` from configuration.
- Update `agent.provider` to be a reference to a key in the `providers` map.
- Centralize all credentials (API keys, OAuth tokens) within the `providers` configuration section.
- Update `lango.example.json` to reflect the new canonical structure.
- Update Supervisor and Config Loader to enforce this structure.

## Capabilities

### New Capabilities
- None (Structural Refactoring)

### Modified Capabilities
- `agent-provider-config`: Enforce centralized credential management in `providers` map and remove legacy `agent.apiKey` support.

## Impact

- **Configuration**: `lango.json` schema changes. Users must migrate credentials to `providers`.
- **Code**: `internal/config`, `internal/supervisor` updates.
- **Documentation**: `lango.example.json` and README updates.
