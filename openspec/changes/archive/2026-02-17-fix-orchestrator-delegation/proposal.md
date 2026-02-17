## Why

The multi-agent orchestrator holds all tools directly and its instruction tells the LLM to handle "simple tasks" itself. Because the LLM classifies most requests as simple, it never delegates to sub-agents — making executor, researcher, planner, and memory-manager effectively dead code. The `message.author` field always shows "lango-orchestrator" regardless of task complexity.

## What Changes

- Remove all direct tools from the orchestrator agent so it cannot call tools itself
- Strip tool-related prompt sections (`SectionIdentity`, `SectionToolUsage`) from orchestrator system prompt to prevent LLM from hallucinating agent names like "browser" or "exec"
- Rewrite orchestrator instruction to enforce delegation for any tool-requiring task
- Improve sub-agent instructions with result-reporting guidance
- Increase default `MaxDelegationRounds` from 3 to 5 to accommodate the extra round-trip required by mandatory delegation
- Update tests to verify orchestrator has no direct tools

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `multi-agent-orchestration`: Orchestrator no longer holds direct tools; all tool-requiring tasks must be delegated to sub-agents

## Impact

- **Code**: `internal/orchestration/orchestrator.go` (core change), `internal/app/wiring.go` (config), `internal/orchestration/orchestrator_test.go` (tests)
- **Behavior**: Tool-requiring tasks now always incur one additional LLM round-trip (orchestrator → sub-agent). Pure conversational messages are unaffected.
- **Single-agent mode**: Completely separate code path — no impact
