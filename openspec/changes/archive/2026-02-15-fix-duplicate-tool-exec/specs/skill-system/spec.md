## ADDED Requirements

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

## MODIFIED Requirements

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
