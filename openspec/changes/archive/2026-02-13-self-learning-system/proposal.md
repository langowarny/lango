## Why

Lango currently has no memory across sessions — every conversation starts from scratch. The agent cannot remember user preferences, learn from past errors, or reuse discovered workflows. This limits effectiveness for repeated tasks and forces users to re-explain context every time. A self-learning knowledge system enables the agent to accumulate expertise over time, turning it from a stateless tool into a persistent assistant.

## What Changes

- **Knowledge Persistence**: New Ent schemas (Knowledge, Learning, Skill, AuditLog, ExternalRef) for storing structured knowledge, error patterns, reusable skills, audit trails, and external references.
- **6 Context Layer Architecture**: A retrieval-augmented generation (RAG) system that searches across 6 context layers (Tool Registry, User Knowledge, Skill Patterns, External Knowledge, Agent Learnings, Runtime Context) and injects relevant context into the system prompt.
- **Learning Engine**: Automatic error pattern detection and confidence-based learning that triggers on tool execution results.
- **Skill System**: Registry, executor, and builder for creating reusable multi-step workflows (composite), shell scripts, and Go text/templates.
- **Meta-Tools**: 6 new agent tools (`save_knowledge`, `search_knowledge`, `save_learning`, `search_learnings`, `create_skill`, `list_skills`) that allow the agent to self-manage its knowledge base.
- **Context-Aware Model Adapter**: Wraps the ADK model adapter to dynamically augment system prompts with retrieved context before each LLM call.
- **Configuration**: New `knowledge` config section with toggles for enabling/disabling, rate limits, and skill approval settings.

## Capabilities

### New Capabilities
- `knowledge-store`: Persistent CRUD operations for knowledge, learning, skill, audit log, and external reference entities with per-session rate limiting.
- `context-retriever`: RAG-based context retrieval across 6 layers with keyword extraction, stop-word filtering, and prompt assembly.
- `learning-engine`: Automatic error pattern extraction, categorization, and confidence-based learning from tool execution results.
- `skill-system`: Skill registry, executor (composite/script/template), and builder with dangerous pattern validation and lifecycle management.
- `meta-tools`: Agent-facing tools for self-managing the knowledge base (save, search, create, list operations).

### Modified Capabilities
- `application-core`: Wiring updated to initialize knowledge components (Store, Engine, Registry) and integrate context-aware model adapter.
- `config-system`: New `KnowledgeConfig` section added to root configuration with learning/knowledge limits and skill approval settings.

## Impact

- **Codebase**: New packages `internal/knowledge`, `internal/learning`, `internal/skill`. Modified `internal/app` (wiring, tools, types) and `internal/config` (types, loader).
- **Database**: 5 new Ent schemas with auto-migration (Knowledge, Learning, Skill, AuditLog, ExternalRef).
- **Dependencies**: No new external dependencies — uses existing Ent ORM and zap logger.
- **Configuration**: New optional `knowledge` section in `lango.json`/`lango.yaml`. Disabled by default for backward compatibility.
