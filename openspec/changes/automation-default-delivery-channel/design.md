## Context

Automation systems (cron, background, workflow) require explicit delivery channel specification. When omitted, results are stored but never delivered to users. The session key already contains channel information (e.g., `telegram:12345`), making auto-detection straightforward.

## Goals / Non-Goals

**Goals:**
- Auto-detect delivery channel from session key prefix for all three automation systems
- Provide config-level defaults as a second fallback layer
- Warn operators when no delivery channel is resolved at all
- Expose default delivery settings in Settings TUI

**Non-Goals:**
- Per-user delivery preferences (out of scope)
- Multi-channel delivery auto-detection (only detect the originating channel)
- Changing the delivery mechanism itself (existing Delivery/Notification/ChannelSender interfaces unchanged)

## Decisions

1. **Three-tier fallback chain**: Explicit param → session auto-detect → config default → warn
   - Rationale: Maximizes backward compatibility while eliminating silent drops. Explicit params always win.

2. **Session key prefix parsing**: Extract channel from `strings.SplitN(sessionKey, ":", 2)[0]` and match against known channels (`telegram`, `discord`, `slack`)
   - Alternative: Pass channel as a separate context value → rejected because session key already carries this information reliably.

3. **Config field `DefaultDeliverTo []string`**: Slice type for cron/workflow (multi-channel), background uses `[0]` (single channel)
   - Rationale: Consistent type across all three configs; background only supports one origin channel by design.

4. **Helper placement**: `detectChannelFromContext()` in `internal/app/tools.go` alongside the tool builders
   - Alternative: Separate utility package → rejected as over-engineering for a single helper used in one file.

5. **Warning level for empty delivery**: `Warnw` instead of `Debugw`
   - Rationale: Silent drops are the root problem; warnings make the issue visible in production logs without failing the job.

## Risks / Trade-offs

- [Risk] Auto-detected channel may not be the desired delivery target → Mitigation: Explicit param always takes precedence; auto-detection is only a fallback.
- [Risk] CLI/API sessions have no channel prefix → Mitigation: Config defaults serve as the final fallback; warning log alerts the operator.
- [Trade-off] `DefaultDeliverTo` is a slice but background only uses `[0]` → Accepted for type consistency across configs.
