## MODIFIED Requirements

### Requirement: Skill Executor
The system SHALL safely execute skills of three types: composite, script, and template.

#### Scenario: Execute composite skill
- **WHEN** executing a composite skill
- **THEN** the system SHALL extract the steps array from the definition
- **AND** return an execution plan with step numbers, tool names, and parameters

#### Scenario: Execute script skill
- **WHEN** executing a script skill
- **THEN** the system SHALL validate the script against dangerous patterns
- **AND** create a temporary file via `os.CreateTemp` in the OS temp directory
- **AND** write the script content to the temp file and close it before execution
- **AND** execute it via `sh` with context-based timeout
- **AND** clean up the temporary file after execution via `defer os.Remove`

#### Scenario: Execute template skill
- **WHEN** executing a template skill
- **THEN** the system SHALL parse the template string as a Go text/template
- **AND** execute it with the provided parameters
- **AND** return the rendered output

### Requirement: Executor Initialization
The system SHALL initialize the executor without filesystem side-effects.

#### Scenario: Infallible construction
- **WHEN** `NewExecutor` is called
- **THEN** the system SHALL return an `*Executor` value directly (no error)
- **AND** SHALL NOT create any directories or perform filesystem operations

## REMOVED Requirements

### Requirement: Executor Initialization
**Reason**: Replaced by the infallible constructor above. The `~/.lango/skills/` directory is no longer created; script temp files use the OS temp directory instead.
**Migration**: Callers remove error handling from `NewExecutor` and `NewRegistry` calls. Existing `~/.lango/skills/` directories can be manually deleted.
