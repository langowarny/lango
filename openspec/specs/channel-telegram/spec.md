## ADDED Requirements

### Requirement: Telegram bot connection
The system SHALL connect to Telegram using the Bot API with a provided bot token.

#### Scenario: Successful bot connection
- **WHEN** the application starts with a valid TELEGRAM_BOT_TOKEN
- **THEN** the system SHALL establish a connection to Telegram servers

#### Scenario: Invalid bot token
- **WHEN** the bot token is invalid or revoked
- **THEN** the system SHALL log an error and retry with exponential backoff

### Requirement: Message reception
The system SHALL receive and process incoming messages from Telegram chats. Message handling SHALL be dispatched to a separate goroutine so that the event loop remains non-blocking and can continue processing CallbackQuery updates concurrently.

#### Scenario: Direct message received
- **WHEN** a user sends a direct message to the bot
- **THEN** the message SHALL be forwarded to the agent with sender context

#### Scenario: Group message with mention
- **WHEN** a user mentions the bot in a group chat
- **THEN** the message SHALL be processed if mention-gating allows

#### Scenario: Concurrent callback processing during handler execution
- **WHEN** a message handler is blocking (e.g., waiting for tool approval)
- **AND** a CallbackQuery update arrives from a button click
- **THEN** the event loop SHALL process the CallbackQuery immediately without waiting for the handler to complete

#### Scenario: Graceful shutdown with active handlers
- **WHEN** the channel is stopped while message handlers are still executing
- **THEN** the system SHALL wait for all active handler goroutines to complete before returning from Stop()

### Requirement: Message sending
The system SHALL send agent responses back to Telegram chats.

#### Scenario: Send text response
- **WHEN** the agent generates a text response
- **THEN** the response SHALL be sent to the originating chat

#### Scenario: Long message chunking
- **WHEN** a response exceeds Telegram's 4096 character limit
- **THEN** the response SHALL be split into multiple messages

### Requirement: Media handling
The system SHALL process media attachments from Telegram messages.

#### Scenario: Image attachment
- **WHEN** a user sends an image
- **THEN** the system SHALL download the image and provide it to the agent

#### Scenario: Document attachment
- **WHEN** a user sends a document
- **THEN** the system SHALL download and make it available for processing

### Requirement: Allowlist filtering
The system SHALL filter incoming messages based on configured allowlists.

#### Scenario: Allowed sender
- **WHEN** a message is from an allowed user/group
- **THEN** the message SHALL be processed normally

#### Scenario: Unknown sender
- **WHEN** a message is from an unknown sender and pairing is enabled
- **THEN** a pairing code SHALL be sent to the sender

### Requirement: Telegram approval provider
The Telegram channel SHALL provide an approval provider that uses InlineKeyboard buttons for tool execution approval.

#### Scenario: Approval message sent
- **WHEN** a sensitive tool approval is requested for a Telegram session
- **THEN** the system SHALL send a message with InlineKeyboard containing "Approve" and "Deny" buttons to the originating chat

#### Scenario: User approves
- **WHEN** the user clicks the "Approve" button
- **THEN** the callback query SHALL be answered
- **AND** the original message SHALL be edited to show approval status
- **AND** the tool execution SHALL proceed

#### Scenario: User denies
- **WHEN** the user clicks the "Deny" button
- **THEN** the callback query SHALL be answered
- **AND** the original message SHALL be edited to show denial status
- **AND** the tool execution SHALL be denied

#### Scenario: Approval timeout
- **WHEN** no button is clicked within the timeout period
- **THEN** the approval request SHALL be denied with a timeout error

#### Scenario: Context cancellation
- **WHEN** the request context is cancelled before a response
- **THEN** the approval request SHALL return the context error

### Requirement: Callback query routing
The Telegram channel event loop SHALL route CallbackQuery updates to the approval provider before processing regular messages.

#### Scenario: CallbackQuery received
- **WHEN** an update contains a CallbackQuery
- **THEN** the update SHALL be routed to the approval provider's HandleCallback method
- **AND** regular message processing SHALL be skipped for that update

### Requirement: BotAPI Request method
The BotAPI interface SHALL include a `Request` method for operations that return an APIResponse (e.g., AnswerCallbackQuery).

#### Scenario: Answer callback query
- **WHEN** a callback query needs to be acknowledged
- **THEN** the system SHALL use the `Request` method to send an AnswerCallbackQuery

### Requirement: Approval message editing on timeout and cancellation
The Telegram approval provider SHALL edit the approval message to display "Expired" status and remove inline keyboard buttons when the approval times out or the context is cancelled. The inline keyboard markup SHALL be constructed as a struct literal with an explicitly empty `InlineKeyboard` slice (`[][]InlineKeyboardButton{}`) to ensure JSON serialization produces `[]` rather than `null`.

#### Scenario: Timeout removes buttons
- **WHEN** an approval request times out without user response
- **THEN** the system SHALL edit the original message to "🔐 Tool approval — ⏱ Expired" with an empty inline keyboard markup serialized as `"inline_keyboard": []`

#### Scenario: Context cancellation removes buttons
- **WHEN** the approval request context is cancelled
- **THEN** the system SHALL edit the original message to "🔐 Tool approval — ⏱ Expired" with an empty inline keyboard markup serialized as `"inline_keyboard": []`

#### Scenario: Approval status removes buttons
- **WHEN** a user clicks the Approve or Deny button
- **THEN** the system SHALL edit the original message to show approval/denial status with an empty inline keyboard markup serialized as `"inline_keyboard": []`

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
