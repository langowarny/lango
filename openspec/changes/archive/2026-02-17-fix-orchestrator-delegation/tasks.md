## 1. Remove Direct Tools from Orchestrator

- [x] 1.1 Delete `allAdkTools` adaptation code (lines 133-138) in `internal/orchestration/orchestrator.go`
- [x] 1.2 Set orchestrator `Tools: nil` in `llmagent.Config`

## 2. Rewrite Orchestrator Instruction

- [x] 2.1 Replace instruction to state orchestrator has no tools and must delegate
- [x] 2.2 Add delegation rules per sub-agent role (executor, researcher, planner, memory-manager)
- [x] 2.3 Keep direct response for simple conversational messages

## 3. Improve Sub-Agent Instructions

- [x] 3.1 Add result-reporting guidance to executor instruction
- [x] 3.2 Add result-reporting guidance to researcher instruction
- [x] 3.3 Add result-reporting guidance to planner instruction
- [x] 3.4 Add result-reporting guidance to memory-manager instruction

## 4. Update MaxDelegationRounds Default

- [x] 4.1 Change default from 3 to 5 in `internal/orchestration/orchestrator.go`
- [x] 4.2 Update `MaxDelegationRounds: 5` in `internal/app/wiring.go`

## 5. Update Tests

- [x] 5.1 Rename `TestBuildAgentTree_OrchestratorHasToolsAndSubAgents` to verify orchestrator has NO direct tools
- [x] 5.2 Remove `TestBuildAgentTree_OrchestratorAdaptError` (no longer applicable)
- [x] 5.3 Verify each tool is adapted exactly once (only for sub-agent)

## 6. Strip Tool Sections from Orchestrator System Prompt

- [x] 6.1 Build orchestrator-specific prompt in `internal/app/wiring.go` multi-agent block
- [x] 6.2 Remove `SectionToolUsage` from orchestrator prompt builder
- [x] 6.3 Replace `SectionIdentity` with orchestrator-specific identity (no tool categories)
- [x] 6.4 Pass `orchestratorPrompt` to `orchestration.Config.SystemPrompt`

## 7. Verification

- [x] 7.1 Run `go build ./...` — build passes
- [x] 7.2 Run `go test ./internal/orchestration/...` — all tests pass
- [x] 7.3 Run `go test ./internal/app/...` — all tests pass
- [x] 7.4 Run `go test ./...` — full test suite passes
