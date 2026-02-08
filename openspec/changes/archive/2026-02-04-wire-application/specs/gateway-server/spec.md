## ADDED Requirements

### Requirement: Agent RPC Handler
The Gateway server SHALL support handling RPC requests by delegating to an injected Agent Runtime.

#### Scenario: Handling chat.message RPC
- **WHEN** the Gateway receives a "chat.message" RPC request
- **THEN** it SHALL invoke the injected Agent's `Run` method with the message content
- **AND** it SHALL return the Agent's response as the RPC result

#### Scenario: Agent dependency injection
- **WHEN** the Gateway server is created
- **THEN** it SHALL accept an Agent Runtime instance (interface) as a dependency
