## 1. Internal App Package

- [x] 1.1 Create internal/app/types.go defining App and Component interfaces
- [x] 1.2 Create internal/app/app.go with New() and Start/Stop life-cycle methods
- [x] 1.3 Implement dependency injection logic in App struct

## 2. Session Store Integration

- [x] 2.1 Implement EntSessionStore adapter in internal/session/ent_store.go (Verified existing implementation)
- [x] 2.2 Wire Ent initialization into internal/app
- [x] 2.3 Bridge Ent store to Agent Runtime in App initialization

## 3. Tool Registration

- [x] 3.1 Update Config to include tools configuration (internal/config)
- [x] 3.2 Implement auto-registration of tools (Exec, Browser, FS) in App based on config
- [x] 3.3 Inject registered tools into Agent Runtime

## 4. Gateway Integration

- [x] 4.1 Update internal/gateway/server.go to accept Agent dependency
- [x] 4.2 Implement "chat.message" RPC handler delegating to Agent
- [x] 4.3 Wire Gateway into App life-cycle (Start/Stop)

## 5. Channel Wiring

- [x] 5.1 Implement ChannelAdapter interface for uniform channel management (Done via app.Channel interface and wiring)
- [x] 5.2 Wire Telegram channel to Agent in App
- [x] 5.3 Wire Discord channel to Agent in App
- [x] 5.4 Wire Slack channel to Agent in App

## 6. Main Entry Point

- [x] 6.1 Refactor cmd/lango/main.go to use internal/app
- [x] 6.2 Ensure graceful shutdown propagates to App.Stop()

## 7. Configuration & Testing

- [x] 7.1 Create default lango.json for local development (if missing)
- [ ] 7.2 Manual Test: Verify full flow (Start app -> Connect Telegram -> Chat with Agent)
