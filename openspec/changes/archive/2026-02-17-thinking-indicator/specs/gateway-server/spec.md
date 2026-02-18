## ADDED Requirements

### Requirement: Session-scoped broadcast
The Gateway server SHALL provide a `BroadcastToSession` method that sends events only to UI clients matching a specific session key. When session key is empty (no auth), it SHALL broadcast to all UI clients.

#### Scenario: Authenticated session broadcast
- **WHEN** `BroadcastToSession` is called with a non-empty session key
- **THEN** only UI clients with a matching `SessionKey` SHALL receive the event
- **AND** companion clients SHALL NOT receive the event

#### Scenario: Unauthenticated broadcast
- **WHEN** `BroadcastToSession` is called with an empty session key
- **THEN** all UI clients SHALL receive the event

### Requirement: Agent thinking events
The Gateway server SHALL broadcast `agent.thinking` before agent processing and `agent.done` after processing completes, scoped to the requesting user's session.

#### Scenario: Thinking event on message receipt
- **WHEN** a `chat.message` RPC is received
- **THEN** the server SHALL broadcast an `agent.thinking` event to the session before calling `RunAndCollect`

#### Scenario: Done event after processing
- **WHEN** `RunAndCollect` returns (success or error)
- **THEN** the server SHALL broadcast an `agent.done` event to the session
