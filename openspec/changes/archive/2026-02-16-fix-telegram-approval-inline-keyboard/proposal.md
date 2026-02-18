## Why

When a user clicks the Approve or Deny inline keyboard button in Telegram, the `editApprovalMessage` function attempts to remove the keyboard by calling `tgbotapi.NewInlineKeyboardMarkup()` with no arguments. This produces a nil `InlineKeyboard` field, which serializes to `"inline_keyboard": null` in JSON. The Telegram Bot API requires `inline_keyboard` to be an Array, so the edit request fails with `Bad Request: field "inline_keyboard" must be of type Array`. As a result, the approval flow appears unresponsive — the button click has no visible effect and the approval message is never updated.

## What Changes

- Fix `editApprovalMessage` in `internal/channels/telegram/approval.go` to construct an `InlineKeyboardMarkup` with an explicitly empty slice (`[][]InlineKeyboardButton{}`) instead of calling `NewInlineKeyboardMarkup()` with zero arguments, ensuring the field serializes to `[]` rather than `null`.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `channel-telegram`: Fix inline keyboard removal on approval message edit to comply with Telegram Bot API array requirement.

## Impact

- **Code**: `internal/channels/telegram/approval.go` — `editApprovalMessage` method (single-line fix).
- **APIs**: No API contract changes; this fixes compliance with the Telegram Bot API expectation.
- **Dependencies**: No dependency changes.
- **Systems**: Telegram approval buttons will now correctly update the message after being clicked.
