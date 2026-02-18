## Why

The current Signer Provider architecture protects encryption keys at rest but fails to prevent AI agents from accessing secret plaintext through tool APIs. `secrets_get` and `crypto_decrypt` return plaintext directly into the agent context, and the approval middleware uses a fail-open strategy that bypasses approval when no companion is connected. This fundamentally undermines the AI agent secret isolation goal.

## What Changes

- **BREAKING**: `secrets_get` no longer returns plaintext values. Returns opaque reference tokens (`{{secret:name}}`) that are resolved at execution time by the exec tool.
- **BREAKING**: `crypto_decrypt` no longer returns plaintext data. Returns opaque reference tokens (`{{decrypt:id}}`) resolved at execution time.
- **BREAKING**: `wrapWithApproval` changes from fail-open to fail-closed. Sensitive tools are denied by default without explicit approval.
- New `RefStore` manages mapping between reference tokens and secret values with session-scoped lifecycle.
- New `SecretScanner` detects leaked secret values in tool results and model responses, replacing them with `[SECRET:name]` markers.
- `PIIRedactingModelAdapter` extended to scan both input (PII) and output (secrets) directions.
- Exec tool resolves reference tokens just before shell execution; resolved values never appear in logs or agent context.

## Capabilities

### New Capabilities
- `secret-reference-tokens`: Opaque reference token system that prevents secret plaintext from entering the AI agent context. Tokens are resolved at execution time by the exec tool.
- `output-secret-scanning`: Scans tool results and model responses for known secret values, masking them with `[SECRET:name]` placeholders to prevent leakage through output channels.

### Modified Capabilities
- `ai-privacy-interceptor`: Approval workflow changes from fail-open to fail-closed. TTY prompt fallback added when no companion is connected.
- `security-tools`: `secrets_get` and `crypto_decrypt` return reference tokens instead of plaintext. Tools now depend on `RefStore` and `SecretScanner`.
- `tool-exec`: Exec tool gains reference token resolution capability, replacing `{{secret:name}}` and `{{decrypt:id}}` tokens with actual values at execution time.

## Impact

- **Core security module**: New `internal/security/secret_ref.go` (RefStore)
- **Agent module**: New `internal/agent/secret_scanner.go` (SecretScanner)
- **App wiring**: `internal/app/tools.go`, `internal/app/app.go`, `internal/app/wiring.go` updated for fail-closed approval and RefStore/Scanner plumbing
- **Tool interfaces**: `internal/tools/secrets/secrets.go`, `internal/tools/crypto/crypto.go` constructor signatures changed (breaking for direct callers)
- **Exec tool**: `internal/tools/exec/exec.go` gains reference resolution
- **Model adapter**: `internal/adk/pii_model.go` extended with bidirectional scanning
- **API consumers**: Any code calling `secrets_get` expecting plaintext in `value` field will receive a reference token instead
