## ADDED Requirements

### Requirement: Capability-based sub-agent descriptions
Sub-agent descriptions in the orchestrator prompt SHALL use human-readable capability summaries instead of raw tool names. A `capabilityMap` SHALL map tool name prefixes to natural-language descriptions (e.g., `browser_` → "web browsing", `exec` → "command execution"). The `capabilityDescription()` function SHALL deduplicate capabilities across a tool set.

#### Scenario: Executor description uses capabilities
- **WHEN** the executor sub-agent has tools `exec_shell`, `fs_read`, `browser_navigate`
- **THEN** its description SHALL contain "command execution, file operations, web browsing"
- **AND** it SHALL NOT contain raw tool names like "exec_shell" or "browser_navigate"

#### Scenario: Duplicate capabilities are deduplicated
- **WHEN** two tools share the same prefix (e.g., `exec_shell` and `exec_run`)
- **THEN** the capability "command execution" SHALL appear only once in the description

#### Scenario: Unknown tool prefix falls back to general actions
- **WHEN** a tool has no matching prefix in `capabilityMap`
- **THEN** its capability SHALL be "general actions"

## MODIFIED Requirements

### Requirement: Orchestrator instruction guides delegation-only execution
The orchestrator instruction SHALL enforce mandatory delegation for all tool-requiring tasks. It SHALL list all valid sub-agent names with exact spelling, explicitly prohibit inventing or abbreviating agent names, and instruct the LLM that it has no tools of its own. The instruction SHALL NOT contain words that could be confused with agent names (e.g., avoid listing action types like "browser" that might be mistaken for an agent name). Sub-agent entries in the instruction SHALL use capability descriptions, not raw tool name lists.

#### Scenario: Tool-requiring task
- **WHEN** a user requests any task requiring tool execution
- **THEN** the orchestrator SHALL delegate to the appropriate sub-agent using its exact registered name

#### Scenario: Agent name exactness
- **WHEN** the orchestrator delegates to a sub-agent
- **THEN** it SHALL use the EXACT name (e.g. "executor", NOT "exec", "browser", or any abbreviation)

#### Scenario: Delegation rules by sub-agent role
- **WHEN** the orchestrator instruction is generated
- **THEN** it SHALL specify: executor for actions, researcher for information lookup, planner for multi-step planning, memory-manager for memory operations

#### Scenario: Invalid agent name prevention
- **WHEN** the orchestrator instruction is generated
- **THEN** it SHALL contain the text "NEVER invent or abbreviate agent names"
- **AND** it SHALL list only the exact names of registered sub-agents

#### Scenario: Sub-agent descriptions use capabilities not tool names
- **WHEN** the orchestrator instruction lists sub-agents
- **THEN** each sub-agent entry SHALL describe capabilities (e.g., "command execution, file operations")
- **AND** SHALL NOT contain raw tool names (e.g., "exec_shell", "browser_navigate")
