## Context

The automation delivery pipeline (cron, background, workflow) uses a string passthrough architecture: tools auto-detect or accept a delivery target string, which flows unchanged through the pipeline until `channelSender.SendMessage` consumes it. Currently, `detectChannelFromContext` extracts only the channel name (e.g. "telegram") from session keys like `telegram:123456789:987654321`, discarding the chat/channel ID needed for actual message delivery.

## Goals / Non-Goals

**Goals:**
- Deliver automation results to the correct chat/channel by including the target ID in the delivery string
- Maintain backward compatibility with bare channel names (e.g. existing DB records with `["telegram"]`)
- Keep the string passthrough pipeline unchanged — modify only the producer (`detectChannelFromContext`) and consumer (`SendMessage`)

**Non-Goals:**
- Changing the DB schema, config types, or intermediate pipeline components
- Adding new channel adapters or delivery mechanisms
- Implementing delivery retry or queue mechanisms

## Decisions

### Decision 1: Delivery target format `channel:id`

Extend the delivery target string from bare `"telegram"` to `"telegram:123456789"`. This is parsed at consumption time by `parseDeliveryTarget`.

**Rationale**: The simplest change that preserves the entire passthrough pipeline. All intermediate layers (cron delivery, background notification, workflow engine, config, DB) already store and pass strings without interpretation — only the endpoints need to change.

**Alternatives considered**:
- Struct-based delivery target: Would require changing every intermediate type and interface — disproportionate effort for the fix
- Separate target ID field on all tool parameters: Would require schema changes and DB migration

### Decision 2: Backward compatibility via allowlist fallback

When Telegram receives a bare channel name without target ID, it falls back to the existing `firstTelegramChatID()` allowlist behavior. Discord and Slack return explicit errors when no target ID is provided (they already failed with empty channel IDs).

**Rationale**: Preserves existing behavior for DB records that predate this change while making failure modes explicit for Discord/Slack.

## Risks / Trade-offs

- [Telegram allowlist fallback may deliver to wrong chat] → Acceptable: same behavior as before the fix; explicit target ID is preferred and auto-detected going forward
- [Discord/Slack now return explicit errors instead of silent failures] → Improvement: surfaces the real problem instead of hiding it behind adapter-level errors
