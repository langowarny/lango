## MODIFIED Requirements

### Requirement: Chat message handling
The gateway handleChatMessage handler SHALL use RunStreaming instead of RunAndCollect to process chat messages. During execution, each LLM text chunk SHALL be broadcast as an `agent.chunk` WebSocket event to the session's clients.

#### Scenario: Streaming chat response
- **WHEN** a chat.message RPC is received
- **THEN** the gateway SHALL call RunStreaming with a chunk callback that broadcasts `agent.chunk` events

#### Scenario: Full response in RPC result
- **WHEN** streaming completes successfully
- **THEN** the RPC response SHALL contain `{"response": "<full text>"}` for backward compatibility
