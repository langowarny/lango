## Context

When the knowledge system is enabled in `app.go`, the tool initialization flow is:
1. `buildTools()` creates base tools (`exec`, `exec_bg`, filesystem tools, etc.)
2. `initKnowledge()` passes these base tools to `Registry`, which stores them as `baseTools`
3. Tools are wrapped with learning engine hooks → `tools = wrapped`
4. `kc.registry.AllTools()` is appended to `tools`

`AllTools()` returns `baseTools + loaded`, where `baseTools` still references the original unwrapped tools. This means `exec` (and other base tools) appear twice in the final tools slice—once wrapped and once original—causing the ADK runner to reject them with a "duplicate tool" error.

## Goals / Non-Goals

**Goals:**
- Eliminate duplicate tool names when knowledge system appends skill tools
- Maintain backward compatibility of `AllTools()` for any other callers
- Provide a clean API to retrieve only dynamically loaded skills

**Non-Goals:**
- Refactoring the overall tool initialization pipeline in `app.go`
- Changing how `Registry` stores or uses `baseTools` internally

## Decisions

**Add `LoadedSkills()` method instead of modifying `AllTools()`**

Rationale: `AllTools()` has a valid use case (returning the full tool set including base tools) and may be used elsewhere. Adding a separate `LoadedSkills()` method is additive and non-breaking. The call site in `app.go` switches to `LoadedSkills()` since it already has the base tools in the `tools` slice.

Alternative considered: Removing `baseTools` from `Registry` entirely. Rejected because the registry uses `baseTools` internally and `AllTools()` serves other callers that need the complete set.

## Risks / Trade-offs

- [Minimal risk] `AllTools()` remains unchanged → no impact on other callers
- [Trade-off] Two methods with overlapping concerns (`AllTools` vs `LoadedSkills`) → acceptable since semantics are clear and both are documented
