## Context

The orchestration layer (`internal/orchestration/`) builds a hierarchical agent tree with an orchestrator root delegating to sub-agents. The current 4-agent design (executor/researcher/planner/memory-manager) has executor acting as a catch-all for 6+ domains, no structured routing, and no mechanism for sub-agents to reject misrouted tasks.

## Goals / Non-Goals

**Goals:**
- Split executor's 6+ domains into focused agents (operator, navigator, vault)
- Rename researcher → librarian, memory-manager → chronicler for clarity
- Add AgentSpec type for data-driven agent creation
- Add structured routing table with keywords, I/O specs, and negative constraints
- Add decision protocol and reject mechanism
- Track unmatched tools separately

**Non-Goals:**
- Runtime dynamic re-routing (static routing table is sufficient)
- Changes to Config type signature or BuildAgentTree public API
- Changes to wiring.go or other callers
- Agent-to-agent direct communication

## Decisions

### 1. Data-Driven Agent Creation via AgentSpec Registry

**Decision**: Define a `var agentSpecs = []AgentSpec{...}` registry and iterate it in BuildAgentTree instead of hardcoded blocks per agent.

**Rationale**: Adding/removing/modifying agents becomes a data change in one place rather than scattered code blocks. The loop ensures consistent construction logic.

**Alternative**: Keep per-agent blocks → rejected because adding new agents requires duplicating boilerplate and is error-prone.

### 2. Six-Agent Topology (operator/navigator/vault/librarian/planner/chronicler)

**Decision**: Split executor into operator (shell/file/skill), navigator (browser), vault (crypto/secrets/payment). Rename researcher → librarian, memory-manager → chronicler.

**Rationale**: Each agent has a coherent, narrow responsibility. Naming reflects actual function (librarian manages knowledge, chronicler records memory).

**Alternative**: Keep 4 agents with better prompts → rejected because executor's scope is too broad for reliable prompt-based routing.

### 3. Prefix Matching Order: Librarian First

**Decision**: Match prefixes in order: Librarian → Chronicler → Navigator → Vault → Operator → Unmatched.

**Rationale**: Librarian has exact-match prefixes (`save_knowledge`, `create_skill`, `list_skills`) that would otherwise match operator's broader prefixes. Checking librarian first prevents misrouting.

### 4. Unmatched Tools as Separate Bucket

**Decision**: Tools matching no prefix go to `Unmatched` instead of defaulting to executor/operator.

**Rationale**: Explicit tracking prevents silent misrouting. The orchestrator prompt mentions unmatched tools so the LLM can handle them.

### 5. Routing Table in Orchestrator Prompt

**Decision**: Include a structured routing table with keywords, accepts/returns, and cannot-do per agent in the orchestrator instruction.

**Rationale**: Gives the LLM explicit decision criteria instead of relying on description matching alone.

### 6. Reject Protocol

**Decision**: Each sub-agent's instruction includes a `[REJECT]` response format for misrouted tasks.

**Rationale**: Allows graceful recovery when routing fails. The orchestrator instruction includes handling for rejections.

## Risks / Trade-offs

- **Longer orchestrator prompt** → Mitigated by structured format; routing table adds ~30 lines but significantly improves routing accuracy
- **More sub-agents = more delegation overhead** → Mitigated by conditional creation (agents without tools are skipped)
- **Agent name changes break existing conversations** → No persistence of agent names in user-facing state; ADK reads names dynamically
