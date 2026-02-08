## ADDED Requirements

### Requirement: Configuration file loading
The system SHALL load configuration from JSON or YAML files with a priority chain.

#### Scenario: Load JSON config
- **WHEN** a lango.json file exists in the config directory
- **THEN** the configuration SHALL be loaded from that file

#### Scenario: Load YAML config
- **WHEN** a lango.yaml file exists (no JSON)
- **THEN** the configuration SHALL be loaded from YAML

#### Scenario: Config file priority
- **WHEN** multiple config files exist
- **THEN** JSON SHALL take precedence over YAML

### Requirement: Environment variable substitution
The system SHALL substitute environment variables in configuration values.

#### Scenario: Environment variable in value
- **WHEN** a config value contains ${VAR_NAME}
- **THEN** it SHALL be replaced with the environment variable value

#### Scenario: Missing environment variable
- **WHEN** a referenced environment variable is not set
- **THEN** an error SHALL be logged and default used if available

### Requirement: Configuration validation
The system SHALL validate configuration against a schema before use.

#### Scenario: Valid configuration
- **WHEN** configuration matches the expected schema
- **THEN** the configuration SHALL be accepted

#### Scenario: Invalid configuration
- **WHEN** configuration has missing required fields or wrong types
- **THEN** a validation error SHALL be returned with details

### Requirement: Default values
The system SHALL apply sensible defaults for optional configuration fields.

#### Scenario: Missing optional field
- **WHEN** an optional field is not specified
- **THEN** the documented default value SHALL be used

### Requirement: Runtime configuration updates
The system SHALL support reloading configuration without full restart.

#### Scenario: Config file change
- **WHEN** the configuration file is modified
- **THEN** the system MAY reload affected components

#### Scenario: API config update
- **WHEN** configuration is updated via the Gateway API
- **THEN** the changes SHALL take effect for new operations
