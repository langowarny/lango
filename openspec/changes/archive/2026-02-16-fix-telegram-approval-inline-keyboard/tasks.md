## 1. Fix Inline Keyboard Serialization

- [x] 1.1 Replace `tgbotapi.NewInlineKeyboardMarkup()` with struct literal using empty `InlineKeyboard` slice in `editApprovalMessage` (`internal/channels/telegram/approval.go`)

## 2. Verification

- [x] 2.1 Run existing approval tests to confirm no regressions (`go test ./internal/channels/telegram/... -run TestApproval`)
- [x] 2.2 Verify `go build ./...` succeeds with no errors
