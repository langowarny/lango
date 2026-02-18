## Why

The `ContextAwareModelAdapter` injects knowledge, memory, and RAG context into the system prompt before each LLM call. However, the `sessionKey` field is hardcoded to an empty string at initialization (`WithMemory(mc.store, "")`), causing the memory retrieval guard (`if m.sessionKey != ""`) to always skip. This means observational memory (observations and reflections) is never injected into the LLM context, despite being fully functional elsewhere in the system.

## What Changes

- Remove the `sessionKey` field from `ContextAwareModelAdapter` struct and the corresponding parameter from `WithMemory`.
- Resolve session key at call time from `context.Context` via `session.SessionKeyFromContext(ctx)` in `GenerateContent`, matching the pattern already used by gateway, channels, and tools.
- Pass the resolved session key as a parameter to `assembleMemorySection`, `assembleRAGSection`, and `assembleGraphRAGSection` instead of reading from the struct field.
- Update `wiring.go` callers to use the simplified `WithMemory(provider)` signature.
- Add unit tests verifying session key extraction from context and memory injection.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `observational-memory`: Session key is now resolved from request context instead of being set at init time, enabling per-request memory retrieval.

## Impact

- `internal/adk/context_model.go`: Struct field removal, method signature changes, new import.
- `internal/app/wiring.go`: Two call sites updated (lines 675, 711).
- `internal/adk/context_model_test.go`: New test file with 4 test cases.
- No breaking changes to external APIs. The `WithMemory` method signature changes but is only called internally.
