## ADDED Requirements

### Requirement: Approval message editing on timeout and cancellation
The Telegram approval provider SHALL edit the approval message to display "Expired" status and remove inline keyboard buttons when the approval times out or the context is cancelled.

#### Scenario: Timeout removes buttons
- **WHEN** an approval request times out without user response
- **THEN** the system SHALL edit the original message to "🔐 Tool approval — ⏱ Expired" with an empty inline keyboard markup

#### Scenario: Context cancellation removes buttons
- **WHEN** the approval request context is cancelled
- **THEN** the system SHALL edit the original message to "🔐 Tool approval — ⏱ Expired" with an empty inline keyboard markup

### Requirement: Duplicate callback prevention via LoadAndDelete
The Telegram approval provider SHALL use `LoadAndDelete` as the first operation when handling callbacks to atomically claim the pending request and prevent duplicate processing.

#### Scenario: First callback succeeds
- **WHEN** a user clicks an approval button for a pending request
- **THEN** the system SHALL atomically load and delete the pending entry, deliver the result, and edit the message

#### Scenario: Duplicate callback is silently ignored
- **WHEN** a second callback arrives for an already-processed request
- **THEN** the system SHALL answer the callback silently without delivering a duplicate result

### Requirement: Error classification for callback and message operations
The Telegram approval provider SHALL classify API errors to suppress benign conditions at appropriate log levels.

#### Scenario: Expired callback query logged at debug level
- **WHEN** answering a callback fails with "query is too old" error
- **THEN** the system SHALL log at Debug level, not Warn level

#### Scenario: Message not modified error suppressed
- **WHEN** editing a message fails with "message is not modified" error
- **THEN** the system SHALL suppress the error without logging
