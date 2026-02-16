## ADDED Requirements

### Requirement: Approval message editing on timeout and cancellation
The Slack approval provider SHALL update the approval message to display "Expired" status and remove action buttons when the approval times out or the context is cancelled.

#### Scenario: Timeout removes buttons
- **WHEN** an approval request times out without user response
- **THEN** the system SHALL call UpdateMessage with text "🔐 Tool approval — ⏱ Expired" and empty MsgOptionBlocks to remove action buttons

#### Scenario: Context cancellation removes buttons
- **WHEN** the approval request context is cancelled
- **THEN** the system SHALL call UpdateMessage with expired text and empty blocks

## MODIFIED Requirements

### Requirement: TOCTOU-safe interactive callback handling
The Slack approval provider SHALL use a single `LoadAndDelete` call as the first operation in HandleInteractive to atomically claim the pending request, preventing the race condition between Load and a concurrent timeout Delete.

#### Scenario: First action click succeeds
- **WHEN** a user clicks an approval button for a pending request
- **THEN** the system SHALL atomically load and delete the pending entry, update the message with approval status, and deliver the result

#### Scenario: Duplicate action click is silently ignored
- **WHEN** a second action arrives for an already-processed request
- **THEN** the system SHALL return immediately without updating the message or delivering a result

### Requirement: Button removal on approval or denial
The Slack approval provider SHALL pass empty `MsgOptionBlocks()` when updating the approval message after user action to remove action buttons.

#### Scenario: Approved message has no buttons
- **WHEN** a user approves a tool request
- **THEN** the updated message SHALL contain no action blocks
