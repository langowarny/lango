
## Context
The current agent runtime (`internal/agent`) is a custom implementation that manually handles LLM API interactions, tool execution, and session management. This has led to fragility, especially with complex providers like Gemini. We are migrating to `google/adk-go` to leverage a supported, production-ready framework.

## Goals / Non-Goals

**Goals:**
- Replace `internal/agent` with an ADK-based implementation.
- Maintain feature parity (multi-provider support, tools, session persistence).
- Fix existing bugs related to protocol handling (e.g., Gemini empty text parts).
- Simplify the codebase by removing custom plumbing for tool execution and message formatting.

**Non-Goals:**
- Changing the underlying session storage (Ent/SQLite) schema (except where strictly necessary for ADK compatibility).
- Rewriting the frontend or gateway API (API contract should remain stable).

## Decisions

### 1. Agent Architecture
We will create a new package `internal/adk` (or repurpose `internal/agent`) that wraps `adk.Agent`.
- **Why**: ADK's `agent.Agent` is the core struct. We need to wrap it to integrate with our `internal/config` and `internal/session`.
- **Alternative**: Using ADK directly in `internal/app`. **Rejected** because we need a translation layer for our existing session models and configuration.

### 2. State Management
We will implement ADK's `state.Store` interface using our `internal/session` Ent store.
- **Why**: ADK uses a simple key-value or object store interface. Our Ent store is relational. We need an adapter to map ADK's state requirements to our DB schema.
- **Alternative**: Use ADK's default file/memory store. **Rejected** because we need persistence and compatibility with existing user sessions.

### 3. Tool Adaptation
We will write a generic adapter to convert our functional tools (referenced in `runtime.go`) to ADK's `tools.Tool` interface.
- **Why**: ADK tools expect specific signatures. Adapting existing tools is faster than rewriting them all.

### 4. Provider Strategy
We will use ADK's model abstraction.
- **Why**: ADK provides `model.Model` interface. We will configure ADK with the appropriate model implementation based on our config (Gemini, OpenAI, etc.). ADK likely has built-in support for Gemini and standard OpenAI-compatible providers.

## Risks / Trade-offs

- [x] **Risk**: ADK might have different assumptions about session history formatting.
  - **Mitigation**: Thorough testing of the `StateAdapter` and implementing explicit `Author` mapping.
- [x] **Risk**: "Ejecting" from custom runtime might lose some granular control.
  - **Mitigation**: ADK is designed to be extensible. We will verify we can still inject our middleware (PII redaction, approval).
- [x] **Risk**: Large session history causing "High Demand" errors or context overflow.
  - **Mitigation**: Implemented history truncation (last 100 messages) in `EventsAdapter`.

## Migration Plan
1.  **Refactor**: Create `internal/adk` and implement the `state.Store` adapter.
2.  **Integrate**: Update `internal/supervisor` to initialize the ADK agent instead of the legacy runtime.
3.  **Verify**: Run end-to-end tests with `lango serve`.
4.  **Cleanup**: Remove the old `internal/agent` code once verified.

## Open Questions
- Does ADK Support "Approval" middleware out of the box, or do we need to wrap the tool execution? (Likely need to wrap/interceptor). Answer: Yes, ADK has an ApprovalMiddleware that we can use.
