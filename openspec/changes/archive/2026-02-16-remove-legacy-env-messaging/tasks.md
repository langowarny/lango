## 1. Onboard Post-Save Cleanup

- [x] 1.1 Remove env var export instructions and channel-specific token guidance from `printNextSteps()` in `internal/cli/onboard/onboard.go`
- [x] 1.2 Delete `generateEnvExample()` function and `.lango.env.example` file creation line
- [x] 1.3 Remove unused imports (`os`, `common`) from `internal/cli/onboard/onboard.go`
- [x] 1.4 Simplify `printNextSteps()` to show only: confirmation, serve, doctor, and profile management commands

## 2. Onboard Long Description Update

- [x] 2.1 Add Auth (OIDC providers, JWT settings) and Session (Session DB, TTL) categories to Long description
- [x] 2.2 Update Security category to list only PII interceptor and Signer
- [x] 2.3 Update description to note encrypted profile storage

## 3. Doctor API Key Security Check

- [x] 3.1 Change warn message from "Plaintext API keys detected" to "Inline API keys for: <providers>"
- [x] 3.2 Update warn Details to explain encrypted profiles are safe and suggest `${ENV_VAR}` for portability
- [x] 3.3 Change pass message from "All API keys use environment variable references" to "All API keys secured"

## 4. Verification

- [x] 4.1 Run `go build ./...` — no compilation errors
- [x] 4.2 Run `go test ./...` — all tests pass
