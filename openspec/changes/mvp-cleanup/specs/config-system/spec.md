## MODIFIED Requirements

### Requirement: Default values
The configuration system SHALL apply sensible defaults for all non-credential fields. The minimum viable configuration SHALL require only: `agent.provider`, `providers.<name>.type`, `providers.<name>.apiKey`, and one channel's `enabled: true` + token. All other fields SHALL have defaults:
- `server.host`: `"localhost"`
- `server.port`: `18789`
- `server.httpEnabled`: `true`
- `server.wsEnabled`: `true`
- `session.databasePath`: `"~/.lango/data.db"`
- `session.maxHistoryTurns`: `100`
- `logging.level`: `"info"`
- `logging.format`: `"console"`
- `agent.maxTokens`: `4096`
- `agent.temperature`: `0.7`
- `tools.exec.defaultTimeout`: `30s`
- `tools.exec.allowBackground`: `true`
- `tools.filesystem.maxReadSize`: `1048576` (1MB)
- `tools.browser.headless`: `true`
- `tools.browser.sessionTimeout`: `5m`

#### Scenario: Missing optional field
- **WHEN** a configuration field is not specified
- **THEN** the system SHALL use the default value listed above
- **THEN** no error or warning SHALL be emitted for missing optional fields

#### Scenario: Minimal configuration startup
- **WHEN** config contains only `agent.provider`, one provider entry with `type` and `apiKey`, and one channel with `enabled: true` and token
- **THEN** the application SHALL start successfully with all defaults applied

### Requirement: Configuration validation
The configuration system SHALL validate that at least one provider is configured with a non-empty `apiKey` or valid OAuth token. It SHALL validate that `agent.provider` references an existing key in the `providers` map. It SHALL NOT require `agent.apiKey` (this field no longer exists).

#### Scenario: Valid configuration
- **WHEN** config has `agent.provider: "google"` and `providers.google.type: "gemini"` with a valid `apiKey`
- **THEN** validation SHALL pass

#### Scenario: Invalid configuration
- **WHEN** config has `agent.provider: "google"` but no `google` key in `providers` map
- **THEN** validation SHALL fail with a clear error message

## REMOVED Requirements

### Requirement: agent.apiKey field
**Reason**: Credentials are centralized in the `providers` map. The `agent.apiKey` field created duplication and confusion.
**Migration**: Move API key to `providers.<name>.apiKey`. The `agent.provider` field references the provider by name.
