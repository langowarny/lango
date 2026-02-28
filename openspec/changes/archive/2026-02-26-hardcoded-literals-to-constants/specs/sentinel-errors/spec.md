## ADDED Requirements

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
