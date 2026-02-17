## Why

In multi-agent mode, the orchestrator LLM hallucinates non-existent agent names like `browser_agent` because raw tool name prefixes (e.g., `browser_navigate`, `exec_shell`) are exposed in sub-agent descriptions. This causes `"failed to find agent: browser_agent"` errors from the ADK runner.

## What Changes

- Replace raw tool name lists in sub-agent descriptions with human-readable capability descriptions (e.g., `browser_navigate` → "web browsing"), removing the source material the LLM uses to fabricate agent names.
- Add a fallback retry layer in `RunAndCollect`: when a "failed to find agent" error is detected, extract the hallucinated name, send a correction message with valid agent names, and retry once.
- Add `toolCapability()` and `capabilityDescription()` helper functions that map tool name prefixes to natural-language capability strings.

## Capabilities

### New Capabilities

_(none — this is a bugfix within existing capabilities)_

### Modified Capabilities

- `multi-agent-orchestration`: Sub-agent descriptions now use capability summaries instead of raw tool names; orchestrator prompt unchanged.
- `adk-architecture`: `RunAndCollect` gains a 1-retry fallback on agent-not-found errors with correction messaging.

## Impact

- `internal/orchestration/tools.go` — new `capabilityMap`, `toolCapability()`, `capabilityDescription()` functions
- `internal/orchestration/orchestrator.go` — agent description strings changed from tool names to capabilities
- `internal/adk/agent.go` — `RunAndCollect` refactored into `runAndCollectOnce` + retry wrapper; new `extractMissingAgent`, `subAgentNames` helpers
- `internal/orchestration/orchestrator_test.go` — new tests for capability mapping
- `internal/adk/agent_test.go` — new tests for error extraction
