## Context

The P2P tool invocation pipeline (`handler.go` → `approvalFn` → `composite.go` → executor) has 5 security gaps discovered during audit. Sandbox isolation, KMS, and security events were already added in recent commits, but the approval path itself—the decision layer that determines whether a remote peer's tool request should execute—had critical bypass vectors.

Current state: remote peers authenticate via handshake sessions, pass through firewall ACL, and then reach the approval check. If `approvalFn` is nil (e.g., no payment module), the check is silently skipped. HeadlessProvider (designed for Docker environments) auto-approves everything, including P2P sessions. Grant entries never expire.

## Goals / Non-Goals

**Goals:**
- Eliminate all known approval-path bypass vectors for P2P tool invocations
- Ensure fail-closed behavior at every decision point (nil handler, missing provider, unknown tool)
- Prevent HeadlessProvider from ever approving P2P remote peer requests
- Block overly permissive firewall rules that would allow any peer to access any tool
- Add time-based expiration to approval grants to limit implicit trust windows
- Prevent double-prompting when handler approval and tool-level approval both trigger

**Non-Goals:**
- Changing the handshake/authentication layer (already hardened in P2-10/11/12)
- Adding new approval provider types (e.g., webhook-based approval for CI/CD)
- Modifying sandbox or container isolation behavior
- Persistent grant storage (grants remain in-memory, cleared on restart)

## Decisions

### 1. Default-deny on nil approvalFn (vs. optional approval)
**Decision**: Return "denied" response when `approvalFn` is nil.
**Rationale**: The previous behavior silently skipped approval, which is fail-open. Any code path that reaches the handler without configuring approval (e.g., P2P enabled without payment module) would execute tools unconditionally. Default-deny is the only safe choice for a security boundary.
**Alternative considered**: Making approvalFn required in HandlerConfig. Rejected because it would break backward compatibility and the handler is also used for non-tool requests (agent card, capability query).

### 2. Dedicated P2P fallback slot in CompositeProvider (vs. prefix-matching)
**Decision**: Add `p2pFallback` field that intercepts all `"p2p:..."` session keys before TTY fallback.
**Rationale**: The TTY fallback slot is shared between local and P2P sessions. When HeadlessProvider occupies it, P2P sessions get auto-approved. A dedicated slot ensures P2P sessions are always routed to an appropriate provider regardless of the TTY fallback configuration.
**Alternative considered**: Adding a `CanHandleP2P()` method to Provider interface. Rejected as it would require changes to all provider implementations.

### 3. SafetyLevel check before price-based auto-approve
**Decision**: Check `tool.SafetyLevel.IsDangerous()` before any auto-approve logic. Unknown tools (not in toolIndex) are treated as dangerous.
**Rationale**: A low-priced dangerous tool (e.g., `payment_send` at $0.01) should never be auto-approved by a remote peer. The SafetyLevel metadata already exists on all tools; using it here closes the gap without introducing new abstractions.

### 4. AddRule returns error (vs. silent reject)
**Decision**: Change `AddRule(rule)` from void to `error` return, with `ValidateRule()` as a separate public function.
**Rationale**: Silent rejection would hide configuration errors. Returning an error lets callers (CLI tools, config loaders) provide actionable feedback. `ValidateRule()` is public so it can be used independently for validation UIs.
**Migration**: Existing `New()` constructor warns but still loads overly permissive rules for backward compatibility.

### 5. Grant TTL with per-field timestamps
**Decision**: Replace `map[string]struct{}` with `map[string]grantEntry{grantedAt}`. TTL defaults to 0 (no expiry) for backward compatibility; P2P sets 1-hour TTL.
**Rationale**: Indefinite grants mean a single approval creates permanent trust. TTL bounds the window. The `grantedAt` field enables per-entry expiration without requiring a background goroutine (lazy expiration on `IsGranted` + explicit `CleanExpired`).

### 6. Double-approval prevention via grant recording
**Decision**: When P2P `approvalFn` approves a tool, immediately record a grant for `"p2p:"+peerDID`. The tool's internal `wrapWithApproval` checks `IsGranted` and skips the second prompt.
**Rationale**: Without this, the user would see two approval prompts for one remote tool call (one from handler, one from tool wrapper). The grant is TTL-bounded so it doesn't create permanent trust.

## Risks / Trade-offs

- **[TTY unavailable in headless P2P]** → P2P sessions with TTYProvider as fallback will fail with "stdin is not a terminal" in headless environments. Mitigation: users must either use a Gateway companion for approval or disable P2P in headless mode. This is intentionally fail-closed.
- **[AddRule API break]** → Callers that ignore the return value will get a compile error. Mitigation: only one caller (`p2p_firewall_add` tool handler) exists; already updated.
- **[Grant TTL clock skew]** → `nowFn` is injectable for testing but uses `time.Now` in production. If system clock jumps, grants may expire prematurely or persist too long. Mitigation: 1-hour TTL is coarse enough that clock jitter is negligible.
- **[Backward-compatible overly permissive rules]** → `New()` still loads wildcard allow rules with a warning. Mitigation: `AddRule()` rejects them going forward; existing configs get logged warnings encouraging cleanup.
