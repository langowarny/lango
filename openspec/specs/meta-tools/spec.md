## ADDED Requirements

### Requirement: Knowledge Management Tools
The system SHALL provide agent-facing tools for managing the knowledge base.

#### Scenario: save_knowledge tool
- **WHEN** the agent invokes `save_knowledge` with key, category, content, and optional tags/source
- **THEN** the system SHALL persist the knowledge entry via the Store
- **AND** create an audit log entry with action "knowledge_save"
- **AND** return a success status with the key

#### Scenario: search_knowledge tool
- **WHEN** the agent invokes `search_knowledge` with a query and optional category
- **THEN** the system SHALL search knowledge entries via the Store
- **AND** return matching results with count

### Requirement: Learning Management Tools
The system SHALL provide agent-facing tools for managing learned patterns.

#### Scenario: save_learning tool
- **WHEN** the agent invokes `save_learning` with trigger, fix, and optional error_pattern/diagnosis/category
- **THEN** the system SHALL persist the learning entry via the Store
- **AND** create an audit log entry with action "learning_save"
- **AND** return a success status

#### Scenario: search_learnings tool
- **WHEN** the agent invokes `search_learnings` with a query and optional category
- **THEN** the system SHALL search learning entries via the Store
- **AND** return matching results with count

### Requirement: Skill Management Tools
The system SHALL provide agent-facing tools for creating and listing skills.

#### Scenario: create_skill tool
- **WHEN** the agent invokes `create_skill` with name, description, type, and definition (JSON string)
- **THEN** the system SHALL parse the definition JSON
- **AND** create the skill via the Registry
- **AND** if auto-approve is enabled, SHALL activate the skill immediately
- **AND** create an audit log entry with action "skill_create"
- **AND** return the skill status ("draft" or "active")

#### Scenario: list_skills tool
- **WHEN** the agent invokes `list_skills`
- **THEN** the system SHALL return all active skills with their metadata

### Requirement: Tool Learning Wrapper
The system SHALL wrap existing tool handlers to feed execution results into the learning engine.

#### Scenario: Wrap tool with learning
- **WHEN** `wrapWithLearning` is called on a tool
- **THEN** the system SHALL return a new tool with the same name, description, and parameters
- **AND** the wrapped handler SHALL call the original handler first
- **AND** then call `engine.OnToolResult` with the tool name, params, result, and error
- **AND** return the original result and error unchanged
