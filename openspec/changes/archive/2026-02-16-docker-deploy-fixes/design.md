## Context

When lango runs inside a Docker container, three issues prevent normal operation:
1. The approval chain (GatewayProvider → TTYProvider) fails entirely — no companion is connected and stdin is not a TTY, so all sensitive tool executions are denied.
2. go-rod cannot find the system-installed Chromium because `ENV ROD_BROWSER` is not a recognized go-rod env var. go-rod uses `launcher.LookPath()` or an explicit `Bin()` call.
3. `WORKDIR /app` is root-owned, but the process runs as `USER lango`, causing write permission errors.

## Goals / Non-Goals

**Goals:**
- Enable tool execution approval in headless environments via an opt-in `HeadlessAutoApprove` config
- Auto-detect system-installed browser binaries using `launcher.LookPath()`
- Fix container filesystem permissions so the non-root user can write to WORKDIR

**Non-Goals:**
- Changing the fail-closed default behavior of the approval system
- Supporting multiple browser binaries or browser selection UI
- Adding new Docker Compose features or secrets management changes

## Decisions

### Decision 1: HeadlessProvider as TTY fallback slot

Reuse the existing `CompositeProvider.SetTTYFallback()` slot instead of adding a new provider chain level. When `HeadlessAutoApprove` is true, wire `HeadlessProvider` into the TTY fallback position; otherwise wire `TTYProvider` (existing behavior).

**Rationale**: Minimal code change. The TTY fallback is exactly the slot that fails in Docker (TTYProvider returns false when stdin isn't a terminal). Replacing it is the natural fix point. Channel-based providers (Telegram/Discord/Slack) and GatewayProvider still take priority when available.

**Alternative considered**: Adding HeadlessProvider as a registered provider with lowest priority. Rejected because `CanHandle()` semantics don't fit — headless approval is not prefix-based routing.

### Decision 2: launcher.LookPath() for browser detection

Use go-rod's `launcher.LookPath()` to auto-detect system Chromium. Fall back to go-rod's default download behavior only if both config `BrowserBin` and LookPath fail.

**Rationale**: `LookPath()` is go-rod's official API for finding system browsers. It checks standard paths (`/usr/bin/chromium`, `/usr/bin/google-chrome`, etc.) across all platforms. The removed `ENV ROD_BROWSER` was never read by go-rod.

### Decision 3: WORKDIR /home/lango

Change WORKDIR from `/app` to `/home/lango` (the user's home directory created by `useradd -m -d /home/lango`).

**Rationale**: The binary is at `/usr/local/bin/lango`, so `/app` serves no purpose. `/home/lango` is owned by the lango user and writable by default.

## Risks / Trade-offs

- **[Risk] HeadlessAutoApprove bypasses human approval** → Mitigation: default is `false`; every auto-approval is logged at WARN level for audit; channel-based providers still take priority when configured.
- **[Risk] LookPath may find unexpected browser** → Mitigation: `BrowserBin` config allows explicit override; LookPath is go-rod's standard detection.
- **[Risk] Existing deployments using WORKDIR /app** → Mitigation: `/app` was only used as CWD; binary path unchanged at `/usr/local/bin/lango`; data volume at `/data` unchanged.
