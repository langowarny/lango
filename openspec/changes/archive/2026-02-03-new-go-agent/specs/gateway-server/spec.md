## ADDED Requirements

### Requirement: WebSocket server initialization
The system SHALL start a WebSocket server on a configurable port (default 18789) for client connections.

#### Scenario: Server startup on default port
- **WHEN** the gateway starts without port configuration
- **THEN** the server SHALL listen on port 18789

#### Scenario: Server startup on custom port
- **WHEN** the gateway starts with a custom port specified
- **THEN** the server SHALL listen on the specified port

### Requirement: Client connection management
The system SHALL accept and manage multiple WebSocket client connections simultaneously.

#### Scenario: New client connection
- **WHEN** a client connects to the WebSocket endpoint
- **THEN** the system SHALL assign a unique connection ID and track the connection

#### Scenario: Client disconnection
- **WHEN** a client disconnects
- **THEN** the system SHALL clean up associated resources and remove from connection pool

### Requirement: RPC method handling
The system SHALL process JSON-RPC style method calls from connected clients.

#### Scenario: Valid method invocation
- **WHEN** a client sends a valid RPC request with method name and parameters
- **THEN** the system SHALL execute the corresponding handler and return the result

#### Scenario: Unknown method
- **WHEN** a client sends an RPC request for an unknown method
- **THEN** the system SHALL return a method-not-found error

### Requirement: Event broadcasting
The system SHALL broadcast events to all connected clients or specific client subsets.

#### Scenario: Broadcast to all clients
- **WHEN** a system event occurs (e.g., session update)
- **THEN** the event SHALL be sent to all connected clients

#### Scenario: Targeted message delivery
- **WHEN** a message is intended for a specific session
- **THEN** only clients subscribed to that session SHALL receive it

### Requirement: HTTP API endpoints
The system SHALL expose HTTP REST endpoints for non-WebSocket operations.

#### Scenario: Health check endpoint
- **WHEN** a GET request is made to /health
- **THEN** the system SHALL return 200 OK with status information
