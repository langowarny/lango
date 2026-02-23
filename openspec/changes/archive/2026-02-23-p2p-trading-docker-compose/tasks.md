## 1. CLI Secrets Non-Interactive Support

- [x] 1.1 Add `--value-hex` string flag to `newSecretsSetCmd` in `internal/cli/security/secrets.go`
- [x] 1.2 Implement hex decode logic (strip `0x` prefix, `hex.DecodeString`) and store raw bytes
- [x] 1.3 Update error message for non-interactive terminals to suggest `--value-hex`

## 2. P2P REST API

- [x] 2.1 Create `internal/app/p2p_routes.go` with `registerP2PRoutes(chi.Router, *p2pComponents)`
- [x] 2.2 Implement `GET /api/p2p/status` handler (peer ID, listen addrs, connected peer count)
- [x] 2.3 Implement `GET /api/p2p/peers` handler (list connected peers with addrs)
- [x] 2.4 Implement `GET /api/p2p/identity` handler (DID from identity provider)
- [x] 2.5 Wire `registerP2PRoutes` in `internal/app/app.go` after gateway creation and P2P init

## 3. Docker Compose Example Structure

- [x] 3.1 Create `examples/p2p-trading/` directory with subdirs: configs, secrets, scripts, contracts
- [x] 3.2 Create `docker-compose.yml` with 5 services: anvil, setup, alice, bob, charlie
- [x] 3.3 Create `docker-entrypoint-p2p.sh` extending base entrypoint with USDC address wait, config substitution, and wallet key injection

## 4. Agent Configs and Secrets

- [x] 4.1 Create `configs/alice.json` (port 18789, P2P 9001, payment enabled, no LLM)
- [x] 4.2 Create `configs/bob.json` (port 18790, P2P 9002)
- [x] 4.3 Create `configs/charlie.json` (port 18791, P2P 9003)
- [x] 4.4 Create passphrase files in `secrets/` for each agent

## 5. Smart Contract and Setup

- [x] 5.1 Create `contracts/MockUSDC.sol` — minimal ERC-20 with mint, 6 decimals
- [x] 5.2 Create `scripts/setup-anvil.sh` — deploy MockUSDC, mint 1000 USDC to each agent

## 6. Test Scripts and Build

- [x] 6.1 Create `scripts/wait-for-health.sh` — poll URL until HTTP 200
- [x] 6.2 Create `scripts/test-p2p-trading.sh` — E2E test: health, P2P status, discovery, identity, balances, payment
- [x] 6.3 Create `Makefile` with build/up/test/down/clean targets
- [x] 6.4 Create `README.md` with architecture diagram and usage guide

## 7. Verification

- [x] 7.1 Run `go build ./...` — all packages compile
- [x] 7.2 Run `go test ./internal/app/...` — existing tests pass
