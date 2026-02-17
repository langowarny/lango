## Why

LLM agents produce responses in standard Markdown, but Telegram uses Markdown v1 and Slack uses mrkdwn — both with different syntax. Currently, `Send()` transmits raw text, resulting in broken formatting on Telegram (plain text) and Slack (partially broken markup). Discord natively supports standard Markdown and needs no change.

## What Changes

- Add per-channel Markdown format converters that auto-transform standard Markdown at the `Send()` level
- Telegram: standard Markdown → Telegram v1 (`**bold**` → `*bold*`, `# Heading` → `*Heading*`, `~~strike~~` removed), with plain text fallback on API parse errors
- Slack: standard Markdown → mrkdwn (`**bold**` → `*bold*`, `~~strike~~` → `~strike~`, `[text](url)` → `<url|text>`, `# Heading` → `*Heading*`)
- Code blocks (` ``` `) are preserved without transformation on both platforms
- No changes to Discord channel or `internal/app/channels.go` handlers

## Capabilities

### New Capabilities
- `channel-message-format`: Per-channel Markdown format conversion applied automatically at the Send layer

### Modified Capabilities
- `channel-telegram`: Send() now auto-formats standard Markdown to Telegram v1 with plain text fallback
- `channel-slack`: Send() now auto-formats standard Markdown to Slack mrkdwn before sending

## Impact

- `internal/channels/telegram/format.go` — new file with FormatMarkdown()
- `internal/channels/telegram/telegram.go` — Send() modified to auto-format and fallback
- `internal/channels/slack/format.go` — new file with FormatMrkdwn()
- `internal/channels/slack/slack.go` — Send() modified to auto-format
- No dependency changes, no breaking API changes
