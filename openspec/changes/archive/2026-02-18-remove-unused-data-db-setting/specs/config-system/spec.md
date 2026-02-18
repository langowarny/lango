## MODIFIED Requirements

### Requirement: Default values
The configuration system SHALL apply sensible defaults for all non-credential fields. The minimum viable configuration SHALL require only: `agent.provider`, `providers.<name>.type`, `providers.<name>.apiKey`, and one channel's `enabled: true` + token. All other fields SHALL have defaults:
- `server.host`: `"localhost"`
- `server.port`: `18789`
- `server.httpEnabled`: `true`
- `server.wsEnabled`: `true`
- `session.databasePath`: `"~/.lango/lango.db"`
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
- `librarian.enabled`: `false`
- `librarian.observationThreshold`: `2`
- `librarian.inquiryCooldownTurns`: `3`
- `librarian.maxPendingInquiries`: `2`
- `librarian.autoSaveConfidence`: `"high"`

#### Scenario: Missing optional field
- **WHEN** a configuration field is not specified
- **THEN** the system SHALL use the default value listed above
- **THEN** no error or warning SHALL be emitted for missing optional fields

#### Scenario: Session database path defaults to lango.db
- **WHEN** `session.databasePath` is not specified in the configuration
- **THEN** the system SHALL default to `"~/.lango/lango.db"`
- **THEN** standalone CLI commands (doctor, memory list) SHALL open this path as fallback
