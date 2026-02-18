## ADDED Requirements

### Requirement: Panic recovery for rod/CDP calls
The browser tool SHALL recover from panics in go-rod/rod library calls and convert them into structured errors instead of crashing the process.

#### Scenario: Rod panic during navigation
- **WHEN** a rod API call panics during `Navigate`
- **THEN** the system SHALL recover the panic and return an error wrapping `ErrBrowserPanic`
- **AND** the process SHALL NOT crash

#### Scenario: Rod panic during screenshot
- **WHEN** a rod API call panics during `Screenshot`
- **THEN** the system SHALL recover the panic and return an error wrapping `ErrBrowserPanic`

#### Scenario: Rod panic during element interaction
- **WHEN** a rod API call panics during `Click`, `Type`, `GetText`, `GetElementInfo`, or `Eval`
- **THEN** the system SHALL recover the panic and return an error wrapping `ErrBrowserPanic`

#### Scenario: Rod panic during session creation
- **WHEN** a rod API call panics during `NewSession`
- **THEN** the system SHALL recover the panic and return an error wrapping `ErrBrowserPanic`

#### Scenario: Rod panic during close
- **WHEN** a rod API call panics during `Close`
- **THEN** the system SHALL recover the panic silently
- **AND** cleanup SHALL continue for remaining sessions

#### Scenario: Normal errors pass through unchanged
- **WHEN** a rod API call returns a normal error (no panic)
- **THEN** the error SHALL be returned as-is without `ErrBrowserPanic` wrapping

### Requirement: Auto-reconnect on browser panic
The SessionManager SHALL detect `ErrBrowserPanic` during session creation and attempt to reconnect by closing the browser and retrying once.

#### Scenario: Reconnect on EnsureSession panic
- **WHEN** `EnsureSession` receives `ErrBrowserPanic` from `NewSession`
- **THEN** the SessionManager SHALL close the browser tool
- **AND** the SessionManager SHALL retry `NewSession` exactly once
- **AND** if the retry succeeds, the new session ID SHALL be returned

#### Scenario: Reconnect retry fails
- **WHEN** `EnsureSession` receives `ErrBrowserPanic` and the retry also fails
- **THEN** the error SHALL be returned to the caller
- **AND** no further retries SHALL be attempted

### Requirement: Browser tool handler panic wrapper
The application layer SHALL wrap all browser tool handlers with panic recovery and retry logic.

#### Scenario: Handler-level panic recovery
- **WHEN** a browser tool handler panics during execution
- **THEN** the wrapper SHALL recover the panic and return an error wrapping `ErrBrowserPanic`

#### Scenario: Handler-level retry on ErrBrowserPanic
- **WHEN** a browser tool handler returns `ErrBrowserPanic`
- **THEN** the wrapper SHALL close the session manager and retry the handler once
- **AND** if the retry succeeds, the result SHALL be returned normally
