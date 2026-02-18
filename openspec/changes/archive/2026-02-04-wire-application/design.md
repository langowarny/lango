## Context

The current `lango` application consists of several well-implemented but isolated components: `agent` (ADK), `gateway` (Server), `tools` (Browser, FS, Exec), `channels` (Telegram, Discord, Slack), and `ent` (Database Schema). The entry point `main.go` only initializes the Gateway, lacking the glue code to assemble these parts into a cohesive system. Detailed analysis revealed that no wiring exists to inject dependencies or manage the lifecycle of these components.

## Goals / Non-Goals

**Goals:**
- Centralize application bootstrap logic in a new `internal/app` package.
- Implement proper dependency injection for the Agent Runtime (DB, Tools).
- Connect the Agent to the Gateway for RPC handling.
- Connect the Agent to Channels for message processing.
- Ensure graceful shutdown for all components.

**Non-Goals:**
- Adding new features to existing components (e.g., new tools or channel features).
- Changing the internal implementation of `agent` or `gateway` (except for necessary API exposure).

## Decisions

### 1. Introduce `internal/app` Package
**Decision:** Create a `App` struct in `internal/app` that holds references to all major components (`Agent`, `Gateway`, `SessionStore`, etc.) and manages their lifecycle (`Start`, `Stop`).
**Rationale:** Keeps `main.go` clean and makes the application logic testable and reusable.
**Alternatives:** Put all wiring logic directly in `main.go`. Rejected as it scales poorly and is hard to test.

### 2. Ent as Session Store Implementation
**Decision:** Implement an adapter in `internal/session` or `internal/app` that wraps the `ent` client to satisfy `agent.SessionStore` interface.
**Rationale:** `internal/ent` is generated but doesn't implement the interface directly. An adapter is needed to bridge the gap without modifying generated code.
**Alternatives:** Modify `ent` schema to generate compatible code (complex).

### 3. Agent-Centric Wiring
**Decision:** The `Agent` is the central component. Tools are injected into the Agent. Channels and Gateway hold a reference to the Agent to dispatch messages/RPCs.
**Rationale:** Reflects the data flow: Inputs (Channels/Gateway) -> Agent -> Tools/DB.

### 4. Explicit Tool Registration in Config
**Decision:** Register tools based on their presence in the `lango.json` configuration.
**Rationale:** Allows enabling/disabling tools via config without recompiling.

## Risks / Trade-offs

- **Risk:** Circular dependencies if components reference each other effectively.
  - **Mitigation:** Use strict dependency injection from `app` layer; components should prefer interfaces.
- **Risk:** Complexity of `ent` initialization (driver, migration).
  - **Mitigation:** Use `ent`'s standard initialization patterns and ensure migration runs on startup.

## Migration Plan
This is a new wiring for an existing codebase, essentially "turning on" the full application. No data migration is needed as the previous version didn't use the DB effectively for this purpose. Deployment will require a valid `lango.json` with all sections populated.
