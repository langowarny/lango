## Why

`internal/gateway/server.go` is a 575-line file that handles six distinct responsibilities: HTTP server lifecycle, WebSocket connection management (Client struct, read/write pumps, broadcast), JSON-RPC dispatch, chat message processing (agent execution + response aggregation), security RPC proxying (sign/encrypt/decrypt responses), and companion approval workflow. This violates the project's own guideline that "no single file SHALL exceed 200 lines" (established in the mvp-cleanup application-core spec). The monolithic structure makes it difficult to test individual concerns, modify one subsystem without risk to others, and onboard new contributors.

## What Changes

- Extract WebSocket connection management into `websocket.go`: Client struct, readPump, writePump, sendResult, sendError, Close, handleWebSocketConnection, handleWebSocket, handleCompanionWebSocket, Broadcast, BroadcastToCompanions, broadcastToType, removeClient.
- Extract RPC handler implementations into `handlers.go`: handleChatMessage, handleSignResponse, handleEncryptResponse, handleDecryptResponse, handleCompanionHello, handleApprovalResponse, RequestApproval.
- Retain in `server.go`: Server struct, Config, RPCRequest/RPCResponse/RPCError types, RPCHandler type, New(), setupRoutes(), RegisterHandler(), Start(), Shutdown(), handleHealth, handleStatus, logger().
- No public API changes. All exported methods and types remain on the Server struct. Import paths unchanged.
- No behavioral changes. This is a pure structural refactoring.

## Capabilities

### New Capabilities
_None. This change restructures existing code without adding features._

### Modified Capabilities
- `gateway-server`: Decompose server.go (575 lines) into three focused files: `server.go` (~180 lines, core server lifecycle + types), `websocket.go` (~180 lines, connection management + broadcast), `handlers.go` (~180 lines, RPC handler implementations). Auth remains in `auth.go` (already separated).

## Impact

- **Code**: `internal/gateway/server.go` split into 3 files; no new packages or dependencies
- **Tests**: `server_test.go` unchanged (tests Server through public API); may add focused handler tests later
- **API**: No changes to exported types, methods, or constructor signatures
- **Risk**: Low - pure file-level extraction with no logic changes
