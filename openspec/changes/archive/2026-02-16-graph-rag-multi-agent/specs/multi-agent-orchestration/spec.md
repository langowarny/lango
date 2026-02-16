## ADDED Requirements

### Requirement: Hierarchical agent tree with 4 sub-agents
The system SHALL support a multi-agent mode (`agent.multiAgent: true`) that creates an orchestrator root agent with 4 specialized sub-agents: Executor, Researcher, Planner, and MemoryManager.

#### Scenario: Multi-agent mode enabled
- **WHEN** `agent.multiAgent` is true
- **THEN** BuildAgentTree SHALL create an orchestrator with Executor, Researcher, Planner, and MemoryManager sub-agents

#### Scenario: Single-agent fallback
- **WHEN** `agent.multiAgent` is false
- **THEN** the system SHALL create a single flat agent with all tools

### Requirement: Tool partitioning by prefix
Tools SHALL be partitioned to sub-agents based on name prefixes: `exec/fs_/browser_/crypto_/skill_` → Executor, `search_/rag_/graph_/save_knowledge/save_learning` → Researcher, `memory_/observe_/reflect_` → MemoryManager, unmatched → Executor.

#### Scenario: Graph tools routed to Researcher
- **WHEN** tools named `graph_traverse` and `graph_query` are registered
- **THEN** they SHALL be assigned to the Researcher sub-agent

#### Scenario: Unmatched tools default to Executor
- **WHEN** a tool with an unrecognized prefix is present
- **THEN** it SHALL be assigned to the Executor sub-agent

### Requirement: Graph, RAG, and Memory agent tools
The system SHALL provide dedicated tools for sub-agents: `graph_traverse`, `graph_query` (graph store), `rag_retrieve` (RAG service), `memory_list_observations`, `memory_list_reflections` (memory store).

#### Scenario: Graph tools available when graph enabled
- **WHEN** `graph.enabled: true`
- **THEN** `graph_traverse` and `graph_query` tools SHALL be added to the tool set

#### Scenario: RAG tool available when embedding configured
- **WHEN** embedding provider is configured and RAG service is initialized
- **THEN** `rag_retrieve` tool SHALL be added to the tool set

#### Scenario: Memory tools available when observational memory enabled
- **WHEN** `observationalMemory.enabled: true`
- **THEN** `memory_list_observations` and `memory_list_reflections` tools SHALL be added

### Requirement: Remote agents as sub-agents
The orchestrator SHALL accept remote A2A agents and append them to its sub-agent list. Remote agent names and descriptions SHALL be included in the orchestrator instruction.

#### Scenario: Remote agents loaded and wired
- **WHEN** `a2a.enabled: true` and `a2a.remoteAgents` contains entries
- **THEN** LoadRemoteAgents SHALL create ADK agents and they SHALL appear as sub-agents in the orchestrator

#### Scenario: Remote agent load failure
- **WHEN** a remote agent card URL is unreachable
- **THEN** the agent SHALL be skipped with a warning log, and the orchestrator SHALL continue with local sub-agents
