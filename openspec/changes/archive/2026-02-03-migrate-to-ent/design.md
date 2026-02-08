## Context

Lango session store currently uses `mattn/go-sqlite3` with raw SQL queries. This requires CGO which complicates cross-compilation and CI. We're migrating to entgo.io with `modernc.org/sqlite` (pure Go) to simplify builds while adding type-safe queries.

## Goals

1. Eliminate CGO dependency for SQLite
2. Maintain existing `Store` interface compatibility
3. Add type-safe queries via ent code generation
4. Enable simpler cross-compilation and Docker builds

## Decisions

### 1. Entity Framework: entgo.io

**Options considered**:
- GORM - Popular but heavy, magic behavior
- sqlc - SQL-first, good but raw SQL
- ent - Type-safe, code generation, Facebook/Meta backed

**Decision**: entgo.io
- Type-safe fluent API
- Automatic migrations
- Works with `modernc.org/sqlite` (pure Go)
- Good documentation and active maintenance

### 2. SQLite Driver: modernc.org/sqlite

**Options considered**:
- `mattn/go-sqlite3` - CGO, widely used
- `modernc.org/sqlite` - Pure Go, CGO-free
- `crawshaw.io/sqlite` - Pure Go, less maintained

**Decision**: `modernc.org/sqlite`
- Pure Go, no CGO required
- Compatible with ent
- Active maintenance

### 3. Schema Design

**Entities**:
- `Session` - Main session entity with fields matching current struct
- `Message` - Conversation messages, linked to Session

**Relationships**:
- Session has many Messages (one-to-many)

### 4. Migration Strategy

- Keep existing `Store` interface unchanged
- Create new `EntStore` implementing same interface
- Replace `SQLiteStore` with `EntStore` in `NewSQLiteStore` function
- Existing tests should pass without modification

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Pure Go SQLite performance | Low | Medium | Benchmark before/after, acceptable for session storage |
| Ent learning curve | Low | Low | Good docs, straightforward schema |
| Migration path for existing DBs | Medium | Low | Ent auto-migrate handles schema changes |

## Open Questions

1. ~~Should we support both drivers?~~ → No, fully commit to pure Go
2. ~~Keep raw SQL as fallback?~~ → No, clean migration to ent
