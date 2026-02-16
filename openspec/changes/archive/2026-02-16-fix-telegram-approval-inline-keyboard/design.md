## Context

The `editApprovalMessage` method in `internal/channels/telegram/approval.go` constructs an empty `InlineKeyboardMarkup` by calling `tgbotapi.NewInlineKeyboardMarkup()` with no arguments. The underlying function uses a variadic `rows ...[]InlineKeyboardButton` parameter, which results in a nil `InlineKeyboard` slice when no rows are passed. When serialized to JSON, this produces `"inline_keyboard": null` instead of `"inline_keyboard": []`, violating the Telegram Bot API contract.

## Goals / Non-Goals

**Goals:**
- Ensure `editApprovalMessage` sends a valid empty array for `inline_keyboard`, complying with the Telegram Bot API.
- Approval button clicks result in visible message updates (Approved/Denied/Expired).

**Non-Goals:**
- Changing the approval flow logic or adding new approval features.
- Modifying the `go-telegram-bot-api` library itself.

## Decisions

**Decision: Use struct literal with empty slice instead of `NewInlineKeyboardMarkup()`**

Replace:
```go
emptyMarkup := tgbotapi.NewInlineKeyboardMarkup()
```

With:
```go
emptyMarkup := tgbotapi.InlineKeyboardMarkup{
    InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
}
```

**Rationale**: This guarantees `InlineKeyboard` is a non-nil empty slice, which serializes to `[]` in JSON. The alternative of passing an empty row (`NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow())`) would produce `[[]]` which is semantically incorrect. Using the struct literal is explicit and self-documenting.

## Risks / Trade-offs

- **[Low] Library update could change behavior** → If `go-telegram-bot-api` changes `NewInlineKeyboardMarkup()` to initialize with an empty slice, this workaround becomes redundant but not harmful.
