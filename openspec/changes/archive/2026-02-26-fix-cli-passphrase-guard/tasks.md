## 1. Runtime Guard

- [x] 1.1 Expand `blockLangoExec()` in `internal/app/tools.go` with guards for graph, memory, p2p, security, and payment subcommands (with in-process tool alternatives)
- [x] 1.2 Add catch-all guard for any `lango` prefix not matched by specific subcommands
- [x] 1.3 Add comprehensive test cases in `internal/app/tools_test.go` covering all subcommands, catch-all, case-insensitivity, and non-lango commands

## 2. Prompt Reinforcement

- [x] 2.1 Add exec safety rule as first bullet in `### Exec Tool` section of `prompts/TOOL_USAGE.md`
- [x] 2.2 Broaden automation prompt warning in `internal/app/wiring.go` to cover all `lango` subcommands

## 3. Verification

- [x] 3.1 Run `go build ./...` to verify no compilation errors
- [x] 3.2 Run `go test ./internal/app/...` to verify all tests pass
