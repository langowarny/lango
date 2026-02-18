## Context

The `ContextAwareModelAdapter` is the central component that augments each LLM call with knowledge, memory, and RAG context. It was initialized with a `sessionKey` field via `WithMemory(provider, sessionKey)`, but `wiring.go` always passed an empty string because the actual session key is not available at init time — it varies per request.

The session key is already propagated through `context.Context` by the gateway (`server.go`) and channel adapters (`channels.go`) using `session.WithSessionKey(ctx, key)`. Multiple tools and services already extract it via `session.SessionKeyFromContext(ctx)`. The `ContextAwareModelAdapter` was the only component not using this pattern.

## Goals / Non-Goals

**Goals:**
- Enable observational memory injection into LLM context by resolving session key at call time from `context.Context`.
- Align `ContextAwareModelAdapter` with the existing context-based session key pattern used throughout the codebase.
- Ensure RAG and GraphRAG session key filtering works correctly per-request.
- Ensure `RuntimeContextAdapter` receives the correct session key per-request.

**Non-Goals:**
- Changing how session keys are set in context (gateway/channels already handle this correctly).
- Modifying the `MemoryProvider` interface or memory storage layer.
- Adding new memory retrieval capabilities beyond fixing the existing broken path.

## Decisions

### Decision 1: Remove `sessionKey` field entirely (vs. keeping as fallback)

**Choice**: Remove the field and always resolve from context.

**Rationale**: The field was always empty in practice. Keeping it as a fallback adds complexity with no benefit — all callers already set session key in context. A single resolution path is simpler and less error-prone.

**Alternative considered**: Keep the field as a fallback (`if sk := SessionKeyFromContext(ctx); sk != "" { use sk } else { use m.sessionKey }`). Rejected because the fallback value was always empty, adding dead code.

### Decision 2: Pass `sessionKey` as parameter to assemble methods (vs. storing in a request-scoped field)

**Choice**: Pass as a function parameter.

**Rationale**: The session key is request-scoped, not adapter-scoped. Passing it as a parameter makes the data flow explicit and avoids any concurrency concerns with a mutable field on a shared adapter instance.

## Risks / Trade-offs

- **[Low] Context missing session key**: If a code path calls `GenerateContent` without setting session key in context, memory retrieval will be skipped (same as current behavior). → Mitigation: This matches the existing guard condition and is the correct behavior for contexts without a session.
- **[Low] Internal API change**: `WithMemory` signature changes from 2 params to 1. → Mitigation: Only called internally in `wiring.go` (2 call sites). No external consumers.
