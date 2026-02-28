## ADDED Requirements

### Requirement: Auto-renew expired sessions
The `SessionServiceAdapter.Get()` SHALL automatically delete an expired session and create a fresh replacement when the store returns `ErrSessionExpired`, so the user's current message is processed normally.

#### Scenario: Expired Telegram session auto-renews
- **WHEN** `SessionServiceAdapter.Get()` receives `ErrSessionExpired` for session `telegram:123:456`
- **THEN** the system SHALL delete the expired session, create a new session with the same key, and return it successfully

#### Scenario: Expired session delete failure propagates error
- **WHEN** `SessionServiceAdapter.Get()` receives `ErrSessionExpired` and the subsequent `Delete()` call fails
- **THEN** the system SHALL return the delete error wrapped with context, without attempting to create a new session

#### Scenario: Concurrent expiry recovery is safe
- **WHEN** multiple goroutines detect the same expired session simultaneously
- **THEN** the `getOrCreate()` retry logic SHALL ensure all goroutines return a valid session without errors

## MODIFIED Requirements

### Requirement: Non-recoverable store errors propagated
The `SessionServiceAdapter.Get()` SHALL propagate store errors that are not "session not found" or "session expired" (e.g., database connection failures).

#### Scenario: Database error during get
- **WHEN** the store returns an error other than "session not found" or "session expired"
- **THEN** the system SHALL propagate that error to the caller without attempting auto-creation or renewal
