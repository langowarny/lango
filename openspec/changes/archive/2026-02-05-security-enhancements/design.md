## Context

Lango is a Go-based AI agent platform. The current architecture uses a monolith `app` structure where `internal/agent` handles AI logic and `internal/gateway` exposes HTTP/WebSocket endpoints. Security is currently minimal. We need to introduce enterprise-grade security features: PII protection, hardware-backed keys, and standard identity.

## Goals / Non-Goals

**Goals:**
*   Prevent sensitive data (PII) from leaking to LLMs.
*   Prevent unauthorized high-risk tool usage via a "human-in-the-loop" approval flow.
*   Secure cryptographic keys using the MacOS Secure Enclave (when running on supported hardware).
*   Authenticate users via Google/GitHub OIDC.

**Non-Goals:**
*   Implementing a full-blown IAM system (RBAC/ABAC is out of scope for now, just authentication).
*   Supporting all possible HSMs (just MacOS Secure Enclave for this iteration).

## Decisions

### 1. Interceptor Pattern for AI Privacy
**Decision:** Implement a decorator/middleware pattern around the `agent.Initializable` or `agent.Runtime` interface.
**Rationale:** This allows transparent inspection of inputs (prompts) and blocking of tool execution without modifying the core agent logic.
**Architecture:**
```go
type PrivacyInterceptor struct {
    Next agent.Runtime
    Redactor *PIIRedactor
    Approver *ApprovalWorkflow
}
```

### 2. RPC-based Secure Signing
**Decision:** Use a bidirectional WebSocket connection (or local IPC) to delegate signing to a host process (e.g., a Swift app wrapping the Secure Enclave) if the Go app cannot access it directly (due to sandboxing or complexity).
**Rationale:** Go's `cgo` access to MacOS Security framework can be tricky, especially if the app is not bundled deeply. An RPC approach provides decoupling.
**Fallback:** If no host signer is available, fallback to a file-based encrypted keystore (with a warning).

### 3. OIDC Integration
**Decision:** Use `coreos/go-oidc` for OIDC and `golang.org/x/oauth2` for token exchange.
**Rationale:** Standard, battle-tested libraries.
**Implementation:** Add `/auth/login/{provider}` and `/auth/callback/{provider}` routes to `internal/gateway`.

## Risks / Trade-offs

*   **Risk**: Interceptor latency.
    *   **Mitigation**: Run PII redaction synchronously (fast regex). Run approval flow asynchronously with long timeouts (5-10m).
*   **Risk**: Secure Enclave unavailablity.
    *   **Mitigation**: Robust fallback to software keys.
