## 1. Model Adapter Stream Handling

- [x] 1.1 Add non-streaming path to `ModelAdapter.GenerateContent` that accumulates all provider events and yields a single complete `LLMResponse` with `Partial=false` and `TurnComplete=true`
- [x] 1.2 Update streaming path to accumulate text and include full accumulated text in the final done event's `Content.Parts`
- [x] 1.3 Handle `StreamEventError` in both streaming and non-streaming paths

## 2. Agent Text Collection Fix

- [x] 2.1 Update `runAndCollectOnce` to track whether partial events were seen and avoid double-counting text from the final non-partial event

## 3. Tests

- [x] 3.1 Update `TestModelAdapter_GenerateContent_ToolCall` to expect a single accumulated response in non-streaming mode
- [x] 3.2 Verify all existing ADK tests pass with the new behavior
- [x] 3.3 Run full project test suite (`go test ./...`) to confirm no regressions
