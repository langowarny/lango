## ADDED Requirements

### Requirement: Approval message editing on timeout and cancellation
The Discord approval provider SHALL edit the approval message to display "Expired" status and remove button components when the approval times out or the context is cancelled.

#### Scenario: Timeout removes buttons
- **WHEN** an approval request times out without user response
- **THEN** the system SHALL call ChannelMessageEditComplex with content "🔐 Tool approval — ⏱ Expired" and empty Components slice

#### Scenario: Context cancellation removes buttons
- **WHEN** the approval request context is cancelled
- **THEN** the system SHALL call ChannelMessageEditComplex with expired content and empty Components

### Requirement: Message ID capture for timeout editing
The Discord approval provider SHALL capture the message ID from ChannelMessageSendComplex return value and store it in an approvalPending struct for use in timeout/cancellation message editing.

#### Scenario: Sent message ID stored in pending struct
- **WHEN** an approval message is sent successfully
- **THEN** the system SHALL store the returned message ID alongside the response channel and channel ID in an approvalPending struct
