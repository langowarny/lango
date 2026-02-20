## MODIFIED Requirements

### Requirement: Background task context isolation
The background task manager SHALL detach the task context from the originating request context. The task context SHALL preserve context values (session key, approval target) but SHALL NOT be cancelled when the parent request completes. The task context SHALL have a configurable timeout (default: 30 minutes).

#### Scenario: Parent request completes without cancelling background task
- **WHEN** a background task is submitted via `bg_submit` and the originating agent request completes
- **THEN** the background task SHALL continue executing without `context canceled` errors

#### Scenario: Background task respects its own timeout
- **WHEN** a background task exceeds the configured `taskTimeout`
- **THEN** the task context SHALL be cancelled with `DeadlineExceeded`

#### Scenario: Manual cancellation still works
- **WHEN** a user calls `bg_cancel` on a running background task
- **THEN** the task SHALL be cancelled immediately via the task's cancel function

#### Scenario: Context values are preserved
- **WHEN** a background task is submitted from a session with approval target set
- **THEN** the background task's context SHALL carry the same approval target value

### Requirement: Configurable task timeout
The `BackgroundConfig` SHALL include a `TaskTimeout` field of type `time.Duration`. When `TaskTimeout` is zero or negative, the default of 30 minutes SHALL be used.

#### Scenario: Custom task timeout from configuration
- **WHEN** `background.taskTimeout` is set to `1h` in configuration
- **THEN** background tasks SHALL use a 1-hour timeout instead of the 30-minute default
