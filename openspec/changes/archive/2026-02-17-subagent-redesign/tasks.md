## 1. AgentSpec Type and Registry

- [x] 1.1 Define AgentSpec struct with Name, Description, Instruction, Prefixes, Keywords, Accepts, Returns, CannotDo, AlwaysInclude fields
- [x] 1.2 Create agentSpecs registry with 6 specs: operator, navigator, vault, librarian, planner, chronicler
- [x] 1.3 Write structured Instructions for each spec with What You Do, Input Format, Output Format, Constraints sections
- [x] 1.4 Add [REJECT] protocol to each spec's Constraints section

## 2. RoleToolSet and Partitioning

- [x] 2.1 Update RoleToolSet struct: Operator, Navigator, Vault, Librarian, Planner, Chronicler, Unmatched fields
- [x] 2.2 Rewrite PartitionTools with matching order Librarian → Chronicler → Navigator → Vault → Operator → Unmatched
- [x] 2.3 Add specPrefixes helper to derive prefixes from agentSpecs
- [x] 2.4 Add toolsForSpec helper to map spec name to RoleToolSet field
- [x] 2.5 Update capabilityMap with new prefixes: secrets_, create_skill, list_skills

## 3. Routing Table and Orchestrator Prompt

- [x] 3.1 Define routingEntry type with Name, Description, Keywords, Accepts, Returns, CannotDo
- [x] 3.2 Implement buildRoutingEntry to create entries from AgentSpec + capabilities
- [x] 3.3 Implement buildOrchestratorInstruction with routing table, decision protocol, reject handling, unmatched tools section

## 4. BuildAgentTree Rewrite

- [x] 4.1 Replace hardcoded agent blocks with data-driven loop over agentSpecs
- [x] 4.2 Add conditional creation: skip agents with no tools unless AlwaysInclude
- [x] 4.3 Wire routing entries into orchestrator instruction via buildOrchestratorInstruction
- [x] 4.4 Preserve remote A2A agent integration

## 5. Tests

- [x] 5.1 TestPartitionTools — 6-role + Unmatched partitioning with 8 subtests
- [x] 5.2 TestPartitionTools_PrefixPriority — librarian priority over operator
- [x] 5.3 TestBuildAgentTree_Success — 6 sub-agents created
- [x] 5.4 TestBuildAgentTree_NoTools — only planner
- [x] 5.5 TestBuildAgentTree_PartialAgents — conditional creation
- [x] 5.6 TestBuildAgentTree_UnmatchedToolsNotAssigned — unmatched not adapted
- [x] 5.7 TestBuildAgentTree_RoutingTableInInstruction — routing table and decision protocol
- [x] 5.8 TestBuildAgentTree_RejectProtocolInInstructions — reject in all specs
- [x] 5.9 TestBuildRoutingEntry — routing entry builder
- [x] 5.10 TestBuildOrchestratorInstruction — routing table, unmatched tools
- [x] 5.11 TestAgentSpecs consistency tests — unique names, keywords, I/O metadata, instruction structure

## 6. Verification

- [x] 6.1 go build ./... passes
- [x] 6.2 go test ./internal/orchestration/... passes (25 tests)
- [x] 6.3 go test ./... passes (full project)
