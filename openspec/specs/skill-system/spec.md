## ADDED Requirements

### Requirement: Skill Registry
The system SHALL provide a registry for managing reusable skill definitions with lifecycle management.

#### Scenario: Create skill
- **WHEN** `CreateSkill` is called with a valid skill entry
- **THEN** the system SHALL validate the skill type is one of "composite", "script", or "template"
- **AND** validate the skill definition matches the type requirements
- **AND** persist the skill with status "draft"

#### Scenario: Invalid skill type
- **WHEN** `CreateSkill` is called with an unrecognized skill type
- **THEN** the system SHALL return an error

#### Scenario: Composite skill validation
- **WHEN** creating a composite skill
- **THEN** the definition SHALL contain a "steps" array

#### Scenario: Script skill validation
- **WHEN** creating a script skill
- **THEN** the definition SHALL contain a "script" string
- **AND** the script SHALL be validated against dangerous patterns

#### Scenario: Template skill validation
- **WHEN** creating a template skill
- **THEN** the definition SHALL contain a "template" string

#### Scenario: Activate skill
- **WHEN** `ActivateSkill` is called with a draft skill name
- **THEN** the system SHALL set the skill status to "active"

#### Scenario: Load skills on startup
- **WHEN** the registry is initialized
- **THEN** `LoadSkills` SHALL load all active skills from the store

#### Scenario: App tool assembly with knowledge system
- **WHEN** the knowledge system is enabled and tools are assembled in `app.go`
- **THEN** the app SHALL use `LoadedSkills()` to append only dynamic skills
- **AND** SHALL NOT use `AllTools()` which would duplicate base tools already present in the tool list

### Requirement: Loaded Skills Retrieval
The registry SHALL provide a `LoadedSkills()` method that returns only dynamically loaded skill tools, excluding base tools.

#### Scenario: No skills loaded
- **WHEN** `LoadedSkills` is called before any skills are loaded
- **THEN** the system SHALL return an empty slice

#### Scenario: Skills loaded
- **WHEN** `LoadedSkills` is called after skills have been activated
- **THEN** the system SHALL return only the dynamically loaded skill tools
- **AND** the result SHALL NOT include any base tools passed during registry creation

#### Scenario: Concurrent safety
- **WHEN** `LoadedSkills` is called concurrently with `LoadSkills`
- **THEN** access SHALL be protected by a read lock

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

### Requirement: Dangerous Pattern Validation
The system SHALL validate scripts against known dangerous patterns as a defense-in-depth measure.

#### Scenario: Reject dangerous scripts
- **WHEN** a script matches any dangerous pattern
- **THEN** the system SHALL return an error
- **AND** dangerous patterns SHALL include: recursive force delete (`rm -rf /`), fork bombs, pipe-to-shell (`curl|sh`), raw device writes (`>/dev/sd`), filesystem formatting (`mkfs.`), and raw disk copies (`dd if=`)

### Requirement: Skill Builder
The system SHALL provide a builder for constructing skill entries from tool execution traces.

#### Scenario: Build composite skill from steps
- **WHEN** `BuildFromSteps` is called with a name, description, and list of tool steps
- **THEN** the system SHALL construct a SkillEntry of type "composite" with the steps in the definition

#### Scenario: Build script skill
- **WHEN** `BuildScript` is called with a name, description, and script content
- **THEN** the system SHALL construct a SkillEntry of type "script" with the script in the definition

### Requirement: Executor Initialization
The system SHALL initialize the executor without filesystem side-effects.

#### Scenario: Infallible construction
- **WHEN** `NewExecutor` is called
- **THEN** the system SHALL return an `*Executor` value directly (no error)
- **AND** SHALL NOT create any directories or perform filesystem operations
