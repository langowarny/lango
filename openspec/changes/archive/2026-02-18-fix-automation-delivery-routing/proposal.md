## Why

Automation delivery (cron, background, workflow) fails because `detectChannelFromContext` extracts only the channel name (e.g. "telegram") from the session key, discarding the target ID (chat/channel ID). This causes `channelSender.SendMessage` to attempt delivery without routing information, resulting in errors like "telegram delivery requires at least one allowlisted chat ID". The fix extends the delivery target format to `channel:id` for precise routing.

## What Changes

- `detectChannelFromContext` returns `channel:targetID` (e.g. `telegram:123456789`) instead of bare channel name
- `channelSender.SendMessage` parses `channel:id` format to extract routing ID for each adapter
- Telegram falls back to allowlist when no target ID is provided (backward compatible)
- Discord/Slack require target ID for delivery (previously failed silently with empty channel ID)
- Tool parameter descriptions and automation prompt hints updated with format examples

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `automation-delivery-fallback`: Session detection now returns `channel:targetID` instead of bare channel name; `SendMessage` parses target format for routing

## Impact

- `internal/app/tools.go`: `detectChannelFromContext` function signature unchanged but returns `channel:id` format; tool descriptions updated
- `internal/app/sender.go`: New `parseDeliveryTarget` helper; `SendMessage` rewritten with target parsing and per-channel routing logic
- `internal/app/wiring.go`: Automation prompt section updated with format hints
- Backward compatible: bare channel names (e.g. `"telegram"`) still work via allowlist fallback
- No schema/config/DB changes required — existing string passthrough pipeline preserved
