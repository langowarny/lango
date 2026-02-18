## MODIFIED Requirements

### Requirement: Message sending
The Telegram channel SHALL auto-convert standard Markdown to Telegram v1 format in the Send() method when ParseMode is not explicitly set. On API parse failure, the system SHALL fallback to sending plain text without ParseMode. Messages with explicit ParseMode SHALL be sent as-is.

#### Scenario: Auto-format standard Markdown
- **WHEN** Send() is called with an OutgoingMessage where ParseMode is empty
- **THEN** the system converts the text via FormatMarkdown() and sends with ParseMode "Markdown"

#### Scenario: Plain text fallback on parse error
- **WHEN** the Telegram API rejects a message due to Markdown parse error
- **THEN** the system re-sends the original unformatted text without ParseMode

#### Scenario: Explicit ParseMode preserved
- **WHEN** Send() is called with ParseMode "HTML"
- **THEN** the system sends the text with ParseMode "HTML" without Markdown conversion

#### Scenario: Long message chunking preserved
- **WHEN** a formatted message exceeds 4096 characters
- **THEN** the system splits the message into chunks and sends each with Markdown formatting
