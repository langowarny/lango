## ADDED Requirements

### Requirement: P2P approval fallback isolation
The CompositeProvider SHALL provide a dedicated P2P fallback slot (`p2pFallback`) that is used exclusively for approval requests with session keys prefixed with `"p2p:"`. P2P sessions MUST never be routed to the TTY fallback slot, preventing HeadlessProvider from auto-approving remote peer requests.

#### Scenario: P2P session with no P2P fallback configured
- **WHEN** a P2P approval request (session key `"p2p:..."`) is received and no P2P fallback is set
- **THEN** the provider SHALL return an error stating "headless auto-approve is not allowed for remote peers"

#### Scenario: P2P session routes to dedicated fallback
- **WHEN** a P2P approval request is received and a P2P fallback provider is configured
- **THEN** the request SHALL be routed to the P2P fallback provider, not the TTY fallback

#### Scenario: Non-P2P session still uses TTY fallback
- **WHEN** a non-P2P approval request (session key without `"p2p:"` prefix) is received
- **THEN** the request SHALL be routed to the TTY fallback as before

#### Scenario: HeadlessProvider as TTY fallback with P2P request
- **WHEN** HeadlessProvider is configured as TTY fallback and a P2P approval request arrives
- **THEN** HeadlessProvider SHALL NOT be called; the request SHALL use the P2P fallback or be denied

## MODIFIED Requirements

### Requirement: P2P approval wiring
When P2P is enabled, the application SHALL configure `TTYProvider` as the P2P fallback on `CompositeProvider`. This ensures P2P approval requests are always routed to an interactive provider, regardless of whether HeadlessProvider is configured as the TTY fallback.

#### Scenario: P2P enabled wiring
- **WHEN** the application initializes with `cfg.P2P.Enabled = true`
- **THEN** `composite.SetP2PFallback(&approval.TTYProvider{})` SHALL be called
