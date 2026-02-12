## Architecture

The gateway package currently consists of two files:
- `server.go` (575 lines) - everything except auth
- `auth.go` (200 lines) - OIDC authentication (already well-separated)

After decomposition, it will consist of four files:

```
internal/gateway/
  server.go      (~180 lines)  Core server: types, lifecycle, routes, health
  websocket.go   (~180 lines)  WebSocket: Client, pumps, broadcast, connection mgmt
  handlers.go    (~180 lines)  RPC handlers: chat, security proxy, approval
  auth.go        (200 lines)   OIDC auth (unchanged)
  server_test.go (existing)    Integration tests (unchanged)
  gateway_test.go (existing)   Package-level helpers
```

### File Boundaries

#### server.go (Core Server)
Retains the Server struct definition, configuration types, RPC protocol types, and server lifecycle methods. This is the "skeleton" that other files attach behavior to.

```
Types:
  - Server struct (fields only)
  - Config struct
  - RPCRequest, RPCResponse, RPCError structs
  - RPCHandler type

Functions:
  - logger()
  - New()                    constructor, handler registration, provider wiring
  - setupRoutes()            route registration (delegates to methods in other files)
  - RegisterHandler()        register RPC method handler
  - Start()                  HTTP server start
  - Shutdown()               graceful shutdown (closes WS clients, stops HTTP)
  - handleHealth()           GET /health
  - handleStatus()           GET /status
```

#### websocket.go (Connection Management)
All WebSocket-specific code: the Client struct, connection lifecycle (upgrade, read/write pumps), and broadcast utilities. This file is purely about managing connected clients and message routing.

```
Types:
  - Client struct

Functions:
  - handleWebSocket()                 GET /ws upgrade handler
  - handleCompanionWebSocket()        GET /companion upgrade handler
  - handleWebSocketConnection()       generic upgrade + client registration
  - (Client) readPump()               read loop, RPC dispatch
  - (Client) writePump()              write loop, ping keepalive
  - (Client) sendResult()             send RPC success response
  - (Client) sendError()              send RPC error response
  - (Client) Close()                  close connection safely
  - Broadcast()                       send to all UI clients
  - BroadcastToCompanions()           send to all companion clients
  - broadcastToType()                 internal broadcast by client type
  - removeClient()                    unregister client on disconnect
```

#### handlers.go (RPC Handlers)
All business logic handlers that process RPC requests. Each handler is a `RPCHandler` function registered in `New()`.

```
Functions:
  - handleChatMessage()          process user message via ADK agent
  - handleSignResponse()         proxy sign response to RPCProvider
  - handleEncryptResponse()      proxy encrypt response to RPCProvider
  - handleDecryptResponse()      proxy decrypt response to RPCProvider
  - handleCompanionHello()       companion device registration
  - handleApprovalResponse()     companion approval response
  - RequestApproval()            broadcast approval request + await response
```

### Key Design Decisions

1. **No new types or interfaces**: All functions remain methods on `*Server` or `*Client`. No abstraction layers added.

2. **Client struct moves to websocket.go**: Client is purely a WebSocket concern. Server accesses clients through its `clients` map field, which is defined in server.go but managed in websocket.go.

3. **readPump stays in websocket.go despite calling handlers**: The RPC dispatch logic (looking up handlers from `s.handlers` map) is WebSocket protocol handling, not business logic. The actual handler functions called by dispatch live in handlers.go.

4. **RequestApproval stays in handlers.go**: Although it uses broadcast (from websocket.go), its core logic is approval workflow, not connection management.

5. **Import dependencies**: All three files share the same `package gateway`. No circular dependencies since they all operate on the same Server/Client structs.

### Migration Path

This is a single atomic refactoring. Steps:
1. Create `websocket.go` with Client struct and all WS-related methods extracted verbatim
2. Create `handlers.go` with all RPC handler methods extracted verbatim
3. Remove extracted code from `server.go`
4. Verify: `go build ./...` and `go test ./internal/gateway/...`
5. No changes needed in any file outside `internal/gateway/`
