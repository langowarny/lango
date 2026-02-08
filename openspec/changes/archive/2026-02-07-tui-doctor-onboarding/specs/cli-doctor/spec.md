## ADDED Requirements

### Requirement: Doctor Command Entry Point
The system SHALL provide a `lango doctor` command that runs diagnostic checks on the Lango installation and configuration.

#### Scenario: Running doctor command
- **WHEN** user executes `lango doctor`
- **THEN** system runs all diagnostic checks and displays results in TUI format

#### Scenario: Running doctor with JSON output
- **WHEN** user executes `lango doctor --json`
- **THEN** system outputs results as JSON to stdout without TUI formatting

### Requirement: Configuration File Check
The system SHALL verify that the configuration file exists and contains valid JSON syntax.

#### Scenario: Valid configuration file
- **WHEN** lango.json exists with valid JSON syntax
- **THEN** check passes with message "Configuration file valid"

#### Scenario: Missing configuration file
- **WHEN** lango.json does not exist
- **THEN** check fails with message "Configuration file not found" and suggestion to run `lango onboard`

#### Scenario: Invalid JSON syntax
- **WHEN** lango.json contains invalid JSON
- **THEN** check fails with specific syntax error location

### Requirement: API Key Verification
The system SHALL verify that the configured API key is present and optionally validate it with the provider.

#### Scenario: API key configured via environment
- **WHEN** GOOGLE_API_KEY environment variable is set
- **THEN** check passes with message "API key configured"

#### Scenario: API key configured in config file
- **WHEN** agent.apiKey is set in lango.json
- **THEN** check passes with message "API key configured"

#### Scenario: API key missing
- **WHEN** neither environment variable nor config contains API key
- **THEN** check fails with message "API key not configured"

### Requirement: Channel Token Validation
The system SHALL verify that enabled channel tokens are configured.

#### Scenario: Telegram enabled with token
- **WHEN** channels.telegram.enabled is true AND botToken is set
- **THEN** Telegram channel check passes

#### Scenario: Channel enabled without token
- **WHEN** any channel is enabled but token is missing
- **THEN** check fails with specific channel and missing token field

### Requirement: Session Database Check
The system SHALL verify that the session database is accessible.

#### Scenario: Database file exists and is writable
- **WHEN** session.databasePath points to an accessible SQLite file
- **THEN** check passes with database path displayed

#### Scenario: Database path not writable
- **WHEN** database path directory is not writable
- **THEN** check fails with permission error

### Requirement: Server Port Check
The system SHALL verify that the configured server port is available.

#### Scenario: Port available
- **WHEN** configured port (default 18789) is not in use
- **THEN** check passes with "Port 18789 available"

#### Scenario: Port in use
- **WHEN** configured port is already bound by another process
- **THEN** check fails with "Port 18789 in use" and process information if available

### Requirement: Auto-Fix Mode
The system SHALL support a `--fix` flag that attempts to automatically repair common issues.

#### Scenario: Fix creates missing database directory
- **WHEN** `--fix` is provided AND database directory does not exist
- **THEN** system creates the directory and reports "Created ~/.lango directory"

#### Scenario: Fix generates default config
- **WHEN** `--fix` is provided AND lango.json is missing
- **THEN** system creates minimal lango.json and reports "Created default configuration"

### Requirement: Check Result Summary
The system SHALL display a summary of all check results at the end of execution.

#### Scenario: All checks pass
- **WHEN** all diagnostic checks pass
- **THEN** display "Summary: X passed, 0 warnings, 0 errors"

#### Scenario: Mixed results
- **WHEN** some checks pass, some warn, some fail
- **THEN** display accurate counts for each category
