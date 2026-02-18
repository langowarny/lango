## MODIFIED Requirements

### Requirement: Message sending
The Slack channel SHALL auto-convert standard Markdown to Slack mrkdwn format in the Send() method before posting messages. The conversion SHALL apply to the text content passed via MsgOptionText.

#### Scenario: Auto-format standard Markdown
- **WHEN** Send() is called with an OutgoingMessage
- **THEN** the system converts the text via FormatMrkdwn() before creating MsgOptionText

#### Scenario: Thread reply preserved
- **WHEN** Send() is called with ThreadTS set
- **THEN** the system applies mrkdwn conversion and replies in the specified thread

#### Scenario: Block Kit content unaffected
- **WHEN** Send() is called with Blocks specified
- **THEN** the system converts the text field but Block Kit content is sent as-is
