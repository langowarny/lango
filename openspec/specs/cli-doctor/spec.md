## ADDED Requirements

### Requirement: Doctor Command Entry Point
The system SHALL provide a `lango doctor` command that runs diagnostic checks on the Lango installation and configuration.

#### Scenario: Running doctor command
- **WHEN** user executes `lango doctor`
- **THEN** system runs all diagnostic checks and displays results in TUI format

#### Scenario: Running doctor with JSON output
- **WHEN** user executes `lango doctor --json`
- **THEN** system outputs results as JSON to stdout without TUI formatting

### Requirement: Verify all providers
The command SHALL verify the status of every provider defined in the `providers` configuration map.

#### Scenario: Multiple providers configured
- **WHEN** `lango.json` contains both "openai" and "anthropic" in `providers`
- **THEN** the doctor output includes checks for both "OpenAI" and "Anthropic"

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

### Requirement: Legacy API Key Verification
The system SHALL verify the legacy API key configuration ONLY IF no modern providers are configured.

#### Scenario: API key configured via environment (fallback)
- **WHEN** no providers configured AND GOOGLE_API_KEY environment variable is set
- **THEN** check passes with warning "Implicit Gemini config found"

#### Scenario: API key missing
- **WHEN** no providers configured AND no API key found
- **THEN** check fails with message "No AI providers configured"

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

### Requirement: Security Provider Mode Check
The system SHALL check security provider configuration and provide appropriate warnings.

#### Scenario: Local provider warning
- **WHEN** security.signer.provider is "local"
- **THEN** doctor displays WARNING: "Using LocalCryptoProvider (dev/test only). For production, use RPCProvider with Companion app."
- **AND** check status is "warn"

#### Scenario: RPC provider configured
- **WHEN** security.signer.provider is "rpc"
- **THEN** doctor displays "Security: Using RPCProvider (production mode)"
- **AND** check status is "pass"

#### Scenario: Companion connectivity with RPC
- **WHEN** security.signer.provider is "rpc"
- **AND** no companion is connected
- **THEN** doctor displays WARNING: "RPCProvider configured but no companion connected. Crypto operations will fail."
- **AND** check status is "warn"

### Requirement: Passphrase Checksum Integrity Check
The system SHALL verify that passphrase checksum exists when local provider is configured.

#### Scenario: Checksum present
- **WHEN** local provider is configured
- **AND** security_config contains valid checksum
- **THEN** doctor displays "Passphrase checksum: configured"
- **AND** check status is "pass"

#### Scenario: Checksum missing
- **WHEN** local provider is configured
- **AND** security_config has no checksum
- **THEN** doctor displays WARNING: "Passphrase not initialized. Run application to set up."
- **AND** check status is "warn"

### Requirement: Graceful Failure
The `doctor` command MUST NOT crash or return obscure system errors (like "out of memory") when encountering an encrypted database without credentials.

#### Scenario: Encrypted DB without passphrase
- **WHEN** opening session store returns "out of memory" or encryption error
- **THEN** doctor displays WARNING: "Session database is encrypted and locked"
- **AND** check status is "warn"

### Requirement: Clear Feedback
If the database is locked, `doctor` SHOULD suggest how to unlock (e.g., set `LANGO_PASSPHRASE`).

### Requirement: Security Config Validation
The `Security Configuration` check MUST verify that the `security` block in `lango.json` is valid according to the new schema.

### Requirement: Documentation Alignment
- **Example Config**: `lango.example.json` MUST include valid `providers` map and `security` block.
- **README**: MUST explain multiple providers and TUI configuration.
