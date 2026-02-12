## MODIFIED Requirements

### Requirement: File Structure
The gateway package SHALL be organized into focused files where no single file exceeds 200 lines of non-test code. The decomposition SHALL be: `server.go` (server lifecycle, types, routes, health), `websocket.go` (WebSocket connection management, Client struct, broadcast), `handlers.go` (RPC handler implementations), `auth.go` (OIDC authentication, already separated).

#### Scenario: File Size Compliance
- **GIVEN** the gateway package after decomposition
- **WHEN** line counts are measured for non-test `.go` files
- **THEN** no single file SHALL exceed 200 lines

#### Scenario: Build and Test Integrity
- **GIVEN** the decomposed gateway package
- **WHEN** `go build ./internal/gateway/...` is run
- **THEN** it SHALL succeed with no errors
- **WHEN** `go test ./internal/gateway/...` is run
- **THEN** all existing tests SHALL pass without modification

### Requirement: server.go (Core Server)
The `server.go` file SHALL contain the Server struct definition, Config struct, RPC protocol types (RPCRequest, RPCResponse, RPCError, RPCHandler), the constructor `New()`, route setup, handler registration, server Start/Shutdown lifecycle, and HTTP endpoint handlers (health, status). It SHALL NOT contain WebSocket connection management code or RPC handler business logic.

#### Scenario: Server Constructor
- **WHEN** `gateway.New()` is called with config, agent, provider, store, and auth parameters
- **THEN** it SHALL return a fully initialized Server
- **THEN** it SHALL register all RPC handlers (chat.message, sign.response, encrypt.response, decrypt.response, companion.hello, approval.response)
- **THEN** it SHALL wire up the provider sender if provider is non-nil

#### Scenario: Server Lifecycle
- **WHEN** `Start()` is called
- **THEN** it SHALL listen on the configured host:port
- **WHEN** `Shutdown()` is called
- **THEN** it SHALL close all WebSocket clients and stop the HTTP server

### Requirement: websocket.go (Connection Management)
The `websocket.go` file SHALL contain the Client struct, WebSocket upgrade handlers, read/write pump goroutines, send helpers, client close logic, broadcast methods, and client removal. All WebSocket protocol handling SHALL be isolated in this file.

#### Scenario: Client Connection
- **WHEN** a client connects to `/ws`
- **THEN** a Client SHALL be created with type "ui" and registered in the server's client map
- **WHEN** a client connects to `/companion`
- **THEN** a Client SHALL be created with type "companion"

#### Scenario: RPC Dispatch
- **WHEN** a Client receives a JSON message in readPump
- **THEN** it SHALL parse it as RPCRequest and dispatch to the matching registered handler
- **THEN** it SHALL send back RPCResponse with the handler result or error

#### Scenario: Broadcast
- **WHEN** `Broadcast()` is called
- **THEN** the message SHALL be sent to all clients with type "ui"
- **WHEN** `BroadcastToCompanions()` is called
- **THEN** the message SHALL be sent to all clients with type "companion"

### Requirement: handlers.go (RPC Handlers)
The `handlers.go` file SHALL contain all RPC handler method implementations: chat message processing, security response proxying (sign, encrypt, decrypt), companion hello, approval response handling, and the `RequestApproval` method.

#### Scenario: Chat Message Handler
- **WHEN** `handleChatMessage` receives a valid message with sessionKey
- **THEN** it SHALL execute the ADK agent via `agent.Run()` and return the aggregated text response
- **WHEN** message is empty
- **THEN** it SHALL return an error

#### Scenario: Security Proxy Handlers
- **WHEN** `handleSignResponse`, `handleEncryptResponse`, or `handleDecryptResponse` is called
- **THEN** it SHALL delegate to the corresponding RPCProvider method
- **WHEN** provider is nil
- **THEN** it SHALL return "provider not configured" error

#### Scenario: Approval Workflow
- **WHEN** `RequestApproval` is called
- **THEN** it SHALL broadcast an approval request to companions
- **THEN** it SHALL wait for a response or timeout after 30 seconds
- **WHEN** no companion is connected
- **THEN** it SHALL return an error immediately
