## REMOVED Requirements

### Requirement: Provider Configuration
**Reason**: OAuth authentication path removed. Replaced with API key-only configuration.
**Migration**: Remove `clientId`, `clientSecret`, and `scopes` fields from provider entries. Use `apiKey` with `${ENV_VAR}` references.

## MODIFIED Requirements

### Requirement: Provider Configuration
The system SHALL allow configuring AI providers with an API key. All provider credentials SHALL be configured in the `providers` map.

#### Scenario: Provider with API Key
- **WHEN** `lango.json` includes a provider with `apiKey`
- **THEN** system initializes the provider using the API key directly

#### Scenario: Provider with environment variable reference
- **WHEN** `apiKey` contains `${ENV_VAR}` pattern
- **THEN** system expands the environment variable before use
- **AND** `lango doctor` reports the key as securely configured

#### Scenario: Provider with plaintext API key
- **WHEN** `apiKey` contains a literal key (not an `${ENV_VAR}` reference)
- **THEN** system initializes the provider using the key
- **AND** `lango doctor` warns the user to use environment variables or encrypted profiles

#### Scenario: Provider without credentials
- **WHEN** `apiKey` is empty or missing
- **THEN** system logs a warning during initialization

#### Scenario: Missing provider configuration
- **WHEN** `agent.provider` is set to "google"
- **BUT** `providers.google` is missing or empty
- **THEN** system fails to start with a configuration error
