## 1. Telegram Markdown Format Converter

- [x] 1.1 Create `internal/channels/telegram/format.go` with `FormatMarkdown()` function that converts standard Markdown to Telegram v1 (bold, heading, strikethrough, code block skip)
- [x] 1.2 Create `internal/channels/telegram/format_test.go` with table-driven tests covering bold, italic, code, heading, strikethrough, compound cases, empty string, and code block preservation

## 2. Slack mrkdwn Format Converter

- [x] 2.1 Create `internal/channels/slack/format.go` with `FormatMrkdwn()` function using compiled regex for bold, strikethrough, link, and heading conversion with code block skip
- [x] 2.2 Create `internal/channels/slack/format_test.go` with table-driven tests covering bold, strike, link, heading, code block preservation, and compound conversions

## 3. Telegram Send Integration

- [x] 3.1 Modify `Send()` in `internal/channels/telegram/telegram.go` to auto-apply `FormatMarkdown()` when ParseMode is empty and set ParseMode to "Markdown"
- [x] 3.2 Add `sendPlainText()` fallback method that re-sends original text without ParseMode on API parse error

## 4. Slack Send Integration

- [x] 4.1 Modify `Send()` in `internal/channels/slack/slack.go` to apply `FormatMrkdwn()` to message text before creating MsgOptionText

## 5. Verification

- [x] 5.1 Run `go test ./internal/channels/telegram/...` — all tests pass
- [x] 5.2 Run `go test ./internal/channels/slack/...` — all tests pass
- [x] 5.3 Run `go build ./...` — full build succeeds
- [x] 5.4 Run `go test ./...` — all project tests pass
