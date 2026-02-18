## Why

In multi-agent mode, there is no runtime visibility into whether the orchestrator handles a request directly or delegates to a sub-agent. ADK `session.Event` carries `Author` and `Actions.TransferToAgent` fields, but no runtime logging exists to surface agent routing during development and debugging.

## What Changes

- Add debug-level event stream logging in `RunAndCollect()` to trace agent delegation and response events
- Add a `logger()` helper and `hasText()` utility in the `internal/adk` package
- Log `Author` and `TransferToAgent` fields from each event in the ADK event stream

## Capabilities

### New Capabilities
- `agent-event-logging`: Debug-level logging of multi-agent event streams for runtime observability of agent delegation and response routing

### Modified Capabilities

## Impact

- `internal/adk/agent.go` — new logging in `RunAndCollect()`, new helper functions
- New dependency on `internal/logging` from `internal/adk` package
- Debug-level only — zero impact on production log output unless log level is set to `debug`
