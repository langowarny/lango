## MODIFIED Requirements

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
