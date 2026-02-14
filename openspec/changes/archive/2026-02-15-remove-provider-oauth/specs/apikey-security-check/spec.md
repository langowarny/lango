## ADDED Requirements

### Requirement: API Key Security Diagnostic Check
The system SHALL provide a diagnostic check in `lango doctor` that detects plaintext API keys and recommends secure alternatives.

#### Scenario: All keys use environment variable references
- **WHEN** all configured providers use `${ENV_VAR}` pattern for `apiKey`
- **THEN** the check SHALL report `StatusPass` with message "All API keys use environment variable references"

#### Scenario: Plaintext API key detected
- **WHEN** one or more providers have an `apiKey` that does not match `${...}` pattern
- **THEN** the check SHALL report `StatusWarn`
- **AND** the message SHALL list the provider IDs with plaintext keys
- **AND** details SHALL recommend using environment variable references or encrypted profiles

#### Scenario: No providers configured
- **WHEN** the `providers` map is empty
- **THEN** the check SHALL report `StatusSkip`

#### Scenario: No API keys set
- **WHEN** providers exist but all have empty `apiKey` values
- **THEN** the check SHALL report `StatusSkip`

#### Scenario: Configuration not loaded
- **WHEN** config is nil
- **THEN** the check SHALL report `StatusSkip`
