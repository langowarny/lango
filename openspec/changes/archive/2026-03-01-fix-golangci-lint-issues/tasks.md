## 1. Configuration

- [x] 1.1 Create `.golangci.yml` with v2 format, `default: standard`, `generated: strict`, `std-error-handling` preset, test errcheck exclusion

## 2. errcheck Fixes

- [x] 2.1 Fix `defer resp.Body.Close()` in `internal/agent/pii_presidio.go` (2 locations)
- [x] 2.2 Add `writeJSON` helper in `internal/app/p2p_routes.go` and replace 12 `json.NewEncoder(w).Encode(...)` calls
- [x] 2.3 Fix `sendError()` unchecked errors in discord, slack, telegram channel adapters
- [x] 2.4 Fix `tx.Rollback()` in `internal/session/ent_store.go` and `internal/embedding/sqlite_vec.go`
- [x] 2.5 Fix `Process.Signal/Kill` in `internal/tools/exec/exec.go`
- [x] 2.6 Fix `os.Rename` rollback in `internal/dbmigrate/migrate.go`
- [x] 2.7 Fix `defer logging.Sync()`, `defer resp.Body.Close()`, `fmt.Scanln` in `cmd/lango/main.go`
- [x] 2.8 Fix remaining errcheck: `app.go`, `cli/p2p/p2p.go`, `cli/payment/send.go`, `cli/security/secrets.go`, `gateway/auth.go`, `gateway/server.go`

## 3. staticcheck Fixes

- [x] 3.1 Fix QF1012 (WriteString+Sprintf→Fprintf) across all files: `adk/context_model.go`, `cron/delivery.go`, `knowledge/retriever.go`, `graph/rag.go`, `workflow/engine.go`, `orchestration/tools.go`, `librarian/inquiry_processor.go`, `librarian/observation_analyzer.go`, `skill/parser.go`
- [x] 3.2 Fix S1009 (redundant nil check) in `skill/parser.go`, `skill/registry.go`
- [x] 3.3 Fix S1011 (append spread) in `workflow/parser.go`
- [x] 3.4 Fix SA1012 (nil context) in `x402/handler.go`
- [x] 3.5 Fix SA9003 (empty branches) in `adk/model.go`, `gateway/middleware_test.go`, `payment/service.go`
- [x] 3.6 Fix QF1003 (if/else→switch) in `adk/model.go`, `provider/anthropic/anthropic.go`
- [x] 3.7 Fix QF1008 (embedded field selector) in `adk/session_service.go`, `adk/state_test.go`
- [x] 3.8 Fix S1017 (redundant HasSuffix) in `learning/parse.go`, `librarian/parse.go`
- [x] 3.9 Fix ST1023 (type in declaration) in `app/wiring.go`
- [x] 3.10 Fix S1000 (select single case) in `channels/slack/slack_test.go`
- [x] 3.11 Fix ST1005 (error string case) in `security/azure_kv_provider_stub.go`
- [x] 3.12 Fix SA4006 (unused assignment) in `cli/settings/auth_providers_list.go`, `cli/settings/providers_list.go`
- [x] 3.13 Fix QF1002 (switch rewrite) in `skill/importer_test.go`
- [x] 3.14 Add `//nolint:staticcheck` for SA1019 deprecated field usage in `cli/settings/forms_impl_test.go`, `cli/tuicore/state_update.go`, `p2p/node.go`

## 4. Unused Code Removal

- [x] 4.1 Remove `executor` field and `adka2a` import from `internal/a2a/server.go`
- [x] 4.2 Remove `fakeAgent` type from `internal/a2a/server_test.go`
- [x] 4.3 Remove `wrapWithLearning()`, `wrapWithApproval()` from `internal/app/tools.go`
- [x] 4.4 Remove `wg sync.WaitGroup` field from `internal/channels/discord/discord.go`
- [x] 4.5 Remove `logger` variable from `internal/cli/security/migrate.go`
- [x] 4.6 Remove `toGraphTriples()` from `internal/librarian/proactive_buffer.go`
- [x] 4.7 Remove `logger` variable from `internal/session/ent_store.go`

## 5. ineffassign Fix

- [x] 5.1 Remove dead `message` assignment in `internal/cli/doctor/checks/security.go`

## 6. Verification

- [x] 6.1 Verify `go build ./...` passes
- [x] 6.2 Verify `go test ./...` passes
- [x] 6.3 Verify `golangci-lint run` reports 0 issues
