## Context

The multi-agent orchestration system routes tools to sub-agents via prefix matching in `PartitionTools`. The librarian agent handles skill-related tools (`create_skill`, `list_skills`) but the `import_skill` tool was missing from its prefix list, causing it to fall into the "Unmatched" category.

## Goals / Non-Goals

**Goals:**
- Ensure `import_skill` is correctly routed to the librarian agent
- Add capability description for `import_skill` in the routing table

**Non-Goals:**
- Refactoring the prefix-based routing mechanism
- Adding new sub-agents or tools
- Changing import_skill tool behavior

## Decisions

**Add prefix to existing librarian spec rather than creating a new agent**
- Rationale: `import_skill` is semantically part of skill/knowledge management, which is the librarian's domain. The existing prefix list already has `create_skill` and `list_skills`.
- Alternative: Assign to operator (rejected — operator handles shell/file/exec, not knowledge domain tools)

**Add to capabilityMap for routing table visibility**
- Rationale: The orchestrator's routing table derives capability descriptions from `capabilityMap`. Without an entry, the tool would show as "general actions" instead of a meaningful description.

## Risks / Trade-offs

- [Minimal risk] The change is additive — only adds a new prefix entry to an existing list. No existing routing behavior is altered.
