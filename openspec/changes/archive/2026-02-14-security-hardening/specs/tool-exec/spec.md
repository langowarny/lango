## MODIFIED Requirements

### Requirement: Reference token resolution in exec
The exec tool SHALL resolve secret reference tokens in command strings immediately before execution. Resolved values SHALL never be logged or returned to the agent.

#### Scenario: Command with secret reference
- **WHEN** exec is called with command `curl -H "Auth: {{secret:api-key}}" https://api.example.com`
- **AND** the RefStore contains a value for `{{secret:api-key}}`
- **THEN** the token SHALL be replaced with the actual secret value before shell execution
- **AND** the log entry SHALL contain the original command with the unresolved token
- **AND** the BackgroundProcess.Command field SHALL contain the original command with the unresolved token

#### Scenario: Command with decrypt reference
- **WHEN** exec is called with command `echo {{decrypt:uuid-123}}`
- **AND** the RefStore contains a value for `{{decrypt:uuid-123}}`
- **THEN** the token SHALL be replaced with the actual decrypted value before shell execution

#### Scenario: Command with unknown reference
- **WHEN** exec is called with command `echo {{secret:unknown}}`
- **AND** the RefStore does NOT contain a value for `{{secret:unknown}}`
- **THEN** the literal string `{{secret:unknown}}` SHALL be passed to the shell unchanged

#### Scenario: Command without references
- **WHEN** exec is called with a command containing no reference tokens
- **THEN** the command SHALL be executed unchanged

#### Scenario: Reference resolution in PTY mode
- **WHEN** RunWithPTY is called with a command containing reference tokens
- **THEN** tokens SHALL be resolved identically to synchronous execution

#### Scenario: Reference resolution in background mode
- **WHEN** StartBackground is called with a command containing reference tokens
- **THEN** tokens SHALL be resolved identically to synchronous execution
