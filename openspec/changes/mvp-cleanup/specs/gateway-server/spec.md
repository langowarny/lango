## REMOVED Requirements

### Requirement: Companion WebSocket Handler
**Reason**: Companion protocol removed from MVP. No native companion app integration in initial release.
**Migration**: Re-enable in Phase 2 by restoring `/companion` WebSocket endpoint and RPCProvider sender wiring.

### Requirement: Companion Message Routing
**Reason**: Part of companion protocol, removed from MVP.
**Migration**: Same as above.

### Requirement: Approval Request Broadcasting
**Reason**: Part of companion protocol and security approval workflow, removed from MVP.
**Migration**: Same as above.

## MODIFIED Requirements

### Requirement: Gateway Initialization
The gateway server SHALL be initialized without requiring an `AuthManager` or `RPCProvider`. The `gateway.New()` function SHALL accept `nil` for optional parameters (rpcProvider, authManager). The gateway SHALL serve HTTP and WebSocket endpoints for direct chat without OIDC authentication.

#### Scenario: Gateway startup without auth
- **WHEN** the gateway is created with `nil` authManager and `nil` rpcProvider
- **THEN** it SHALL start successfully
- **THEN** it SHALL serve the `chat.message` RPC endpoint
- **THEN** it SHALL NOT register `/companion` WebSocket endpoint
- **THEN** it SHALL NOT require OIDC configuration
