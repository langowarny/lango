## 1. Gateway HTTP API Docs

- [x] 1.1 Add P2P Network section to `docs/gateway/http-api.md` with `GET /api/p2p/status`, `/peers`, `/identity` endpoints, curl examples, and JSON response formats
- [x] 1.2 Add note explaining REST API vs CLI ephemeral node distinction

## 2. CLI Security Docs

- [x] 2.1 Update `docs/cli/security.md` `secrets set` section with `--value-hex` flag, flags table, and non-interactive examples
- [x] 2.2 Add tip about using `--value-hex` in Docker/CI environments

## 3. Agent Prompts

- [x] 3.1 Add REST API endpoint reference to `prompts/TOOL_USAGE.md` P2P Networking Tool section
- [x] 3.2 Add REST API mention to `prompts/agents/vault/IDENTITY.md` Output Format

## 4. P2P Feature Docs

- [x] 4.1 Add REST API section to `docs/features/p2p-network.md` with endpoint table and curl examples
- [x] 4.2 Clarify CLI commands create ephemeral nodes

## 5. README Updates

- [x] 5.1 Add REST API subsection under P2P Network section in `README.md`
- [x] 5.2 Add Examples section with P2P Trading Docker Compose link
- [x] 5.3 Update `secrets set` CLI line to mention `--value-hex`

## 6. CLI Help Text

- [x] 6.1 Add `Long` description to `internal/cli/p2p/status.go` referencing REST API
- [x] 6.2 Add `Long` description to `internal/cli/p2p/peers.go` referencing REST API
- [x] 6.3 Add `Long` description to `internal/cli/p2p/identity.go` referencing REST API

## 7. Docs Site Homepage

- [x] 7.1 Add P2P Network card to `docs/index.md` features grid

## 8. Verification

- [x] 8.1 Run `go build ./...` — all packages compile
