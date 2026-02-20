## MODIFIED Requirements

### Requirement: Workflow context isolation
The workflow engine SHALL detach the execution context from the originating request context before running DAG steps. The detached context SHALL preserve context values but SHALL NOT be cancelled when the parent request completes.

#### Scenario: Parent request completes without cancelling workflow
- **WHEN** a workflow is started via `workflow_run` and the originating agent request completes
- **THEN** the workflow SHALL continue executing all DAG steps without `context canceled` errors

#### Scenario: Workflow cancellation via workflow_cancel still works
- **WHEN** a user calls `workflow_cancel` on a running workflow
- **THEN** the workflow SHALL be cancelled immediately via the stored cancel function

### Requirement: Async workflow execution
The workflow engine SHALL provide a `RunAsync()` method that validates the workflow, creates run and step records, then executes the DAG in a background goroutine. The method SHALL return the run ID immediately.

#### Scenario: RunAsync returns immediately with run ID
- **WHEN** `RunAsync()` is called with a valid workflow
- **THEN** the method SHALL return a run ID and `nil` error before DAG execution begins

#### Scenario: RunAsync workflow progress is queryable
- **WHEN** a workflow is started via `RunAsync()` and the run ID is used with `Status()`
- **THEN** the status SHALL reflect the current execution state (running, completed, or failed)

#### Scenario: workflow_run tool uses RunAsync
- **WHEN** the `workflow_run` tool handler is invoked
- **THEN** it SHALL call `RunAsync()` and return the run ID immediately with status "running"
