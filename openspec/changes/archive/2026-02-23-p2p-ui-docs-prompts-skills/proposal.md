## Why

The P2P Trading Docker Compose integration (`p2p-trading-docker-compose`) added P2P REST API endpoints, a `--value-hex` secrets flag, and a Docker example, but the surrounding documentation, prompts, and UI/UX layers were not updated. Users reading the docs, README, or agent prompts would not know about the new REST API, the non-interactive secrets flag, or the P2P trading example.

## What Changes

- Add P2P REST API documentation to `docs/gateway/http-api.md` (3 new endpoints)
- Add `--value-hex` flag documentation to `docs/cli/security.md`
- Add REST API reference to `prompts/TOOL_USAGE.md` P2P section
- Add REST API section to `docs/features/p2p-network.md`
- Add Examples section with P2P trading link to `README.md`
- Add P2P REST API mention to `README.md` P2P section
- Update `prompts/agents/vault/IDENTITY.md` with REST API awareness
- Add `Long` descriptions to P2P CLI commands referencing REST API
- Add P2P Network card to `docs/index.md` features grid

## Capabilities

### New Capabilities

(none — all changes are documentation/prompt updates)

### Modified Capabilities

(none — no spec-level requirement changes, only documentation)

## Impact

- **Docs**: `docs/gateway/http-api.md`, `docs/cli/security.md`, `docs/features/p2p-network.md`, `docs/index.md`
- **README**: `README.md` (examples section, P2P REST API, secrets flag hint)
- **Prompts**: `prompts/TOOL_USAGE.md`, `prompts/agents/vault/IDENTITY.md`
- **CLI**: `internal/cli/p2p/status.go`, `internal/cli/p2p/peers.go`, `internal/cli/p2p/identity.go` (Long descriptions only)
