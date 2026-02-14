## 1. Type Definitions

- [x] 1.1 Add `ToolDescriptor` struct to `internal/knowledge/types.go` with Name and Description fields
- [x] 1.2 Add `RuntimeContext` struct to `internal/knowledge/types.go` with SessionKey, ChannelType, ActiveToolCount, EncryptionEnabled, KnowledgeEnabled, MemoryEnabled fields

## 2. Retriever Interfaces and Builder Methods

- [x] 2.1 Define `ToolRegistryProvider` interface in `internal/knowledge/retriever.go` with `ListTools()` and `SearchTools(query, limit)`
- [x] 2.2 Define `RuntimeContextProvider` interface in `internal/knowledge/retriever.go` with `GetRuntimeContext()`
- [x] 2.3 Add `toolProvider` and `runtimeProvider` optional fields to `ContextRetriever`
- [x] 2.4 Add `WithToolRegistry()` and `WithRuntimeContext()` builder methods

## 3. Retrieval Methods

- [x] 3.1 Implement `retrieveTools()` method — delegate to provider's SearchTools, map to ContextItem slice
- [x] 3.2 Implement `retrieveRuntimeContext()` method — delegate to provider's GetRuntimeContext, format as single ContextItem
- [x] 3.3 Replace `continue // handled elsewhere` in Retrieve() switch with calls to retrieveTools and retrieveRuntimeContext

## 4. Prompt Assembly Update

- [x] 4.1 Add "Runtime Context" section to AssemblePrompt before existing sections
- [x] 4.2 Add "Available Tools" section to AssemblePrompt after Runtime Context, before User Knowledge

## 5. Adapter Implementations

- [x] 5.1 Create `internal/adk/context_providers.go` with `ToolRegistryAdapter` — boundary copy on construction, case-insensitive substring search
- [x] 5.2 Implement `RuntimeContextAdapter` in same file — NewRuntimeContextAdapter, SetSession with deriveChannelType, GetRuntimeContext with sync.RWMutex

## 6. Context Model Update

- [x] 6.1 Add `runtimeAdapter` field and `WithRuntimeAdapter()` builder to `ContextAwareModelAdapter`
- [x] 6.2 Update `GenerateContent()` to call `runtimeAdapter.SetSession()` before retrieval
- [x] 6.3 Update `GenerateContent()` to request all 6 layers explicitly

## 7. Wiring

- [x] 7.1 Create `ToolRegistryAdapter` in `initAgent()` and wire to retriever via `WithToolRegistry()`
- [x] 7.2 Create `RuntimeContextAdapter` in `initAgent()` and wire to retriever via `WithRuntimeContext()` and to ctxAdapter via `WithRuntimeAdapter()`

## 8. Tests

- [x] 8.1 Add `TestContextRetriever_RetrieveTools` to `internal/knowledge/retriever_test.go` — mock provider, keyword search, nil provider
- [x] 8.2 Add `TestContextRetriever_RetrieveRuntimeContext` to `internal/knowledge/retriever_test.go` — mock provider, session info, nil provider
- [x] 8.3 Add `TestAssemblePrompt_WithToolsAndRuntime` to `internal/knowledge/retriever_test.go` — section ordering verification
- [x] 8.4 Create `internal/adk/context_providers_test.go` with ListTools, SearchTools (table-driven), BoundaryCopy tests
- [x] 8.5 Add RuntimeContextAdapter and DeriveChannelType tests to `internal/adk/context_providers_test.go`
- [x] 8.6 Add compile-time interface compliance checks for both adapters

## 9. Verification

- [x] 9.1 Run `go build ./...` — confirm clean compilation
- [x] 9.2 Run `go test ./internal/knowledge/... -race` — all pass
- [x] 9.3 Run `go test ./internal/adk/... -race` — all pass
- [x] 9.4 Run `go vet ./...` — no issues
