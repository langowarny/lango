## Why

The AI repeats previous answers when responding to new questions. When a user asks Q1 → A1 then Q2, the AI responds with A1 + A2 instead of just A2. This happens because `ModelAdapter.GenerateContent` ignores the `stream` parameter and always streams, causing empty assistant messages to be stored in the session history.

## What Changes

- Fix `ModelAdapter.GenerateContent` in `internal/adk/model.go` to respect the `stream` parameter:
  - **Non-streaming mode** (`stream=false`): Accumulate all provider streaming events internally and yield a single non-partial `LLMResponse` with the complete text, so the ADK runner stores the full assistant message in the session.
  - **Streaming mode** (`stream=true`): Continue yielding partial text events for real-time UI, but include accumulated full text in the final done event for proper session storage.
- Fix `runAndCollectOnce` in `internal/adk/agent.go` to avoid double-counting text from both partial and non-partial events.
- Update existing tests to match the new non-streaming behavior (single accumulated response instead of multiple events).

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `adk-architecture`: The model adapter now correctly handles the `stream` parameter, changing how LLMResponse events are yielded to the ADK runner.

## Impact

- **Code**: `internal/adk/model.go`, `internal/adk/agent.go`, `internal/adk/model_test.go`
- **Behavior**: Session history will now contain complete assistant messages instead of empty ones, preventing the LLM from re-answering previous questions.
- **Dependencies**: No new dependencies. Uses existing `strings.Builder` for text accumulation.
- **APIs**: No public API changes. Internal `ModelAdapter` behavior is corrected.
