## 1. Setup & Configuration

- [x] 1.1 Update `lango.json` to include new config sections (`auth`, `security.interceptor`, `security.signer`).
- [x] 1.2 Update `internal/config/config.go` to parse the new configuration sections.

## 2. AI Privacy Interceptor

- [x] 2.1 Create `internal/agent/middleware.go` defining the `RuntimeMiddleware` interface.
- [x] 2.2 Implement `PIIRedactor` using regex patterns for email/phone/keys.
- [x] 2.3 Implement `ApprovalMiddleware` that checks for sensitive tool calls.
- [x] 2.4 Update `internal/app/app.go` to wrap the `Agent` with the `PrivacyInterceptor`.

## 3. Secure Signing

- [x] 3.1 Define `Signer` interface in `internal/security`.
- [x] 3.2 Implement `RPCSigner` carrier in `internal/security/rpc_signer.go`.
- [x] 3.3 Update `internal/server/server.go` (Gateway) to handle `sign.request` RPC messages from the Agent.
- [x] 3.4 Wire the `Signer` into the system where keys are used (e.g., Session management or custom tool signing).

## 4. OIDC Authentication

- [x] 4.1 Create `internal/server/auth.go` with `OIDCProvider` struct using `coreos/go-oidc`.
- [x] 4.2 Add `/auth/login/{provider}` handler to redirect to OIDC provider.
- [x] 4.3 Add `/auth/callback/{provider}` handler to exchange code for token and issue session.
- [x] 4.4 Update `internal/server/server.go` routing to include these new endpoints.

## 5. Verification

- [x] 5.1 Add unit tests for `PIIRedactor` to ensure patterns catch standard PII.
- [ ] 5.2 Add integration test for `OIDCProvider` (mocking the upstream provider). (Manual Verification Required: Setup Google Console project)
- [ ] 5.3 Verify RPC Signing (this might require manual testing with a stubbed host). (Manual Verification Required: Host app integration)
