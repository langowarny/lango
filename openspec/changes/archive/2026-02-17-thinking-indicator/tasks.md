## 1. Discord Typing Indicator

- [x] 1.1 Add `ChannelTyping` method to Discord `Session` interface in `session.go`
- [x] 1.2 Add `startTyping` method to Discord `Channel` with goroutine + 8s ticker
- [x] 1.3 Modify `onMessageCreate` to call `startTyping` before handler and stop after
- [x] 1.4 Add `ChannelTyping` to `MockSession` and add typing indicator test

## 2. Telegram Typing Indicator

- [x] 2.1 Add `startTyping` method to Telegram `Channel` using `Request(ChatActionConfig)` with 4s ticker
- [x] 2.2 Modify `handleUpdate` to call `startTyping` before handler and stop after
- [x] 2.3 Add `RequestCalls` tracking to `MockBotAPI` and add typing indicator test

## 3. Slack Thinking Placeholder

- [x] 3.1 Add `postThinking` method to post "_Thinking..._" placeholder message
- [x] 3.2 Add `updateThinking` method to replace placeholder via `UpdateMessage`
- [x] 3.3 Modify `handleMessage` goroutine to use placeholder flow with fallbacks
- [x] 3.4 Add `UpdateMessages` tracking to `MockClient` and add placeholder flow test

## 4. Gateway Session-Scoped Broadcast

- [x] 4.1 Add `BroadcastToSession` method with session-key scoping for UI clients
- [x] 4.2 Modify `handleChatMessage` to emit `agent.thinking` before and `agent.done` after processing
- [x] 4.3 Add session-scoped broadcast tests (authenticated + unauthenticated)

## 5. Verification

- [x] 5.1 Run `go build ./...` to verify clean compilation
- [x] 5.2 Run `go test ./internal/channels/... ./internal/gateway/...` to verify all tests pass
