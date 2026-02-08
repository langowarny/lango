## Context

The Lango agent is currently restricted by a subset of tool exposures in `internal/app/app.go`. While the underlying `internal/tools` packages are feature-rich, the agent lacks the ability to explore the directory structure, modify files, or run long-running background processes. This prevents the agent from performing complex tasks like codebase refactoring or service management.

## Goals / Non-Goals

**Goals:**
- Expose all essential filesystem operations (`list`, `write`, `edit`, `mkdir`, `delete`).
- Enable asynchronous command execution with lifecycle management (`status`, `stop`).
- Ensure all new tools respect existing security boundaries (path validation and environment whitelisting).
- Refactor redundant logic in the current `exec` tool handler.

**Non-Goals:**
- implementing new core logic in the `internal/tools` packages (mostly existing logic will be reused).
- Changes to the `agent` runtime itself.

## Decisions

### 1. Tool Registration in `app.go`
We will expose new handlers in `internal/app/app.go` that map directly to the existing methods in `internal/tools`.
- `fs_list` -> `fsTool.ListDir`
- `fs_write` -> `fsTool.Write`
- `fs_edit` -> `fsTool.Edit`
- `fs_mkdir` -> `fsTool.Mkdir`
- `fs_delete` -> `fsTool.Delete`
- `exec_bg` -> `execTool.StartBackground`

### 2. Execution Tool Proxy via Supervisor
The `exec` tool is mediated by the `Supervisor` to enforce security (Environment Whitelist). We will extend the `Supervisor` or directly utilize the `execTool` within `app.go` if appropriate, ensuring the `EnvWhitelist` is consistently applied. 

### 3. Background Process Management
The `exec.Tool` already maintains an internal map of background processes. We will add a new `exec_status` tool to check the status of these processes and an `exec_stop` tool to terminate them, using the IDs returned by `exec_bg`.

### 4. Input Validation
For `fs_edit`, we will ensure the agent provides clear line ranges to prevent accidental mass-deletion or malformed file states.

## Risks / Trade-offs

- **Security Risk**: Exposing `fs_write` and `fs_delete` increases the risk of the agent accidentally (or maliciously, if compromised) damaging the host environment. This is mitigated by the existing `AllowedPaths` configuration in Lango.
- **Resource Management**: Background processes could leak if the agent starts many without stopping them. We rely on the `exec.Tool.Cleanup()` method which is called on application shutdown, but more granular agent-side management is required for long-running sessions.
- **Complexity**: Adding 6-8 new tools increases the prompt context size. However, the default system prompt already guides the agent toward tool usage.
