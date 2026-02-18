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

### Requirement: Skill Executor
The system SHALL safely execute skills of three types: composite, script, and template.

#### Scenario: Execute composite skill
- **WHEN** executing a composite skill
- **THEN** the system SHALL extract the steps array from the definition
- **AND** return an execution plan with step numbers, tool names, and parameters

#### Scenario: Execute script skill
- **WHEN** executing a script skill
- **THEN** the system SHALL validate the script against dangerous patterns
- **AND** write the script to a temporary file in the skills directory
- **AND** execute it via `sh` with context-based timeout
- **AND** clean up the temporary file after execution

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
The system SHALL properly initialize the executor with error handling.

#### Scenario: Skills directory creation
- **WHEN** `NewExecutor` is called
- **THEN** the system SHALL resolve the user home directory and create `~/.lango/skills/`
- **AND** if home directory resolution fails, SHALL return an error
- **AND** if directory creation fails, SHALL return an error
