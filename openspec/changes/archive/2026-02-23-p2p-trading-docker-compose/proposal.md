## Why

The P2P/A2A architecture is implemented but lacks an end-to-end integration test proving that multiple agents can discover each other via mDNS, establish P2P connections, and execute USDC payments over a local blockchain. A Docker Compose example with 3 agents (Alice, Bob, Charlie) and a local Anvil node provides a reproducible verification environment. Two blockers prevent headless agent setup: CLI P2P commands create ephemeral nodes (not the running server's), and wallet key injection requires an interactive terminal.

## What Changes

- Add `--value-hex` flag to `lango security secrets set` for non-interactive hex-encoded secret injection (e.g., wallet private keys in Docker)
- Add P2P REST API endpoints (`/api/p2p/status`, `/api/p2p/peers`, `/api/p2p/identity`) on the gateway router, enabling external tooling to query P2P node state without ephemeral CLI nodes
- Create `examples/p2p-trading/` Docker Compose integration example with 3 Lango agents, Anvil (local EVM), MockUSDC contract deployment, and an E2E test script verifying health, P2P discovery, DID identity, and USDC payment

## Capabilities

### New Capabilities
- `p2p-rest-api`: REST endpoints for querying P2P node status, connected peers, and local DID identity via the gateway
- `p2p-trading-example`: Docker Compose integration example with 3 agents, local blockchain, and E2E test scripts

### Modified Capabilities
- `cli-secrets-management`: Add `--value-hex` flag to `secrets set` for non-interactive hex value injection

## Impact

- **Modified files**: `internal/cli/security/secrets.go` (new flag), `internal/app/app.go` (P2P route wiring)
- **New files**: `internal/app/p2p_routes.go` (P2P REST handlers), full `examples/p2p-trading/` directory tree
- **Dependencies**: No new Go dependencies; Docker example uses Foundry (Anvil/Forge/Cast) images
- **APIs**: New public REST endpoints at `/api/p2p/*` (no auth required, metadata only)
