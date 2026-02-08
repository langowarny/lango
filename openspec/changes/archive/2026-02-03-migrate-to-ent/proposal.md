## Why

The current session store uses `mattn/go-sqlite3` which requires CGO. This creates:
- Complex cross-compilation setup (requires C compiler toolchains)
- Slower CI builds due to CGO
- Docker build complexity (need libsqlite3-dev)
- Limited deployment options (CGO binary dependencies)

Migrating to **entgo.io** with `modernc.org/sqlite` provides:
- Pure Go SQLite driver (no CGO)
- Type-safe queries with code generation
- Automatic schema migrations
- Simpler cross-compilation and Docker builds

## What Changes

- Replace `mattn/go-sqlite3` with `entgo.io/ent` + `modernc.org/sqlite`
- Generate ent schema for Session and Message entities
- Rewrite `SQLiteStore` to use ent client
- Simplify Dockerfile (remove libsqlite3-dev)
- Update CI to remove CGO requirements

## Capabilities

### New Capabilities
- `ent-session-store`: Type-safe session storage using entgo.io with pure Go SQLite driver

### Modified Capabilities
- `session-store`: Replace raw SQL implementation with ent-based implementation, maintaining the same `Store` interface

## Impact

- **Code**: `internal/session/store.go` → `internal/ent/` (generated) + `internal/session/ent_store.go`
- **Dependencies**: Remove `mattn/go-sqlite3`, add `entgo.io/ent`, `modernc.org/sqlite`
- **Build**: CGO_ENABLED=0 now possible, simpler cross-compilation
- **Docker**: Smaller image, no C libraries needed
- **Tests**: Existing tests should pass with minor import changes
