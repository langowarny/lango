## Purpose

Define the context detach utility that creates a context independent of the parent's cancellation while preserving context values, enabling long-running goroutines to survive their originating request.

## Requirements

### Requirement: Detach context from parent cancellation
The `ctxutil.Detach()` function SHALL return a new context whose `Done()`, `Err()`, and `Deadline()` are independent of the parent context. The returned context SHALL delegate `Value()` calls to the original parent context.

#### Scenario: Parent cancellation does not affect detached context
- **WHEN** a parent context is cancelled after `Detach()` is called
- **THEN** the detached context's `Err()` SHALL return `nil` and `Done()` SHALL return `nil`

#### Scenario: Values are preserved through detached context
- **WHEN** a parent context carries values (e.g., session key, approval target)
- **THEN** the detached context SHALL return the same values via `Value()`

#### Scenario: Detached context has no deadline
- **WHEN** the parent context has a deadline or timeout
- **THEN** the detached context's `Deadline()` SHALL return `(zero, false)`

### Requirement: Detached context supports child wrapping
A detached context SHALL work correctly as a parent for `context.WithCancel()`, `context.WithTimeout()`, and `context.WithValue()`.

#### Scenario: WithCancel on detached context is independent
- **WHEN** `context.WithCancel()` wraps a detached context and the original parent is cancelled
- **THEN** the child context SHALL NOT be cancelled

#### Scenario: WithTimeout on detached context works independently
- **WHEN** `context.WithTimeout()` wraps a detached context with a 50ms timeout
- **THEN** the child context SHALL expire after 50ms with `DeadlineExceeded`, regardless of parent state
