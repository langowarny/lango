## 1. Dependencies

- [x] 1.1 Add `entgo.io/ent` to go.mod
- [x] 1.2 Add `modernc.org/sqlite` to go.mod
- [x] 1.3 Remove `github.com/mattn/go-sqlite3` from go.mod
- [x] 1.4 Run `go mod tidy`

## 2. Ent Schema

- [x] 2.1 Initialize ent in `internal/ent/` with `go run entgo.io/ent/cmd/ent init Session Message`
- [x] 2.2 Define Session schema with fields: key, agentID, channelType, channelID, model, metadata, createdAt, updatedAt
- [x] 2.3 Define Message schema with fields: role, content, timestamp, toolCalls (JSON)
- [x] 2.4 Add Session → Message edge (one-to-many)
- [x] 2.5 Run `go generate ./internal/ent` to generate code

## 3. Store Implementation

- [x] 3.1 Create `internal/session/ent_store.go` implementing `Store` interface
- [x] 3.2 Implement `NewEntStore(dbPath)` with ent client initialization
- [x] 3.3 Implement `Create(session)` using ent client
- [x] 3.4 Implement `Get(key)` with Message eager loading
- [x] 3.5 Implement `Update(session)` 
- [x] 3.6 Implement `Delete(key)`
- [x] 3.7 Implement `AppendMessage(key, msg)`
- [x] 3.8 Implement `Close()`

## 4. Migration

- [x] 4.1 Update `NewSQLiteStore` to return `*EntStore` (keep function name for compatibility)
- [x] 4.2 Remove old raw SQL implementation from `store.go`
- [x] 4.3 Keep `Session`, `Message`, `ToolCall` types as-is (used by other packages)

## 5. Build Updates

- [x] 5.1 Update Dockerfile to remove libsqlite3-dev
- [x] 5.2 Update CI workflow to use CGO_ENABLED=0
- [x] 5.3 Update Makefile for CGO-free builds

## 6. Testing

- [x] 6.1 Run existing session store tests
- [x] 6.2 Verify build works with CGO_ENABLED=0
- [x] 6.3 Test cross-compilation for linux/arm64
