## Why

The current Lango service restricts the agent to a read-only and synchronous interaction model for system tools. Specifically:
- **Filesystem**: The agent can only read files (`fs_read`), preventing it from making code changes, creating new files, or even exploring directory structures.
- **Execution**: The agent can only run synchronous commands (`exec`), which blocks the entire conversation loop for long-running processes (like starting a server) and prevents the agent from monitoring background tasks.

Unlocking these capabilities is essential for turning Lango into a truly autonomous agent capable of end-to-end development and system management.

## What Changes

- **Filesystem Tool Expansion**:
    - Register `fs_list` to allow directory exploration.
    - Register `fs_write` and `fs_edit` to enable the agent to modify the codebase.
    - Register `fs_mkdir` and `fs_delete` for full lifecycle management of files and directories.
- **Execution Tool Expansion**:
    - Register `exec_bg` to allow starting background processes.
    - Implement status check and stop handlers for background processes to ensure proper lifecycle management.
- **Codebase Clean-up**: Remove redundant logic in the existing `exec` handler in `internal/app/app.go`.

## Capabilities

### New Capabilities
- `fs-write`: Full write/edit access to the filesystem, enabling autonomous coding and refactoring.
- `fs-mgmt`: Ability to list, create, and delete directories and files, facilitating project structure management.
- `exec-background`: Support for long-running tasks and non-blocking command execution, allowing the agent to manage services.

### Modified Capabilities
- `tool-exec`: Enhance the existing execution interface to support non-blocking operations and better error reporting.

## Impact

- `internal/app/app.go`: Registration of new tool handlers and cleanup of the `exec` handler.
- `internal/app/types.go`: Addition of state management for background processes if needed (similar to `BrowserSessionID`).
- `internal/tools/exec/exec.go`: Potential small fixes to ensure background process monitoring is robust.
- `internal/tools/filesystem/filesystem.go`: Verification of existing methods before exposure.
