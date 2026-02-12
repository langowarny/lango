
## 1. Setup & Dependencies

- [x] 1.1 Add `google.golang.org/adk` dependency to `go.mod`
- [x] 1.2 Create `internal/adk` package structure
- [x] 1.3 Define `adk.Agent` wrapper struct in `internal/adk/agent.go`

## 2. Core Adapters

- [x] 2.1 Implement `StateStoreAdapter` in `internal/adk/state.go` adapting `internal/session`
- [x] 2.2 Write unit tests for `StateStoreAdapter` ensuring history preservation
- [x] 2.3 Implement `ToolAdapter` in `internal/adk/tools.go` adapting `internal/agent/tools`
- [x] 2.4 Implement `ModelAdapter` factory in `internal/adk/model.go` to configure ADK models from `internal/provider`

## 3. Agent Integration

- [x] 3.1 Implement `NewAgent` factory in `internal/adk/agent.go` constructing the ADK agent
- [x] 3.2 Wire up `internal/app` to use `internal/adk` (replaced supervisor integration)
- [x] 3.3 Verify encryption/decryption of state works via `StateStoreAdapter`

## 4. Verification & Cleanup

- [x] 4.1 Run `lango serve` and verify Telegram bot response with ADK runtime
- [x] 4.2 Verify multi-turn conversation and state persistence
- [x] 4.3 Verify tool execution (via ADK's internal tool handling)
- [x] 4.4 Remove legacy usage in `internal/app` and `internal/gateway`

## 5. Fixes & Tuning

- [x] 5.1 Fix Gemini 404 error by aliasing `gemini` to `gemini-1.5-flash` in provider
- [x] 5.2 Fix "Event from unknown agent" by mapping `Author` in `EventsAdapter`
- [x] 5.3 Fix "High Demand" 503 error by implementing history truncation (last 100 messages)
- [x] 5.4 Fix Gemini 400 "Invalid Argument" error by filtering empty text parts in `EventsAdapter` and `GeminiProvider`
