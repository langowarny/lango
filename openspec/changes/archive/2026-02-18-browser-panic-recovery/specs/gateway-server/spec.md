## MODIFIED Requirements

### Requirement: websocket.go (Connection Management)
The `websocket.go` file SHALL contain the Client struct, WebSocket upgrade handlers, read/write pump goroutines, send helpers, client close logic, broadcast methods, and client removal. The `handleWebSocketConnection` function SHALL extract the authenticated session key from the request context via `SessionFromContext` and bind it to `Client.SessionKey`.

#### Scenario: Client Connection with Auth
- **WHEN** an authenticated client connects to `/ws`
- **THEN** a Client SHALL be created with `SessionKey` set from `SessionFromContext(r.Context())`

#### Scenario: Client Connection without Auth
- **WHEN** a client connects to `/ws` with no auth configured
- **THEN** a Client SHALL be created with empty `SessionKey`

#### Scenario: RPC Dispatch
- **WHEN** a Client receives a JSON message in readPump
- **THEN** it SHALL parse it as RPCRequest and dispatch to the matching registered handler
- **THEN** it SHALL pass the client reference as the first argument to the handler
- **THEN** it SHALL send back RPCResponse with the handler result or error

#### Scenario: RPC handler panic isolation
- **WHEN** an RPC handler panics during execution
- **THEN** the readPump SHALL NOT terminate
- **AND** the client SHALL receive an error response with the panic details
- **AND** the panic SHALL be logged

#### Scenario: readPump panic recovery
- **WHEN** the readPump goroutine encounters an unrecoverable panic
- **THEN** the panic SHALL be recovered and logged
- **AND** the client SHALL be cleaned up (removed from server, connection closed)

#### Scenario: writePump panic recovery
- **WHEN** the writePump goroutine encounters an unrecoverable panic
- **THEN** the panic SHALL be recovered and logged
- **AND** the client connection SHALL be closed

#### Scenario: Broadcast
- **WHEN** `Broadcast()` is called
- **THEN** the message SHALL be sent to all clients with type "ui"
- **WHEN** `BroadcastToCompanions()` is called
- **THEN** the message SHALL be sent to all clients with type "companion"
