## ADDED Requirements

### Requirement: Providers Configuration Section
The system SHALL support a `providers` section in the configuration file to define multiple AI providers.

#### Scenario: Provider specific settings
- **WHEN** `providers` map is present in config
- **THEN** it SHALL map provider IDs (e.g., "openai", "anthropic") to their specific settings
- **AND** settings SHALL include `apiKey`, `baseUrl`, and provider-specific fields

#### Scenario: Fallback configuration
- **WHEN** `agent.fallbacks` list is present
- **THEN** it SHALL define an ordered list of fallback models
- **AND** each fallback SHALL specify `provider` and `model`

### Requirement: Provider Selection
The system SHALL allow selecting the active provider and model.

#### Scenario: Explicit provider selection
- **WHEN** `agent.provider` is set in config
- **THEN** the system SHALL use that provider for agent operations

#### Scenario: Default provider
- **WHEN** `agent.provider` is missing but `providers` has entries
- **THEN** the system SHALL adhere to a documented default behavior or return an error if ambiguous

## MODIFIED Requirements

### Requirement: Configuration validation
The system SHALL validate configuration against a schema before use.

#### Scenario: Valid configuration
- **WHEN** configuration matches the expected schema
- **THEN** the configuration SHALL be accepted
- **AND** `providers` section SHALL be validated if present
- **AND** `agent.provider` SHALL match a valid provider key if specified

#### Scenario: Invalid configuration
- **WHEN** configuration has missing required fields or wrong types
- **THEN** a validation error SHALL be returned with details
