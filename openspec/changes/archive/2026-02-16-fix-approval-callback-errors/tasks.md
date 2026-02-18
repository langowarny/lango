## 1. Telegram Approval Handler

- [x] 1.1 Add `approvalPending` struct with `ch`, `chatID`, `messageID` fields
- [x] 1.2 Implement `editApprovalMessage` helper with empty keyboard markup
- [x] 1.3 Add `isCallbackExpiredErr` and `isMessageNotModifiedErr` error classifiers
- [x] 1.4 Refactor `RequestApproval` to capture `sentMsg.MessageID`, store `*approvalPending`, edit message on timeout/cancel
- [x] 1.5 Refactor `HandleCallback` to use `LoadAndDelete`-first pattern, answer expired callbacks at Debug level
- [x] 1.6 Add tests: keyboard removal, timeout editing, context cancellation editing, duplicate callback

## 2. Slack Approval Handler

- [x] 2.1 Fix TOCTOU race: replace `Load` + `LoadAndDelete` with single `LoadAndDelete` in `HandleInteractive`
- [x] 2.2 Add `editExpiredMessage` helper using `UpdateMessage` with empty `MsgOptionBlocks()`
- [x] 2.3 Add timeout/cancel message editing in `RequestApproval` select cases
- [x] 2.4 Pass empty `MsgOptionBlocks()` in `HandleInteractive` message update to remove buttons
- [x] 2.5 Add tests: TOCTOU fix (duplicate action), timeout message editing, button removal verification

## 3. Discord Approval Handler

- [x] 3.1 Add `approvalPending` struct with `ch`, `channelID`, `messageID` fields
- [x] 3.2 Capture message ID from `ChannelMessageSendComplex` return value
- [x] 3.3 Implement `editExpiredMessage` helper using `ChannelMessageEditComplex`
- [x] 3.4 Add timeout/cancel message editing in `RequestApproval` select cases
- [x] 3.5 Add tests: timeout message editing, context cancellation editing

## 4. Gateway Approval Handler

- [x] 4.1 Add `ApprovalTimeout time.Duration` field to `Config` struct
- [x] 4.2 Replace hardcoded `30 * time.Second` with config timeout (with 30s default fallback)
- [x] 4.3 Fix `handleApprovalResponse` to delete entry inside lock scope (atomic delete)
- [x] 4.4 Add tests: atomic delete verification, duplicate response prevention, config timeout usage

## 5. Verification

- [x] 5.1 Run `go build ./...` to verify compilation
- [x] 5.2 Run `go test ./internal/channels/telegram/... -v -race` — all approval tests pass
- [x] 5.3 Run `go test ./internal/channels/slack/... -v -race` — all approval tests pass
- [x] 5.4 Run `go test ./internal/channels/discord/... -v -race` — all tests pass
- [x] 5.5 Run `go test ./internal/gateway/... -v -race` — all tests pass
