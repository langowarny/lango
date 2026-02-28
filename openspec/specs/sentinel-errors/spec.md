## ADDED Requirements

### Requirement: Session sentinel errors
The system SHALL define `ErrSessionNotFound` and `ErrDuplicateSession` in `session/errors.go`.

#### Scenario: Replace string matching with errors.Is
- **WHEN** `adk/session_service.go` checks for "session not found" errors
- **THEN** it SHALL use `errors.Is(err, session.ErrSessionNotFound)` instead of `strings.Contains`

#### Scenario: Replace UNIQUE constraint matching
- **WHEN** `adk/session_service.go` checks for duplicate session errors
- **THEN** it SHALL use `errors.Is(err, session.ErrDuplicateSession)` instead of string matching

### Requirement: Gateway sentinel errors
The system SHALL define `ErrNoCompanion`, `ErrApprovalTimeout`, `ErrAgentNotReady` in `gateway/errors.go`.

#### Scenario: Gateway error handling
- **WHEN** gateway operations encounter known error conditions
- **THEN** they SHALL return sentinel errors instead of ad-hoc error messages

### Requirement: Workflow sentinel errors
The system SHALL define `ErrWorkflowNameEmpty`, `ErrNoWorkflowSteps`, `ErrStepIDEmpty` in `workflow/errors.go`.

#### Scenario: Workflow validation errors
- **WHEN** workflow validation fails
- **THEN** it SHALL return sentinel errors for programmatic handling

### Requirement: Knowledge sentinel errors
The system SHALL define `ErrKnowledgeNotFound`, `ErrLearningNotFound` in `knowledge/errors.go`.

#### Scenario: Knowledge lookup errors
- **WHEN** knowledge or learning lookups find no results
- **THEN** they SHALL return sentinel errors

### Requirement: Security sentinel errors
The system SHALL define `ErrKeyNotFound`, `ErrNoEncryptionKeys`, `ErrDecryptionFailed` in `security/errors.go`.

#### Scenario: Security operation errors
- **WHEN** security operations encounter known failure modes
- **THEN** they SHALL return sentinel errors

### Requirement: Gateway RPCError type
The system SHALL define `RPCError` struct with `Code int` and `Message string` fields implementing the `error` interface in `gateway/errors.go`.

#### Scenario: Structured RPC errors
- **WHEN** `gateway/server.go` creates RPC error responses
- **THEN** it SHALL use the named `RPCError` type instead of anonymous structs

### Requirement: Protocol sentinel errors
The system SHALL define sentinel errors in `protocol/messages.go` for common P2P protocol error conditions: `ErrMissingToolName`, `ErrAgentCardUnavailable`, `ErrNoApprovalHandler`, `ErrDeniedByOwner`, `ErrExecutorNotConfigured`, `ErrInvalidSession`, `ErrInvalidPaymentAuth`.

#### Scenario: Handler uses sentinel errors
- **WHEN** the protocol handler encounters a known error condition (missing tool name, no card, no approval handler, denied by owner, no executor, invalid session, invalid payment)
- **THEN** it SHALL use the sentinel error's `.Error()` message in the response Error field

#### Scenario: Sentinel errors are matchable
- **WHEN** a caller receives a protocol error
- **THEN** it SHALL be able to use `errors.Is()` to match against the sentinel errors

### Requirement: Firewall sentinel errors
The system SHALL define sentinel errors in `firewall/firewall.go`: `ErrRateLimitExceeded`, `ErrGlobalRateLimitExceeded`, `ErrQueryDenied`, `ErrNoMatchingAllowRule`.

#### Scenario: Rate limit errors wrap sentinel
- **WHEN** a peer exceeds the rate limit
- **THEN** `FilterQuery` SHALL return an error wrapping `ErrRateLimitExceeded` with `%w`

#### Scenario: ACL deny errors wrap sentinel
- **WHEN** a firewall deny rule matches
- **THEN** `FilterQuery` SHALL return an error wrapping `ErrQueryDenied`

#### Scenario: No matching allow rule wraps sentinel
- **WHEN** no allow rule matches and default-deny applies
- **THEN** `FilterQuery` SHALL return an error wrapping `ErrNoMatchingAllowRule`

### Requirement: ZKP unsupported scheme error
The system SHALL define `ErrUnsupportedScheme` in `zkp/zkp.go`.

#### Scenario: Unknown scheme returns sentinel
- **WHEN** a ZKP operation encounters an unknown proving scheme
- **THEN** it SHALL return an error wrapping `ErrUnsupportedScheme`

### Requirement: Session expiry sentinel error
The system SHALL define `ErrSessionExpired` in `session/errors.go` alongside existing session sentinel errors.

#### Scenario: EntStore wraps TTL expiry with ErrSessionExpired
- **WHEN** `EntStore.Get()` finds a session whose `UpdatedAt` exceeds the configured TTL
- **THEN** it SHALL return an error wrapping `ErrSessionExpired` using `fmt.Errorf("get session %q: %w", key, ErrSessionExpired)`

#### Scenario: ErrSessionExpired is matchable via errors.Is
- **WHEN** a caller receives a TTL expiry error from `EntStore.Get()`
- **THEN** `errors.Is(err, ErrSessionExpired)` SHALL return `true`
