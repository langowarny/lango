## MODIFIED Requirements

### Requirement: Hierarchical agent tree with sub-agents
The system SHALL support a multi-agent mode (`agent.multiAgent: true`) that creates an orchestrator root agent with specialized sub-agents: Executor, Researcher, Planner, and MemoryManager. The orchestrator SHALL have NO direct tools (`Tools: nil`) and MUST delegate all tool-requiring tasks to sub-agents.

#### Scenario: Multi-agent mode enabled
- **WHEN** `agent.multiAgent` is true
- **THEN** BuildAgentTree SHALL create an orchestrator that has NO direct tools AND has sub-agents (Executor, Researcher, Planner, and MemoryManager)

#### Scenario: Orchestrator has no direct tools
- **WHEN** the orchestrator is created
- **THEN** the orchestrator's `Tools` field SHALL be `nil`
- **AND** tools SHALL only be adapted for their respective sub-agents (each tool adapted exactly once)

#### Scenario: Single-agent fallback
- **WHEN** `agent.multiAgent` is false
- **THEN** the system SHALL create a single flat agent with all tools

### Requirement: Orchestrator instruction guides delegation-only execution
The orchestrator instruction SHALL enforce mandatory delegation for all tool-requiring tasks. It SHALL list all valid sub-agent names, explicitly prohibit inventing agent names, and instruct the LLM that it has no tools of its own.

#### Scenario: Tool-requiring task
- **WHEN** a user requests any task requiring tool execution
- **THEN** the orchestrator SHALL delegate to the appropriate sub-agent using its exact registered name

#### Scenario: Delegation rules by sub-agent role
- **WHEN** the orchestrator instruction is generated
- **THEN** it SHALL specify: executor for actions, researcher for information lookup, planner for multi-step planning, memory-manager for memory operations

#### Scenario: Invalid agent name prevention
- **WHEN** the orchestrator instruction is generated
- **THEN** it SHALL contain the text "NEVER invent agent names"
- **AND** it SHALL list only the exact names of registered sub-agents

### Requirement: Orchestrator Short-Circuit for Simple Queries
The orchestrator instruction SHALL direct the LLM to respond directly to simple conversational queries (greetings, opinions, general knowledge) without delegating to sub-agents.

#### Scenario: Simple greeting
- **WHEN** user sends a greeting like "hello"
- **THEN** the orchestrator SHALL respond directly without delegation

#### Scenario: Task requiring tools
- **WHEN** user requests an action requiring tool execution
- **THEN** the orchestrator SHALL delegate to the appropriate sub-agent

### Requirement: Max Delegation Rounds
The `Config` struct SHALL include a `MaxDelegationRounds` field. The orchestrator instruction SHALL mention this limit as a prompt-based guardrail.

#### Scenario: Default max rounds
- **WHEN** `MaxDelegationRounds` is zero or unset
- **THEN** the default limit of 5 rounds SHALL be used in the orchestrator prompt

### Requirement: Sub-Agent Result Reporting
Each sub-agent instruction SHALL include guidance to report results clearly after completing their task.

#### Scenario: Executor result reporting
- **WHEN** the executor sub-agent completes an action
- **THEN** its instruction SHALL guide it to provide results clearly

#### Scenario: Researcher result reporting
- **WHEN** the researcher sub-agent completes research
- **THEN** its instruction SHALL guide it to summarize findings clearly

#### Scenario: Planner result reporting
- **WHEN** the planner sub-agent completes planning
- **THEN** its instruction SHALL guide it to present the plan for review

#### Scenario: Memory Manager result reporting
- **WHEN** the memory-manager sub-agent completes memory operations
- **THEN** its instruction SHALL guide it to report what was stored or retrieved

## REMOVED Requirements

### Requirement: Orchestrator direct tool access
**Reason**: The orchestrator no longer holds direct tools. All tool-requiring tasks must be delegated to sub-agents to ensure proper multi-agent behavior and correct `message.author` tracking.
**Migration**: No migration needed. The orchestrator automatically delegates via ADK's `transfer_to_agent` mechanism.
