## Context

Lango's P2P networking and payment systems are implemented across `internal/p2p/`, `internal/payment/`, and `internal/wallet/`. The gateway (`internal/gateway/`) serves HTTP and WebSocket endpoints on a chi router. Currently, P2P node state is only accessible via CLI commands that create ephemeral libp2p nodes — separate from the running server's node. Wallet private keys can only be injected interactively via `lango security secrets set`, blocking Docker-based automation.

## Goals / Non-Goals

**Goals:**
- Enable non-interactive wallet key injection for Docker/CI environments
- Expose P2P node state (peer ID, connected peers, DID) via REST API on the running gateway
- Provide a complete Docker Compose example that proves 3 agents can discover each other and transact USDC

**Non-Goals:**
- Replacing CLI P2P commands (they remain for ad-hoc debugging)
- Adding LLM provider integration to the Docker example (tests P2P + payment only)
- Production-grade Docker deployment (this is an integration test example)
- Modifying the P2P or payment core logic

## Decisions

### 1. P2P routes live in `internal/app/p2p_routes.go` (not `internal/gateway/`)

**Rationale**: The routes depend on `p2pComponents` which is an app-layer type from `wiring.go`. Placing them in gateway would create an import cycle or require leaking p2p internals into the gateway package. This follows the existing pattern where A2A routes are registered from app.go via `a2aServer.RegisterRoutes(router)`.

**Alternative**: Create a gateway sub-package — rejected as over-engineering for 3 simple handlers.

### 2. `--value-hex` flag (not stdin pipe)

**Rationale**: A hex string flag is simplest for Docker entrypoints where the value comes from an environment variable. Stdin piping (`echo $KEY | lango secrets set`) would require additional plumbing and is less explicit. The `0x` prefix is optionally stripped to match Ethereum key conventions.

**Alternative**: `--value-file` reading from a file — could be added later but hex flag covers the immediate need.

### 3. Anvil deterministic accounts for agents

**Rationale**: Anvil generates the same 10 accounts on every run. Using accounts 0-2 for Alice/Bob/Charlie and account 9 for the deployer avoids key generation complexity and ensures test reproducibility.

### 4. MockUSDC instead of real ERC-20 fork

**Rationale**: A minimal 50-line Solidity contract with `mint()` is simpler and faster than forking mainnet USDC. The payment system interacts via standard ERC-20 `transfer`/`balanceOf`, so the mock is functionally equivalent for integration testing.

### 5. P2P REST endpoints are public (no auth middleware)

**Rationale**: The endpoints expose only node metadata (peer ID, listen addresses, DID). No secrets or session data are returned. The existing `/health` endpoint follows the same pattern. In production, operators would use network-level access control.

## Risks / Trade-offs

- **mDNS in Docker**: Docker bridge networks support multicast by default, but some Docker Desktop configurations may block it → Mitigation: 15-second wait with retry in test script; fallback to explicit bootstrap peers if needed
- **Test flakiness**: P2P discovery timing is non-deterministic → Mitigation: generous timeouts and retry loops in test script
- **MockUSDC divergence**: Mock may not match real USDC behavior for edge cases → Mitigation: tests only use basic `transfer`/`balanceOf` which are standard ERC-20
- **Foundry image availability**: `ghcr.io/foundry-rs/foundry` may change tags → Mitigation: use `latest` tag; pin version in production
