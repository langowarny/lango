## Context

The multi-agent orchestrator (`BuildAgentTree` in `internal/orchestration/orchestrator.go`) currently receives all tools directly via `llmagent.Config.Tools` in addition to having sub-agents. The orchestrator instruction tells the LLM to handle "simple tasks" directly, but the LLM classifies nearly all tasks as simple and never uses `transfer_to_agent`. This results in sub-agents being effectively dead code, with `message.author` always showing "lango-orchestrator".

ADK's `AgentTransferRequestProcessor` automatically injects a `transfer_to_agent` tool and sub-agent descriptions into LLM requests. However, when the orchestrator already owns all tools, the LLM has no incentive to delegate.

## Goals / Non-Goals

**Goals:**
- Ensure the orchestrator delegates all tool-requiring tasks to sub-agents
- Preserve direct response capability for pure conversational messages (greetings, opinions)
- Maintain single-agent mode (`multiAgent: false`) without any changes

**Non-Goals:**
- Changing tool partitioning logic (already works correctly)
- Modifying sub-agent tool assignments
- Adding new sub-agents or changing existing sub-agent roles

## Decisions

### Decision 1: Remove all direct tools from orchestrator

**Choice**: Set `Tools: nil` on the orchestrator `llmagent.Config`.

**Rationale**: When the orchestrator has no tools, the LLM's only option for tool-requiring tasks is to call `transfer_to_agent` (injected by ADK). This makes delegation mandatory rather than optional.

**Alternative considered**: Keeping a subset of tools on the orchestrator — rejected because any tool presence gives the LLM a reason to skip delegation, and deciding which tools to keep is arbitrary.

### Decision 2: Rewrite instruction to be delegation-first

**Choice**: Replace the "Direct Tool Usage" section with clear delegation rules per sub-agent role.

**Rationale**: The previous instruction explicitly encouraged direct tool usage ("For simple, single-step tasks, call the appropriate tool directly"). The new instruction makes clear the orchestrator has no tools and must delegate.

### Decision 3: Strip tool-related prompt sections from orchestrator

**Choice**: Build a separate prompt for the orchestrator that removes `SectionToolUsage` and replaces `SectionIdentity` with an orchestrator-specific identity that does not mention tool categories.

**Rationale**: The default `SectionIdentity` (from `AGENTS.md`) lists tool categories like "Exec", "Browser", "Crypto" and `SectionToolUsage` (from `TOOL_USAGE.md`) describes individual tools. When injected into the orchestrator's system prompt, the LLM interprets these category names as agent names and attempts `transfer_to_agent("browser")` or `transfer_to_agent("exec")`, causing `failed to find agent` errors. Since the orchestrator has no tools, these sections must be removed.

**Alternative considered**: Renaming tool categories in `AGENTS.md` — rejected because it would affect the single-agent prompt unnecessarily and is a fragile fix.

### Decision 4: Increase MaxDelegationRounds default from 3 to 5

**Choice**: Mandatory delegation adds one extra LLM round-trip per tool-requiring request. Increase default to 5 to avoid premature termination on multi-step tasks.

**Rationale**: orchestrator → sub-agent → tool execution already uses 2 rounds minimum. Complex tasks with multiple delegation steps need headroom.

## Risks / Trade-offs

- **[Increased latency]** Every tool-requiring request now has one additional LLM call (orchestrator decides to delegate, then sub-agent executes). → Acceptable trade-off for correct multi-agent behavior. Pure conversational queries are unaffected.
- **[Token cost]** The extra delegation round-trip consumes additional tokens. → Mitigated by the orchestrator having a focused instruction (no tool descriptions in its prompt since it has no tools).
- **[Sub-agent instruction quality]** Sub-agents need clear instructions to report results effectively back to the orchestrator. → Added result-reporting guidance to each sub-agent's instruction.
