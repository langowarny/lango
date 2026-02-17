## Context

In multi-agent mode, the ADK runner emits `session.Event` objects through an iterator. Each event carries `Author` (the agent that produced it) and `Actions.TransferToAgent` (delegation target). Currently `RunAndCollect()` only extracts text content, discarding all routing metadata. Developers have no runtime signal to verify whether requests are handled by the orchestrator or delegated to sub-agents.

## Goals / Non-Goals

**Goals:**
- Provide debug-level logging of agent delegation and response events in `RunAndCollect()`
- Enable developers to trace multi-agent routing in real time via log level configuration
- Keep logging centralized in the single event consumption point used by all channels

**Non-Goals:**
- Structured tracing or OpenTelemetry integration
- Logging at info level or above (must remain debug-only)
- Changes to the ADK session storage layer or event schema
- Per-channel logging customization

## Decisions

**1. Log in `RunAndCollect()`, not in `AppendEvent` or per-channel handlers**

Rationale: `RunAndCollect()` is the single event consumption point shared by all channels (Telegram, Slack, Discord, Gateway). Logging here provides universal coverage without code duplication. `AppendEvent` runs inside the ADK runner's internal loop and serves a different purpose (persistence).

**2. Debug level only**

Rationale: Agent routing events are high-frequency in multi-agent mode. Using debug level ensures zero noise in production while allowing full visibility when needed via `log.level: debug`.

**3. Reuse `logging.Agent()` subsystem**

Rationale: The `internal/logging` package already provides a named `agent` subsystem logger. Adding a package-level `logger()` helper in `internal/adk` keeps the call sites clean without introducing new subsystems.

**4. Two log categories: delegation vs response**

Rationale: Delegation events (`TransferToAgent` set) and response events (text content present) represent the two observable routing outcomes. Logging both with distinct messages enables filtering by event type.

## Risks / Trade-offs

- [High-frequency debug logs in complex agent trees] → Debug level ensures they are off by default; no performance impact when disabled
- [New import dependency: `internal/adk` → `internal/logging`] → `internal/logging` is a leaf package with no reverse dependencies; no cycle risk
