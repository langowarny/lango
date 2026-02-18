## Context

Rewrite OpenClaw (TypeScript) in Go to build a faster and safer personal AI agent. Maintain OpenClaw's core architecture (Gateway + Agent + Channels + Tools) while leveraging the advantages of the Go ecosystem.

**Current State:**
- OpenClaw: ~2,500 files, 50+ npm dependencies, Node.js 22 required
- Performance Issues: Startup time 1-3s, Memory 100-300MB

**Constraints:**
- ADK-Go is still in early stages, some features need direct implementation
- WhatsApp (Baileys) porting to Go is very complex → Excluded from initial version
- SQLite CGO dependency prevents pure Go builds

## Goals / Non-Goals

**Goals:**
- Single binary distribution (~20-50MB)
- Startup time < 100ms, Memory < 50MB
- Support for Telegram, Discord, Slack channels
- OpenClaw compatible Gateway protocol
- macOS, Linux cross-compilation

**Non-Goals:**
- WhatsApp, iMessage support (Internal version)
- iOS/Android native apps
- Canvas/A2UI real-time visualization
- Voice Wake function (Initial version)

## Decisions

### D1: Project Structure

```
lango/
├── cmd/
│   └── lango/           # CLI entrypoint
├── internal/
│   ├── agent/           # ADK-Go based agent
│   ├── gateway/         # WebSocket server
│   ├── channels/        # Telegram, Discord, Slack
│   ├── tools/           # exec, filesystem, browser
│   ├── session/         # SQLite session store
│   └── config/          # Configuration loading
├── pkg/                 # Public packages
└── go.mod
```

**Rationale:** Go standard layout. Encapsulation with `internal/`, separation of entrypoints with `cmd/`.

### D2: Agent Framework - ADK-Go

**Selection:** `github.com/google/adk-go`

**Alternatives Reviewed:**
- LangChainGo: Lack of documentation, small community
- Custom Implementation: High time cost, complex tool calling

**Rationale:** Official Google support, built-in A2A protocol, multi-agent support.

### D3: Logging - Uber Zap

**Selection:** `go.uber.org/zap`

**Alternatives Reviewed:**
- log/slog: Standard but lower performance
- zerolog: Similar to Zap, smaller community

**Rationale:** Structured logging, high performance, production proven.

### D4: HTTP/WebSocket

**Selection:** `chi/v5` + `gorilla/websocket`

**Alternatives Reviewed:**
- Gin: Heavy, many unnecessary features
- Fiber: fasthttp based, compatibility issues

**Rationale:** chi is lightweight and standard net/http compatible. gorilla/websocket is effectively the standard.

### D5: Browser Automation - rod

**Selection:** `github.com/go-rod/rod`

**Alternatives Reviewed:**
- chromedp: Low-level, lots of boilerplate
- go-playwright: Experimental, unstable

**Rationale:** Playwright style API, active maintenance, automatic browser management.

### D6: Session Store - SQLite

**Selection:** `github.com/mattn/go-sqlite3` (CGO)

**Alternatives Reviewed:**
- modernc.org/sqlite: Pure Go, lower performance
- BadgerDB: Key-value only, inconvenient querying

**Rationale:** OpenClaw compatible, supports complex queries, proven stability. Requires CGO but cross-compilation is possible.

### D7: Channel SDK

| Channel | Library | Rationale |
|---------|---------|-----------|
| Telegram | `telegram-bot-api/v5` | Most widely used, rich documentation |
| Discord | `discordgo` | Official level, slash command support |
| Slack | `slack-go/slack` | Bolt pattern support, Events API |

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| ADK-Go Immaturity | Direct extension if needed, community contribution |
| Complex CGO Build | Docker build environment, CI cache |
| OpenClaw Protocol Compatibility | Incremental implementation, core methods first |
| Channel SDK Bugs | Pin stable versions, fallback logic |

## Open Questions

- [ ] Is ADK-Go session persistence compatible with SQLite?
- [ ] Memory buffer size limits when handling Telegram media?
- [ ] How to resolve SQLite CGO during cross-compilation (zig cc vs docker)?
