## MODIFIED Requirements

### Requirement: Provider Configuration
The system SHALL require all provider credentials to be configured in the `providers` map.

#### Scenario: Agent using configured provider
- **WHEN** `agent.provider` is set to "google"
- **AND** `providers.google` contains valid credentials (apiKey or OAuth)
- **THEN** system initializes the agent using the "google" provider configuration

#### Scenario: Missing provider configuration
- **WHEN** `agent.provider` is set to "google"
- **BUT** `providers.google` is missing or empty
- **THEN** system fails to start with a configuration error

## REMOVED Requirements

### Requirement: Legacy API Key Support
**Reason**: Duplication and ambiguity with `providers` map.
**Migration**: Move `agent.apiKey` to `providers.<agent.provider>.apiKey`.

#### Scenario: Legacy config detected
- **WHEN** user configuration contains `agent.apiKey`
- **THEN** system fails to start (or ignores it with a warning, depending on strictness - we choose fail for clarity)
