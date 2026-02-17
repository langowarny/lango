## MODIFIED Requirements

### Requirement: Hierarchical agent tree with sub-agents
The system SHALL support a multi-agent mode (`agent.multiAgent: true`) that creates an orchestrator root agent with specialized sub-agents: operator, navigator, vault, librarian, planner, and chronicler. The orchestrator SHALL have NO direct tools (`Tools: nil`) and MUST delegate all tool-requiring tasks to sub-agents.

#### Scenario: Multi-agent mode enabled
- **WHEN** `agent.multiAgent` is true
- **THEN** BuildAgentTree SHALL create an orchestrator that has NO direct tools AND has sub-agents (operator, navigator, vault, librarian, planner, chronicler)

#### Scenario: Orchestrator has no direct tools
- **WHEN** the orchestrator is created
- **THEN** the orchestrator's `Tools` field SHALL be `nil`
- **AND** tools SHALL only be adapted for their respective sub-agents (each tool adapted exactly once)

#### Scenario: Single-agent fallback
- **WHEN** `agent.multiAgent` is false
- **THEN** the system SHALL create a single flat agent with all tools

### Requirement: Tool partitioning by prefix
Tools SHALL be partitioned to sub-agents based on name prefixes with matching order Librarian → Chronicler → Navigator → Vault → Operator → Unmatched: `exec/fs_/skill_` → operator, `browser_` → navigator, `crypto_/secrets_/payment_` → vault, `search_/rag_/graph_/save_knowledge/save_learning/create_skill/list_skills` → librarian, `memory_/observe_/reflect_` → chronicler, unmatched → Unmatched bucket (not assigned to any agent).

#### Scenario: Operator gets shell, file, and skill tools
- **WHEN** tools named `exec_shell`, `fs_read`, `skill_deploy` are registered
- **THEN** they SHALL be assigned to the operator sub-agent

#### Scenario: Navigator gets browser tools
- **WHEN** tools named `browser_navigate`, `browser_screenshot` are registered
- **THEN** they SHALL be assigned to the navigator sub-agent

#### Scenario: Vault gets crypto, secrets, and payment tools
- **WHEN** tools named `crypto_sign`, `secrets_get`, `payment_send` are registered
- **THEN** they SHALL be assigned to the vault sub-agent

#### Scenario: Librarian gets search, RAG, graph, and skill management tools
- **WHEN** tools named `search_web`, `rag_query`, `graph_traverse`, `save_knowledge_item`, `create_skill_x`, `list_skills` are registered
- **THEN** they SHALL be assigned to the librarian sub-agent

#### Scenario: Chronicler gets memory tools
- **WHEN** tools named `memory_store`, `observe_event`, `reflect_summary` are registered
- **THEN** they SHALL be assigned to the chronicler sub-agent

#### Scenario: Unmatched tools tracked separately
- **WHEN** a tool with an unrecognized prefix is present
- **THEN** it SHALL be placed in the Unmatched bucket and NOT assigned to any sub-agent

#### Scenario: Librarian prefix priority over operator
- **WHEN** tools like `save_knowledge_data` or `create_skill_new` are registered
- **THEN** they SHALL match librarian prefixes before reaching operator matching

### Requirement: Conditional Sub-Agent Creation
The `BuildAgentTree` function SHALL create sub-agents data-driven from the agentSpecs registry. Agents with no tools SHALL be skipped unless AlwaysInclude is set. The planner sub-agent SHALL always be created as it is LLM-only.

#### Scenario: All tool categories have tools
- **WHEN** tools exist for operator, navigator, vault, librarian, and chronicler roles
- **THEN** all six sub-agents (operator, navigator, vault, librarian, planner, chronicler) SHALL be created

#### Scenario: Partial tools — only operator and librarian
- **WHEN** only operator and librarian tools are provided
- **THEN** only operator, librarian, and planner sub-agents SHALL be created

#### Scenario: No tools at all
- **WHEN** the tool list is empty
- **THEN** only the planner sub-agent SHALL be created

#### Scenario: Unmatched-only tools
- **WHEN** all tools have unrecognized prefixes
- **THEN** only the planner sub-agent SHALL be created
- **AND** no unmatched tools SHALL be adapted

### Requirement: Capability-based sub-agent descriptions
Sub-agent descriptions in the orchestrator prompt SHALL use human-readable capability summaries instead of raw tool names. The `capabilityMap` SHALL include entries for all tool prefixes including `secrets_`, `create_skill`, and `list_skills`.

#### Scenario: Operator description uses capabilities
- **WHEN** the operator sub-agent has tools `exec_shell`, `fs_read`
- **THEN** its description SHALL contain "command execution, file operations"
- **AND** it SHALL NOT contain raw tool names

#### Scenario: Vault description uses capabilities
- **WHEN** the vault sub-agent has tools `crypto_sign`, `secrets_get`, `payment_send`
- **THEN** its description SHALL contain "cryptography, secret management, blockchain payments (USDC on Base)"

#### Scenario: Duplicate capabilities are deduplicated
- **WHEN** two tools share the same prefix (e.g., `exec_shell` and `exec_run`)
- **THEN** the capability "command execution" SHALL appear only once in the description

### Requirement: Orchestrator instruction guides delegation-only execution
The orchestrator instruction SHALL enforce mandatory delegation for all tool-requiring tasks. It SHALL include a routing table with exact agent names, a decision protocol, and rejection handling. Sub-agent entries SHALL use capability descriptions, not raw tool name lists.

#### Scenario: Agent name exactness
- **WHEN** the orchestrator delegates to a sub-agent
- **THEN** it SHALL use the EXACT name (e.g. "operator", NOT "exec", "browser", or any abbreviation)

#### Scenario: Invalid agent name prevention
- **WHEN** the orchestrator instruction is generated
- **THEN** it SHALL contain the text "NEVER invent or abbreviate agent names"

#### Scenario: Sub-agent descriptions use capabilities not tool names
- **WHEN** the orchestrator instruction lists sub-agents
- **THEN** each sub-agent entry SHALL describe capabilities
- **AND** SHALL NOT contain raw tool names

### Requirement: RoleToolSet has six roles plus Unmatched
The RoleToolSet struct SHALL have fields: Operator, Navigator, Vault, Librarian, Planner, Chronicler, and Unmatched. Each field is a slice of `*agent.Tool`.

#### Scenario: RoleToolSet structure
- **WHEN** PartitionTools is called
- **THEN** it SHALL return a RoleToolSet with seven fields (six roles + Unmatched)

#### Scenario: Planner tools always empty
- **WHEN** PartitionTools is called with any input
- **THEN** the Planner field SHALL always be nil/empty
