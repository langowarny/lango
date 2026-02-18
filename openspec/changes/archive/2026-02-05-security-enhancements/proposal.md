## Why

Lango aims to be a production-ready, secure AI agent platform. Currently, it lacks:
1.  **Defense-in-Depth for AI**: No mechanism to redact PII or block sensitive actions before they reach the LLM provider.
2.  **Hardware-Backed Security**: Cryptographic keys are likely file-based (or non-existent), making them vulnerable to theft.
3.  **Standard Identity**: Authentication (if any) is likely basic. We need industry-standard OIDC (Google/GitHub) to support secure multi-user or enterprise scenarios.

## What Changes

*   **AI Privacy Interceptor**:
    *   Implement a middleware/decorator pattern in `internal/agent` or `internal/gateway`.
    *   Inspects prompt content for PII (Regex/NLP).
    *   Blocks sensitive tool calls (e.g., `exec` with dangerous commands) unless approved via a side-channel (Telegram/Discord).
*   **Secure Enclave / MPC Integration**:
    *   Shift key management from direct file access to a "Signer Interface".
    *   Implement a MacOS Secure Enclave provider (using `cgo` or external helper) to sign requests without exposing private keys.
*   **OIDC Authentication**:
    *   Update `internal/gateway` (HTTP server) to support OIDC login flows (Google/GitHub).
    *   Protect API endpoints with JWTs issued after OIDC auth.

## Capabilities

### New Capabilities
- `privacy-interceptor`: Middleware for intercepting and sanitizing AI interactions.
- `secure-signer`: Interface and implementation for hardware-backed (Secure Enclave) cryptographic signing.
- `oidc-auth`: OpenID Connect authentication handler for multiple providers.

### Modified Capabilities
- `agent-runtime`: Update `internal/agent` to support the interceptor chain.
- `gateway-server`: Update `internal/gateway` to include auth middleware and OIDC routes.

## Impact

*   **Core**: `internal/agent`, `internal/gateway`, `internal/app`.
*   **Configuration**: New config sections for OIDC providers, Privacy Rules, and Key Management.
*   **Dependencies**: `coreos/go-oidc`, `golang.org/x/oauth2`.
