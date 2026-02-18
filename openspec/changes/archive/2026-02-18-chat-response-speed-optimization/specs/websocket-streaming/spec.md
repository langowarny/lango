## ADDED Requirements

### Requirement: Agent streaming execution
The Agent SHALL provide a `RunStreaming` method that streams partial text chunks via a callback function and returns the full accumulated response.

#### Scenario: Streaming with partial events
- **WHEN** the LLM produces partial (streaming) events
- **THEN** each text chunk SHALL be passed to the ChunkCallback and accumulated into the full response

#### Scenario: Non-streaming fallback
- **WHEN** the LLM produces only non-partial events (no streaming support)
- **THEN** the full response text SHALL be collected from non-partial events without invoking the callback

### Requirement: WebSocket chunk events
The gateway SHALL broadcast `agent.chunk` events to session clients during streaming agent execution. Each chunk event SHALL contain the sessionKey and the text chunk.

#### Scenario: Client receives streaming chunks
- **WHEN** the agent produces streaming text chunks
- **THEN** the gateway SHALL broadcast an `agent.chunk` event with `{sessionKey, chunk}` payload for each chunk

#### Scenario: Backward-compatible RPC response
- **WHEN** streaming completes
- **THEN** the RPC response SHALL still contain the full `response` text

#### Scenario: Existing events preserved
- **WHEN** a chat message is processed
- **THEN** `agent.thinking` SHALL be sent before execution and `agent.done` after completion
