## Context

Lango is a Go-based AI agent platform using the Google ADK runtime, Ent ORM for persistence, and a gateway server for multi-channel communication. The agent currently operates statelessly — each session has message history but no cross-session knowledge accumulation. The agent cannot remember user preferences, learn from repeated errors, or reuse discovered workflows.

This design introduces a 6 Context Layer architecture that enables the agent to build and query a persistent knowledge base, automatically learn from tool execution outcomes, and create reusable skills.

## Goals / Non-Goals

**Goals:**
- Enable persistent knowledge storage (user rules, preferences, definitions, facts) across sessions.
- Automatically detect error patterns from tool execution and store diagnosed fixes with confidence scoring.
- Allow the agent to create reusable multi-step skills (composite, script, template).
- Augment LLM system prompts with relevant context retrieved via keyword-based RAG.
- Provide rate limiting to prevent runaway knowledge accumulation.
- Maintain backward compatibility — knowledge system is opt-in via configuration.

**Non-Goals:**
- Semantic/vector search (keyword-based search is sufficient for MVP).
- Multi-user knowledge isolation (single knowledge base per deployment).
- Skill marketplace or sharing between deployments.
- Real-time collaboration on knowledge entries.
- Changing the Ent/SQLite session storage schema.

## Decisions

### 1. 6 Context Layer Architecture
Context is organized into 6 layers: Tool Registry, User Knowledge, Skill Patterns, External Knowledge, Agent Learnings, Runtime Context.
- **Why**: Separating concerns allows targeted retrieval and different storage strategies per layer. The retriever can selectively query relevant layers based on the user's query.
- **Alternative**: Flat knowledge store with tags. **Rejected** because layer-based retrieval is more predictable and allows per-layer limits.

### 2. Keyword-Based RAG over Vector Search
The context retriever uses keyword extraction with stop-word filtering and SQL `CONTAINS` queries.
- **Why**: Keeps the system dependency-free (no embedding model or vector DB needed). For a per-deployment knowledge base with hundreds of entries, keyword search is sufficient.
- **Alternative**: pgvector/SQLite-vec with embedding models. **Rejected** for MVP complexity — can be added later as a retriever strategy.

### 3. Context-Aware Model Adapter Pattern
A `ContextAwareModelAdapter` wraps the existing ADK `ModelAdapter`, intercepting `GenerateContent` calls to augment the system prompt with retrieved context before forwarding to the LLM.
- **Why**: Non-invasive integration — the ADK agent and existing model adapter remain unchanged. Context injection happens transparently at the model boundary.
- **Alternative**: Modifying the ADK agent's prompt construction. **Rejected** because it would couple knowledge retrieval to the ADK internals.

### 4. Atomic Rate Limiting with Reserve Pattern
Rate limiting uses a `reserveSlot` pattern that atomically checks the limit and increments the counter within a single mutex lock.
- **Why**: Prevents TOCTOU race conditions where concurrent requests could both pass a check-then-increment pattern.
- **Alternative**: Database-level constraints. **Rejected** because per-session limits are transient (reset per process) and don't need persistence.

### 5. Skill Execution with Dangerous Pattern Validation
Script-type skills are validated against a blocklist of dangerous patterns (fork bombs, `rm -rf /`, pipe-to-shell) before execution.
- **Why**: Defense-in-depth layer. The Supervisor already provides sandboxing, but script validation adds an extra safety net for user-created skills.
- **Alternative**: Full sandboxing (containers, seccomp). **Rejected** for MVP — the existing Supervisor isolation combined with pattern validation is sufficient.

### 6. Learning Engine as Tool Wrapper
The learning engine wraps tool handlers, observing every execution result to detect error patterns and boost confidence on successes.
- **Why**: Passive observation requires no changes to tool implementations. The engine decorates handlers at wiring time.
- **Alternative**: Explicit learning calls in each tool. **Rejected** because it violates separation of concerns and requires modifying every tool.

## Risks / Trade-offs

- [x] **Risk**: Keyword search may miss semantically relevant but lexically different content.
  - **Mitigation**: Acceptable for MVP. The retriever can be swapped to vector search later without changing the interface.
- [x] **Risk**: Script skill execution could be bypassed via encoding tricks (base64, variable substitution).
  - **Mitigation**: Pattern validation is defense-in-depth only. Primary security comes from the Supervisor's sandboxing.
- [x] **Risk**: Large knowledge bases could slow down context retrieval.
  - **Mitigation**: Per-layer limits (default 5 items) and SQL indexes on searchable fields keep queries bounded.
- [x] **Risk**: Auto-learning from tool errors could accumulate low-quality entries.
  - **Mitigation**: Confidence scoring naturally deprioritizes unreliable learnings. Per-session rate limits cap accumulation.

## Migration Plan

1. **Schema Migration**: Ent auto-migration adds 5 new tables. No changes to existing session/message tables.
2. **Configuration**: New `knowledge` section with `enabled: false` default. Existing deployments are unaffected.
3. **Rollback**: Set `knowledge.enabled: false` to disable. Tables remain but are unused.
