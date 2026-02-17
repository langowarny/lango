## 1. Capability Description Layer

- [x] 1.1 Add `capabilityMap` variable mapping tool name prefixes to natural-language capability strings in `internal/orchestration/tools.go`
- [x] 1.2 Add `toolCapability(name)` function to return capability for a tool name prefix
- [x] 1.3 Add `capabilityDescription(tools)` function to build deduplicated capability string from tool list

## 2. Orchestrator Prompt Update

- [x] 2.1 Replace `toolNameList()` calls with `capabilityDescription()` in executor sub-agent description
- [x] 2.2 Replace `toolNameList()` calls with `capabilityDescription()` in researcher sub-agent description
- [x] 2.3 Replace `toolNameList()` calls with `capabilityDescription()` in memory-manager sub-agent description
- [x] 2.4 Replace tool name lists in `agentDescriptions` entries with capability descriptions

## 3. Hallucination Retry Fallback

- [x] 3.1 Extract `runAndCollectOnce()` from existing `RunAndCollect` logic in `internal/adk/agent.go`
- [x] 3.2 Add `extractMissingAgent(err)` helper to parse hallucinated agent name from error
- [x] 3.3 Add `subAgentNames(agent)` helper to collect valid sub-agent names
- [x] 3.4 Implement retry wrapper in `RunAndCollect`: detect hallucination error → correction message → retry once

## 4. Tests

- [x] 4.1 Add `TestToolCapability` — verify prefix-to-capability mapping for all known prefixes
- [x] 4.2 Add `TestCapabilityDescription` — verify deduplication and unknown tool fallback
- [x] 4.3 Add `TestBuildAgentTree_DescriptionsUseCapabilities` — verify no raw tool names in sub-agent descriptions
- [x] 4.4 Add `TestExtractMissingAgent` — verify error message parsing for hallucinated agent names
- [x] 4.5 Verify all existing `TestBuildAgentTree_*` tests pass with updated descriptions

## 5. Verification

- [x] 5.1 `go build ./...` compiles without errors
- [x] 5.2 `go test ./internal/orchestration/...` passes all tests
- [x] 5.3 `go test ./internal/adk/...` passes all tests
- [x] 5.4 `go test ./...` passes all tests
