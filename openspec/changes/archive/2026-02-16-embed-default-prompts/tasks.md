## 1. Create prompts package

- [x] 1.1 Create `prompts/AGENTS.md` with production-quality identity prompt covering 5 tools, knowledge system, memory, multi-channel support, and response principles
- [x] 1.2 Create `prompts/SAFETY.md` with security rules covering secret protection, destructive operation confirmation, PII, path traversal, crypto material, env var filtering, and ambiguous request handling
- [x] 1.3 Create `prompts/CONVERSATION_RULES.md` with conversation behavior rules preventing answer repetition, respecting channel limits, and maintaining consistency
- [x] 1.4 Create `prompts/TOOL_USAGE.md` with per-tool usage guidelines for exec, filesystem, browser, crypto, and secrets including error handling
- [x] 1.5 Create `prompts/embed.go` with `//go:embed *.md` directive exposing `FS embed.FS`
- [x] 1.6 Create `prompts/embed_test.go` verifying all 4 `.md` files are embedded and non-empty

## 2. Update defaults to use embedded FS

- [x] 2.1 Remove hardcoded `const` strings from `internal/prompt/defaults.go`
- [x] 2.2 Add `defaultContent(filename, fallback string) string` helper that reads from `prompts.FS`
- [x] 2.3 Update `DefaultBuilder()` to call `defaultContent()` for each section
- [x] 2.4 Update `internal/prompt/defaults_test.go` assertions to match new embedded content
- [x] 2.5 Update `internal/prompt/loader_test.go` assertion for default conversation rules text

## 3. Docker and verification

- [x] 3.1 Add `COPY --from=builder /app/prompts/ /usr/share/lango/prompts/` to Dockerfile
- [x] 3.2 Verify `go build ./...` succeeds
- [x] 3.3 Verify `go test ./prompts/...` passes (embed verification)
- [x] 3.4 Verify `go test ./internal/prompt/...` passes (all prompt tests)
- [x] 3.5 Verify `go test ./...` passes (full project)
