## MODIFIED Requirements

### Requirement: AgentSpec registry defines sub-agent identity and routing metadata
The system SHALL define an `AgentSpec` type with fields: Name, Description, Instruction, Prefixes, Keywords, Accepts, Returns, CannotDo, and AlwaysInclude. A `var agentSpecs` registry SHALL contain specs for all 6 sub-agents in creation order.

#### Scenario: AgentSpec type has all required fields
- **WHEN** the AgentSpec type is defined
- **THEN** it SHALL include Name (string), Description (string), Instruction (string), Prefixes ([]string), Keywords ([]string), Accepts (string), Returns (string), CannotDo ([]string), and AlwaysInclude (bool)

#### Scenario: Registry contains exactly 6 specs
- **WHEN** agentSpecs is initialized
- **THEN** it SHALL contain specs for operator, navigator, vault, librarian, planner, and chronicler in that order

#### Scenario: Each spec has unique name
- **WHEN** agentSpecs is iterated
- **THEN** no two specs SHALL have the same Name

#### Scenario: Librarian prefixes include all skill management tools
- **WHEN** the librarian spec's Prefixes are checked
- **THEN** they SHALL include `create_skill`, `list_skills`, and `import_skill`

#### Scenario: import_skill routes to librarian
- **WHEN** a tool named `import_skill` is partitioned via `PartitionTools`
- **THEN** it SHALL be assigned to the Librarian role tool set

#### Scenario: import_skill has capability description
- **WHEN** `toolCapability` is called with a tool name starting with `import_skill`
- **THEN** it SHALL return a non-empty capability description
