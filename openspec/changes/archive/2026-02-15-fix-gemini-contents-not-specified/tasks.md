## 1. AppendEvent In-Memory History Sync

- [x] 1.1 Update `SessionServiceAdapter.AppendEvent` in `internal/adk/session_service.go` to append message to `SessionAdapter.sess.History` after successful DB store
- [x] 1.2 Update `mockStore.AppendMessage` in `internal/adk/state_test.go` to use DB-only storage (separate messages map) instead of modifying `s.History` directly
- [x] 1.3 Add tests in `internal/adk/session_service_test.go`: single event in-memory update, multiple events accumulation, state-delta skip, DB+memory both updated

## 2. SystemInstruction Forwarding

- [x] 2.1 Add `extractSystemText` helper function in `internal/adk/model.go` to concatenate text parts from `genai.Content`
- [x] 2.2 Update `ModelAdapter.GenerateContent` in `internal/adk/model.go` to prepend system message from `req.Config.SystemInstruction`
- [x] 2.3 Add tests in `internal/adk/model_test.go`: SystemInstruction present, SystemInstruction absent

## 3. Verification

- [x] 3.1 Run `go build ./...` and verify no compilation errors
- [x] 3.2 Run `go test ./...` and verify all tests pass
