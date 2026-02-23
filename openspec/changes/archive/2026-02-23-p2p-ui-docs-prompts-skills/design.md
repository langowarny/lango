## Context

This is a documentation-only change. The P2P REST API endpoints, `--value-hex` flag, and Docker Compose example were implemented in the `p2p-trading-docker-compose` change. This change updates all surrounding documentation, prompts, CLI help text, and docs site to reflect those features.

## Goals / Non-Goals

**Goals:**
- Ensure all user-facing documentation reflects P2P REST API, `--value-hex`, and Docker example
- Keep prompts and agent identity files accurate for multi-agent orchestration
- Maintain consistency between code behavior and documentation

**Non-Goals:**
- No code changes beyond CLI help text (Long descriptions)
- No new features or behavior changes

## Decisions

### 1. Documentation-only, no delta specs

**Rationale**: No spec-level requirements changed. The underlying behavior was already specified in `p2p-trading-docker-compose`. This change only updates docs, prompts, and CLI help text.

### 2. REST API note in P2P CLI Long descriptions

**Rationale**: P2P CLI commands create ephemeral nodes separate from the running server. Adding a Long description clarifies this distinction and points users to the REST API for querying the running server's node.

### 3. Examples section in README

**Rationale**: The `examples/p2p-trading/` directory is a significant addition but was not referenced from the main README. Adding an Examples section makes it discoverable.

## Risks / Trade-offs

- Minimal risk — documentation changes only
- CLI Long descriptions add minor verbosity to `--help` output
