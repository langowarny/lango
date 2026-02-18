## Context

Lango's core differentiator is preventing AI agents from leaking secrets. The current architecture uses CryptoProvider (Signer Provider) for encryption at rest, but `secrets_get` and `crypto_decrypt` tools return plaintext directly to the agent context. The approval middleware uses fail-open semantics, meaning sensitive tools execute without approval when no companion is connected.

The real threat is not "stealing the encryption key" but "obtaining plaintext through the tool API and including it in context."

## Goals / Non-Goals

**Goals:**
- Prevent secret plaintext from ever entering the AI agent's LLM context
- Ensure sensitive tools cannot execute without explicit approval (fail-closed)
- Detect and mask any secret values that leak into tool results or model responses
- Maintain backward compatibility for exec tool usage patterns (reference tokens resolve transparently)

**Non-Goals:**
- Removing or replacing the existing CryptoProvider/Signer Provider (it remains valid for encryption at rest)
- Implementing capability-based tool access profiles (P2, deferred)
- Implementing egress content inspection for exec/browser tools (P2, deferred)
- Implementing short-lived tokens or scoped credentials (P3, deferred)
- Changing the secrets storage schema or migration

## Decisions

### Decision 1: Proxy/Reference Token Pattern over Direct Plaintext Return

**Choice**: Return opaque tokens (`{{secret:name}}`, `{{decrypt:id}}`) instead of plaintext values.

**Alternatives considered**:
- **Direct plaintext with output scanning only**: Plaintext still enters agent context; scanning is a detection layer, not prevention. Insufficient.
- **Server-side execution proxy**: Agent sends intent, server executes with secrets injected. Too complex; changes the entire execution model.
- **Encrypted return with agent-side key**: Agent would need the key, defeating the purpose.

**Rationale**: The reference token pattern is the industry standard (HashiCorp Vault Agent, AWS ECS Task Role Credential Injection). It fundamentally prevents the agent from seeing secret values while allowing transparent usage in exec commands.

**Trade-off**: Agent cannot inspect secret values (e.g., check API key expiry). This is acceptable for most use cases.

### Decision 2: Fail-Closed Approval with TTY Fallback

**Choice**: Three-tier approval: companion → TTY prompt → deny.

**Alternatives considered**:
- **Fail-closed without TTY fallback**: Breaks local development where companion is rarely connected.
- **Keep fail-open with logging**: Does not provide actual security.
- **Always require companion**: Too strict for development workflows.

**Rationale**: Fail-closed is the correct security default. TTY fallback provides usability for local development without compromising the security model in production (headless) environments.

### Decision 3: Output Scanning as Defense-in-Depth

**Choice**: Scan both tool results (input to LLM) and model responses (output from LLM) for known secret values.

**Rationale**: Even with reference tokens, edge cases exist (e.g., a tool returning data that happens to contain a secret value from another source). Output scanning catches these. Minimum value length of 4 characters avoids false positives.

### Decision 4: Session-Scoped RefStore (In-Memory)

**Choice**: RefStore is in-memory, scoped to the application lifecycle. No persistence.

**Rationale**: Reference tokens are ephemeral. Persisting them would create another attack surface. Memory-only ensures tokens are cleared on restart.

## Risks / Trade-offs

- **[Risk] Agent cannot inspect secret values** → Acceptable trade-off. If inspection is needed, a dedicated `secrets_inspect` tool with strict approval could be added later.
- **[Risk] Reference tokens in logs** → Tokens like `{{secret:api_key}}` reveal secret names but not values. Name exposure is low-risk compared to value exposure.
- **[Risk] Scanner performance with many secrets** → O(n*m) where n=secrets, m=text length. Acceptable for typical secret counts (<100). If needed, Aho-Corasick can be adopted.
- **[Risk] Exec command with unresolved token** → If RefStore doesn't have the token, the literal `{{secret:name}}` string is passed to the shell. This fails safely (command error) rather than leaking data.
- **[Risk] TTY prompt in automated pipelines** → Non-TTY environments get hard denial, which is the correct behavior for CI/CD.
