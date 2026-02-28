## ADDED Requirements

### Requirement: Session expiry sentinel error
The system SHALL define `ErrSessionExpired` in `session/errors.go` alongside existing session sentinel errors.

#### Scenario: EntStore wraps TTL expiry with ErrSessionExpired
- **WHEN** `EntStore.Get()` finds a session whose `UpdatedAt` exceeds the configured TTL
- **THEN** it SHALL return an error wrapping `ErrSessionExpired` using `fmt.Errorf("get session %q: %w", key, ErrSessionExpired)`

#### Scenario: ErrSessionExpired is matchable via errors.Is
- **WHEN** a caller receives a TTL expiry error from `EntStore.Get()`
- **THEN** `errors.Is(err, ErrSessionExpired)` SHALL return `true`
