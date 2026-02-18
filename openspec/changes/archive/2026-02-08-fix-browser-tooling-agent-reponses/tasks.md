## 1. Agent Runtime Updates

- [x] 1.1 Add `SystemPrompt` field to `agent.Config` in `internal/agent/runtime.go`.
- [x] 1.2 Update `Runtime.Run` in `internal/agent/runtime.go` to prepend the system prompt if available when starting a new session.

## 2. Browser Tooling Integration

- [x] 2.1 Update `browser_navigate` tool handler in `internal/app/app.go` to return page title and content snippet.
- [x] 2.2 Register `browser_read` tool handler in `internal/app/app.go`.
- [x] 2.3 Register `browser_screenshot` tool handler in `internal/app/app.go`.

## 3. Configuration and Prompting

- [x] 3.1 Define a default system prompt constant in `internal/app/app.go`.
- [x] 3.2 Initialize `agent.Config` with the system prompt in `app.New`.

## 4. Verification

- [x] 4.1 Verify that `lango serve` correctly registers all browser tools.
- [x] 4.2 Test browser navigation and reading with a real URL (e.g., example.com).
- [x] 4.3 Confirm the agent correctly identifies its tool capabilities via the system prompt.
