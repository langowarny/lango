# Tasks

## Phase 1: Graph Store Foundation
- [x] Add GraphConfig, A2AConfig, MultiAgent config types to `internal/config/types.go`
- [x] Create `internal/graph/store.go` — Store interface, Triple type, predicate constants
- [x] Create `internal/graph/bolt_store.go` — BoltDB-backed implementation with SPO/POS/OSP indexes
- [x] Create `internal/graph/buffer.go` — Async graph update buffer (Start/Enqueue/Stop)
- [x] Create `internal/graph/bolt_store_test.go` — Table-driven tests for CRUD and traversal
- [x] Add `go.etcd.io/bbolt` dependency

## Phase 2: Graph RAG Hybrid Retrieval
- [x] Create `internal/graph/rag.go` — GraphRAGService with 2-phase hybrid retrieval
- [x] Create `internal/graph/extractor.go` — LLM-based entity/relationship extraction
- [x] Update `internal/adk/context_model.go` — WithGraphRAG(), assembleGraphRAGSection()
- [x] Add SetGraphCallback to `internal/knowledge/store.go`
- [x] Add SetGraphCallback to `internal/memory/store.go`

## Phase 3: Sub-Agent Orchestration
- [x] Create `internal/orchestration/tools.go` — Tool partitioning by agent role
- [x] Create `internal/orchestration/orchestrator.go` — BuildAgentTree with 4 sub-agents
- [x] Create `internal/orchestration/orchestrator_test.go` — Tests
- [x] Add `NewAgentFromADK()` to `internal/adk/agent.go`

## Phase 4: A2A Protocol
- [x] Create `internal/a2a/server.go` — A2A server with Agent Card at /.well-known/agent.json
- [x] Create `internal/a2a/remote.go` — Remote A2A agent loader via remoteagent.NewA2A()
- [x] Add `github.com/a2aproject/a2a-go` dependency (via ADK)

## Phase 5: Observable Memory + Self-Learning Graph
- [x] Create `internal/memory/graph_hooks.go` — GraphHooks for observation/reflection triples
- [x] Create `internal/learning/graph_engine.go` — GraphEngine with confidence propagation

## Phase 6: Wiring Integration
- [x] Update `internal/app/wiring.go` — initGraphStore, wireGraphCallbacks, initGraphRAG
- [x] Update `internal/app/wiring.go` — Multi-agent branch in initAgent
- [x] Update `internal/app/app.go` — Wire gc parameter
- [x] Full build: `go build ./...` — PASS
- [x] Full test: `go test ./...` — ALL PASS

## Phase 7: Full Integration — Gap Closure (8 gaps)

### Gap 1: GraphBuffer Lifecycle
- [x] Add `GraphStore`, `GraphBuffer` fields to App struct (`types.go`)
- [x] Assign gc.store/gc.buffer in `New()` (`app.go`)
- [x] Add `GraphBuffer.Start(&a.wg)` in `Start()` (`app.go`)
- [x] Add `GraphBuffer.Stop()` in `Stop()` (`app.go`)
- [x] Add `GraphStore.Close()` in `Stop()` (`app.go`)

### Gap 2: Entity Extractor Pipeline
- [x] Create Extractor in `wireGraphCallbacks` using `providerTextGenerator` (`wiring.go`)
- [x] Add async entity extraction goroutine in graphCB callback (`wiring.go`)

### Gap 3: GraphEngine Connection
- [x] Add `ToolResultObserver` interface to `learning/engine.go`
- [x] Add compile-time checks for Engine and GraphEngine
- [x] Add `observer` field to `knowledgeComponents` (`wiring.go`)
- [x] Update `initKnowledge` to accept `gc *graphComponents` and create GraphEngine when graph available
- [x] Change `wrapWithLearning` to accept `ToolResultObserver` interface (`tools.go`)
- [x] Reorder initialization: gc before kc (`app.go`)

### Gap 4: GraphHooks → Memory Store
- [x] Add `graphHooks`, `lastObsMu`, `lastObsIDs` to memory Store (`memory/store.go`)
- [x] Add `SetGraphHooks()` method (`memory/store.go`)
- [x] Call `graphHooks.OnObservation()` in `SaveObservation` with previous obs tracking
- [x] Call `graphHooks.OnReflection()` in `SaveReflection` with session observation IDs
- [x] Wire GraphHooks in `wireGraphCallbacks` (`wiring.go`)

### Gap 5: Researcher/MemoryManager Tools
- [x] Add `"save_knowledge"`, `"save_learning"` to researcher prefixes (`orchestration/tools.go`)
- [x] Create `buildGraphTools` — graph_traverse, graph_query (`app/tools.go`)
- [x] Create `buildRAGTools` — rag_retrieve (`app/tools.go`)
- [x] Create `buildMemoryAgentTools` — memory_list_observations, memory_list_reflections (`app/tools.go`)
- [x] Add tool builders to init flow in `New()` (`app.go`)

### Gap 6: A2A Server Routes
- [x] Add `Router() chi.Router` method to gateway Server (`gateway/server.go`)
- [x] Add `adkAgent` field + `ADKAgent()` getter to adk Agent (`adk/agent.go`)
- [x] Mount A2A routes after agent creation in `New()` (`app.go`)

### Gap 7: Remote A2A Agents
- [x] Add `RemoteAgents []adk_agent.Agent` to orchestration Config (`orchestration/orchestrator.go`)
- [x] Append RemoteAgents to sub-agents in `BuildAgentTree` (`orchestration/orchestrator.go`)
- [x] Replace `_ = remoteAgents` with `orchCfg.RemoteAgents = remoteAgents` (`wiring.go`)

### Gap 8: Tests
- [x] Create `internal/a2a/server_test.go` — Agent Card HTTP tests
- [x] Create `internal/learning/graph_engine_test.go` — RecordFix, callback, sanitize tests
- [x] Create `internal/memory/graph_hooks_test.go` — OnObservation/OnReflection triple tests

### Verification
- [x] `go build ./...` — PASS
- [x] `go test ./...` — ALL PASS
- [x] `go vet ./...` — PASS
