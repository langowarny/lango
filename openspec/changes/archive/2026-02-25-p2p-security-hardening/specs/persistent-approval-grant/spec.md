## ADDED Requirements

### Requirement: Grant TTL expiration
GrantStore SHALL support an optional time-to-live (TTL) for grants. When TTL is set to a positive duration, `IsGranted()` MUST check whether the grant has expired (current time minus `grantedAt` exceeds TTL). A TTL of zero MUST preserve backward-compatible behavior (no expiration).

#### Scenario: Grant within TTL
- **WHEN** a grant was created 5 minutes ago and TTL is 10 minutes
- **THEN** `IsGranted()` SHALL return true

#### Scenario: Grant expired past TTL
- **WHEN** a grant was created 11 minutes ago and TTL is 10 minutes
- **THEN** `IsGranted()` SHALL return false

#### Scenario: TTL zero means no expiry
- **WHEN** TTL is zero (default) and a grant was created 100 hours ago
- **THEN** `IsGranted()` SHALL return true

### Requirement: Clean expired grants
GrantStore SHALL provide a `CleanExpired()` method that removes all grants whose `grantedAt` timestamp exceeds the configured TTL. The method SHALL return the count of removed entries. When TTL is zero, `CleanExpired()` SHALL be a no-op returning zero.

#### Scenario: Clean expired entries
- **WHEN** `CleanExpired()` is called with TTL of 5 minutes and 2 of 3 grants are older than 5 minutes
- **THEN** the method SHALL remove the 2 expired grants and return 2

#### Scenario: Clean with zero TTL
- **WHEN** `CleanExpired()` is called with TTL of zero
- **THEN** the method SHALL remove nothing and return 0

## MODIFIED Requirements

### Requirement: P2P grant TTL default
When P2P is enabled, the application SHALL set the GrantStore TTL to 1 hour. This limits the window of implicit trust from P2P approval grants.

#### Scenario: P2P enabled sets 1-hour TTL
- **WHEN** the application initializes with `cfg.P2P.Enabled = true`
- **THEN** `grantStore.SetTTL(time.Hour)` SHALL be called

### Requirement: Double-approval prevention via grant recording
When the P2P approval function approves a tool invocation, the system SHALL record a grant for `"p2p:"+peerDID` and the tool name. This prevents the tool's internal `wrapWithApproval` from prompting a second time.

#### Scenario: Approved P2P tool records grant
- **WHEN** the P2P approval function approves tool "echo" for peer "did:key:abc"
- **THEN** a grant SHALL be recorded with session key `"p2p:did:key:abc"` and tool name `"echo"`
