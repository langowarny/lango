## Why

The current multi-agent orchestration system has structural problems: Executor acts as a catch-all (handling 6+ domains), Researcher's name misleads (it does internal knowledge lookup, not external research), sub-agent prompts lack I/O specs and negative constraints, there's no routing table (relying solely on LLM judgment), and sub-agents cannot reject misrouted tasks.

## What Changes

- Split Executor (6+ domains) into 3 focused agents: **operator** (shell, file, skill), **navigator** (web browsing), **vault** (crypto, secrets, payments)
- Rename Researcher → **librarian** with explicit knowledge/skill management scope
- Rename Memory Manager → **chronicler** for clarity
- Add `AgentSpec` type to define each agent's identity, routing metadata, I/O spec, and negative constraints
- Replace hardcoded agent creation blocks with data-driven loop over `agentSpecs` registry
- Add structured routing table with keywords, accepts/returns specs, and cannot-do lists to orchestrator prompt
- Add decision protocol (CLASSIFY → MATCH → SELECT → VERIFY → DELEGATE) to orchestrator
- Add `[REJECT]` protocol for sub-agents to refuse misrouted tasks
- Track unmatched tools separately instead of dumping them into Executor

## Capabilities

### New Capabilities

- `agent-routing`: Keyword-based routing table and decision protocol for orchestrator delegation

### Modified Capabilities

- `multi-agent-orchestration`: Sub-agent structure changes from 4 roles (executor/researcher/planner/memory-manager) to 6 roles (operator/navigator/vault/librarian/planner/chronicler) with AgentSpec-driven creation, I/O-specified prompts, reject protocol, and unmatched tool tracking

## Impact

- `internal/orchestration/tools.go` — Major rewrite: AgentSpec type, 6-role RoleToolSet, new prefix mapping, Unmatched bucket
- `internal/orchestration/orchestrator.go` — Major rewrite: data-driven agent creation loop, routing table builder, orchestrator prompt assembler
- `internal/orchestration/orchestrator_test.go` — Major rewrite: 25+ tests covering 6-role partitioning, routing, reject protocol, spec consistency
- No changes to `Config` type signature or `BuildAgentTree` public API — external callers unaffected
- No changes to `internal/app/wiring.go` — it calls `BuildAgentTree(cfg)` which remains the same
