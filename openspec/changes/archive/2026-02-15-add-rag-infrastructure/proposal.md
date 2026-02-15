# Proposal: RAG Infrastructure with sqlite-vec

## Problem

The agent's Knowledge, Memory (Observation/Reflection), and Learning searches all use keyword-based SQL LIKE matching. This fails to handle synonyms, multilingual queries, and contextual similarity. The agent cannot semantically associate related concepts across its knowledge base.

## Solution

Add vector embedding and similarity search to the existing SQLite database using sqlite-vec. Build a complete RAG (Retrieval-Augmented Generation) pipeline that:

1. **Embeds** all saved Knowledge, Observations, Reflections, and Learnings asynchronously via a background goroutine
2. **Stores** vector embeddings in sqlite-vec virtual tables alongside existing ent-managed data
3. **Retrieves** semantically relevant context before each LLM call and injects it into the system prompt

## Scope

- **In scope**: Embedding provider interfaces (OpenAI, Google, Ollama), sqlite-vec VectorStore, async embedding buffer, RAG service, agent integration, config, doctor check
- **Out of scope**: Backfill CLI command (deferred to follow-up), custom fine-tuned models, embedding model training

## Approach

- 3-provider support with registry pattern and fallback
- Async embedding via batched background goroutine (follows existing memory.Buffer pattern)
- Store callbacks avoid circular imports (embedding → knowledge/memory, not reverse)
- RAG results injected into ContextAwareModelAdapter alongside existing keyword-based retrieval
- Zero-config degradation: if no embedding provider is configured, system behaves exactly as before
