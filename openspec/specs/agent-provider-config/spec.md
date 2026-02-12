## Requirements

### Requirement: Provider Configuration
The system SHALL allow configuring AI providers with either an API key or OAuth credentials. All provider credentials SHALL be configured in the `providers` map.

#### Scenario: Provider with OAuth
- **WHEN** `lango.json` includes a provider with `clientId` and `clientSecret`
- **AND** a valid token exists in `~/.lango/tokens/<provider>.json`
- **THEN** system initializes the provider using the stored access token
- **AND** ignores any `apiKey` if present (or falls back if token invalid/missing)

#### Scenario: Provider with API Key (Legacy)
- **WHEN** `lango.json` includes a provider with `apiKey` only
- **THEN** system initializes the provider using the API key directly

#### Scenario: Provider without credentials
- **WHEN** neither `apiKey` nor valid OAuth token is available
- **THEN** system logs a warning or error during initialization

#### Scenario: Missing provider configuration
- **WHEN** `agent.provider` is set to "google"
- **BUT** `providers.google` is missing or empty
- **THEN** system fails to start with a configuration error

## Deprecated Requirements

### Requirement: Legacy API Key Support
**Reason**: Duplication and ambiguity with `providers` map.
**Migration**: Move `agent.apiKey` to `providers.<agent.provider>.apiKey`.

#### Scenario: Legacy config detected
- **WHEN** user configuration contains `agent.apiKey`
- **THEN** system fails to start (or ignores it with a warning, depending on strictness - we choose fail for clarity)
