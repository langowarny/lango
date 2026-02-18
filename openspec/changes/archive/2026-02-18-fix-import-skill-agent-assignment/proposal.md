## Why

The `import_skill` tool was not assigned to any sub-agent in multi-agent orchestration mode. While `create_skill` and `list_skills` were correctly routed to the librarian agent via prefix matching, `import_skill` fell through to the "Unmatched" category because its prefix was missing from the librarian's prefix list. This meant the orchestrator could not reliably delegate skill import requests to the correct agent.

## What Changes

- Add `"import_skill"` to the librarian agent's `Prefixes` list in `agentSpecs` so it routes correctly via `PartitionTools`
- Add `"import_skill"` entry to `capabilityMap` so the routing table shows the correct capability description

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `agent-routing`: Add `import_skill` prefix to librarian agent's prefix list to ensure correct tool-to-agent routing

## Impact

- `internal/orchestration/tools.go`: Modified librarian `Prefixes` and `capabilityMap`
- Multi-agent orchestration mode: `import_skill` now routed to librarian instead of appearing as unmatched
- No breaking changes; single-agent mode unaffected
