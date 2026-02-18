## Purpose

Per-channel Markdown format conversion for messaging platforms. Converts standard Markdown from LLM output to each platform's native markup dialect at the Send layer.

## Requirements

### Requirement: Telegram Markdown v1 conversion
The system SHALL convert standard Markdown to Telegram Markdown v1 format before sending messages. The conversion SHALL replace `**bold**` with `*bold*`, convert `# Heading` lines to `*Heading*`, and strip `~~strikethrough~~` markers. Inline code and code blocks SHALL be preserved without transformation.

#### Scenario: Bold text conversion
- **WHEN** a message contains `**bold text**`
- **THEN** the system converts it to `*bold text*` for Telegram v1

#### Scenario: Heading conversion
- **WHEN** a message contains `# Heading` or `## Sub Heading`
- **THEN** the system converts them to `*Heading*` and `*Sub Heading*`

#### Scenario: Strikethrough removal
- **WHEN** a message contains `~~struck text~~`
- **THEN** the system removes the `~~` markers, outputting `struck text`

#### Scenario: Code block preservation
- **WHEN** a message contains standard Markdown inside a fenced code block (` ``` `)
- **THEN** the content inside the code block is not transformed

### Requirement: Slack mrkdwn conversion
The system SHALL convert standard Markdown to Slack mrkdwn format before sending messages. The conversion SHALL replace `**bold**` with `*bold*`, `~~strike~~` with `~strike~`, `[text](url)` with `<url|text>`, and `# Heading` with `*Heading*`. Inline code and code blocks SHALL be preserved without transformation.

#### Scenario: Bold text conversion
- **WHEN** a message contains `**bold text**`
- **THEN** the system converts it to `*bold text*` for Slack mrkdwn

#### Scenario: Strikethrough conversion
- **WHEN** a message contains `~~struck text~~`
- **THEN** the system converts it to `~struck text~`

#### Scenario: Link conversion
- **WHEN** a message contains `[click here](https://example.com)`
- **THEN** the system converts it to `<https://example.com|click here>`

#### Scenario: Heading conversion
- **WHEN** a message contains `# Heading`
- **THEN** the system converts it to `*Heading*`

#### Scenario: Code block preservation
- **WHEN** a message contains standard Markdown inside a fenced code block
- **THEN** the content inside the code block is not transformed

### Requirement: Telegram plain text fallback
The system SHALL re-send the original unformatted text when Telegram API rejects a formatted message due to Markdown parse errors.

#### Scenario: API parse error triggers fallback
- **WHEN** a Telegram API call with ParseMode "Markdown" returns an error
- **THEN** the system re-sends the original text without ParseMode

### Requirement: Auto-format on Send
Each channel's `Send()` method SHALL automatically apply format conversion. Callers SHALL NOT need to explicitly format messages.

#### Scenario: Telegram Send auto-formats
- **WHEN** Send() is called with an OutgoingMessage where ParseMode is empty
- **THEN** the system auto-converts text to Telegram v1 and sets ParseMode to "Markdown"

#### Scenario: Slack Send auto-formats
- **WHEN** Send() is called with an OutgoingMessage
- **THEN** the system auto-converts text to Slack mrkdwn before posting

#### Scenario: Explicit ParseMode preserved
- **WHEN** Send() is called with an OutgoingMessage where ParseMode is already set (e.g., "HTML")
- **THEN** the system does not apply Markdown conversion and uses the specified ParseMode
