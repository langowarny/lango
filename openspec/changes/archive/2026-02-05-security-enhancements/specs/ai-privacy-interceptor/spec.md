## ADDED Requirements

### Requirement: Intercept AI Requests
The system SHALL intercept all agent requests to the LLM and pass them through a filter chain before execution.

#### Scenario: Interception
- **WHEN** an agent initiates a chat completion request
- **THEN** the request parameters (messages, tools) are inspected by the middleware

### Requirement: PII Redaction
The system SHALL automatically redact personally identifiable information (email, phone numbers, API keys) from the prompt content using regex patterns.

#### Scenario: Redact Email
- **WHEN** a user prompt contains "my email is test@example.com"
- **THEN** it is rewritten to "my email is [REDACTED]" before sending to the LLM

### Requirement: Approval Workflow
The system SHALL block execution of "sensitive" tools (configured list) and request approval via a configured notification channel (Discord/Telegram).

#### Scenario: Sensitive Action
- **WHEN** the agent attempts to call `exec` with "rm -rf /"
- **THEN** the execution is paused
- **AND** a notification is sent to the admin channel with an Approve/Deny link
- **AND** the system waits for the callback before proceeding or aborting
