## Context

In multi-agent orchestration mode, the ADK orchestrator delegates tasks to sub-agents using `transfer_to_agent`. The orchestrator's system prompt lists sub-agent descriptions that previously included raw tool names like `browser_navigate`, `exec_shell`, etc. The LLM infers fake agent names (e.g., `browser_agent`) from these prefixes, causing `"failed to find agent"` errors from ADK's `base_flow.go`.

## Goals / Non-Goals

**Goals:**
- Eliminate the primary source of agent name hallucination by removing raw tool names from LLM-visible descriptions
- Add a fallback retry mechanism for residual hallucination cases
- Maintain `Tools: nil` on the orchestrator (delegation-only architecture)

**Non-Goals:**
- Modifying ADK internals or patching `base_flow.go`
- Adding tool-call capabilities to the orchestrator
- Handling hallucination for remote A2A agents (they already provide their own descriptions)

## Decisions

### 1. Capability abstraction over tool names

**Decision**: Map tool name prefixes to natural-language capability descriptions via `capabilityMap`.

**Rationale**: The LLM sees `browser_navigate` → infers `browser_agent`. By replacing with "web browsing", there is no agent-like name to hallucinate. This is the primary defense layer.

**Alternative considered**: Stripping tool names entirely — rejected because the orchestrator still needs to understand what each sub-agent can do for correct delegation.

### 2. Single-retry with correction message

**Decision**: On `"failed to find agent: X"` error, send a system correction message listing valid agent names and retry once.

**Rationale**: Even with capability descriptions, LLMs can occasionally hallucinate. A single retry with explicit valid names is low-cost and high-success-rate. More than one retry risks infinite loops with a persistently confused model.

**Alternative considered**: Pre-validation of `transfer_to_agent` targets before execution — rejected because it requires modifying ADK internals.

### 3. Extract-then-retry pattern in RunAndCollect

**Decision**: Refactor `RunAndCollect` into `runAndCollectOnce` (core logic) + wrapper (retry logic) rather than adding retry inside the event loop.

**Rationale**: Clean separation of concerns. The retry wrapper is simple: detect error → build correction → call once more. The core event collection logic remains unchanged.

## Risks / Trade-offs

- **[Risk] Map staleness**: New tools added without updating `capabilityMap` fall back to "general actions" → acceptable degradation, not a failure.
- **[Risk] Retry adds latency**: One extra LLM round-trip on hallucination → acceptable since it only triggers on error (not on success path).
- **[Trade-off] Map iteration order**: `capabilityMap` uses a Go map, so `capabilityDescription` output order is non-deterministic across runs → cosmetic only, no functional impact.
