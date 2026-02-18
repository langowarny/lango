## Tasks

### Task 1: Extract websocket.go
- [ ] Create `internal/gateway/websocket.go`
- [ ] Move Client struct and all its methods (readPump, writePump, sendResult, sendError, Close)
- [ ] Move handleWebSocket, handleCompanionWebSocket, handleWebSocketConnection
- [ ] Move Broadcast, BroadcastToCompanions, broadcastToType
- [ ] Move removeClient
- [ ] Add necessary imports (json, fmt, time, sync, websocket)
- [ ] Verify `go build ./internal/gateway/...`

### Task 2: Extract handlers.go
- [ ] Create `internal/gateway/handlers.go`
- [ ] Move handleChatMessage
- [ ] Move handleSignResponse, handleEncryptResponse, handleDecryptResponse
- [ ] Move handleCompanionHello, handleApprovalResponse
- [ ] Move RequestApproval
- [ ] Add necessary imports (context, json, fmt, strings, time, adk, security)
- [ ] Verify `go build ./internal/gateway/...`

### Task 3: Clean up server.go
- [ ] Remove all code moved to websocket.go and handlers.go
- [ ] Remove imports that are no longer needed in server.go
- [ ] Verify server.go contains only: types, New(), setupRoutes(), RegisterHandler(), Start(), Shutdown(), handleHealth, handleStatus, logger()
- [ ] Verify no file exceeds 200 lines

### Task 4: Verify
- [ ] `go build ./...` passes
- [ ] `go test ./internal/gateway/...` passes
- [ ] `go vet ./internal/gateway/...` passes
- [ ] No changes needed outside `internal/gateway/`
