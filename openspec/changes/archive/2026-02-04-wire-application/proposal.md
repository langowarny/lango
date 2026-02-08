## Why

The current `main.go` only starts the gateway server, leaving critical components (ADK Agent, Ent Session Store, Tools, and Channels) initialized in isolation or not at all. This results in a non-functional system where the agent cannot process messages or maintain context. Wiring these components together is essential to create a runnable application.

## What Changes

- Create `internal/app` package to manage the application lifecycle and component wiring.
- Update `internal/ent` to initialize the SQLite client correctly.
- Implement an adapter to bridge `internal/ent` with `agent.SessionStore` interface.
- Update `cmd/lango/main.go` to use `internal/app` for bootstrapping.
- Register `internal/tools` (Browser, Filesystem, Exec) with the Agent.
- Connect `internal/channels` (Telegram, Discord, Slack) to the Agent via the App layer.

## Capabilities

### New Capabilities
- `application-core`: Defines the application bootstrapping, component lifecycle management, and wiring logic.

### Modified Capabilities
- `gateway-server`: Update to accept an Agent instance for handling RPC requests (currently just `chat.message`).

## Impact

- `cmd/lango`: Completely refactored to delegate to `internal/app`.
- `internal/gateway`: API update to support Agent injection.
- New `internal/app` package introduced.
- New `internal/session/ent_adapter.go` (or similar) created.
