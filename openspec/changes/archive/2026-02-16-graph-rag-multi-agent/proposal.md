# Graph RAG + Sub Agents + A2A + Self-Learning Integration

## Summary
Evolve Lango from a monolithic single-agent architecture to a relationship-aware multi-agent system with Graph RAG (Cayley+BoltDB), Sub Agents, Agent Orchestration, Google A2A Protocol, Observable Memory Graph, and Self-Learning Graph.

## Motivation
- Current RAG (sqlite-vec cosine similarity) cannot traverse relationships (A→B→C path traversal)
- Knowledge system uses flat keyword matching, missing semantic relationships
- Memory is a linear list without temporal/causal relationship tracking
- Learning engine uses flat pattern matching without error-resolution graph
- Single monolithic agent cannot specialize for different task types
- No external agent communication protocol

## Approach
5-phase implementation:
1. **Phase 1**: Graph Store Foundation (Cayley + BoltDB)
2. **Phase 2**: Graph RAG Hybrid Retrieval
3. **Phase 3**: Sub-Agent Orchestration (ADK native)
4. **Phase 4**: A2A Protocol Integration
5. **Phase 5**: Observable Memory + Self-Learning Graph

## Impact
- New packages: `internal/graph/`, `internal/orchestration/`, `internal/a2a/`
- Modified: `internal/config/types.go`, `internal/app/wiring.go`, `internal/adk/agent.go`, `internal/adk/context_model.go`
- New dependencies: `github.com/cayleygraph/cayley`, `github.com/cayleygraph/quad`
