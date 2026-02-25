## Why

The P2P remote tool invocation path has 5 critical security gaps in its approval pipeline. When `approvalFn` is nil, tools execute without any owner consent. HeadlessProvider can auto-approve P2P requests from remote peers. Dangerous tools (e.g., `payment_send`) can be auto-approved via P2P if the price is low. Firewall allows overly permissive wildcard allow rules. Approval grants never expire, creating an indefinite trust window.

## What Changes

- **Default-deny when approvalFn is nil**: Both `handleToolInvoke()` and `handleToolInvokePaid()` now return "denied" instead of silently skipping the approval check.
- **P2P-specific approval fallback**: CompositeProvider gains a dedicated `p2pFallback` slot. P2P sessions (`"p2p:..."` keys) never route to HeadlessProvider.
- **SafetyLevel enforcement for P2P auto-approve**: Dangerous tools are never auto-approved via P2P, regardless of price. Unknown tools are treated as dangerous.
- **Firewall wildcard rule validation**: New `ValidateRule()` rejects allow rules with wildcard peer + wildcard tools. `AddRule()` now returns an error.
- **Grant TTL support**: GrantStore gains time-based expiration. P2P grants default to 1-hour TTL. `CleanExpired()` removes stale entries.
- **Double-approval prevention**: P2P approvalFn records grants so tools' internal `wrapWithApproval` skips the second prompt.

## Capabilities

### New Capabilities

### Modified Capabilities

- `p2p-protocol`: handler.go default-deny when approvalFn is nil; both invoke paths affected
- `approval-policy`: CompositeProvider P2P fallback slot; HeadlessProvider blocked for P2P sessions
- `persistent-approval-grant`: TTL-based expiration for grants; CleanExpired cleanup method
- `p2p-firewall`: ValidateRule function; AddRule returns error; overly permissive rules rejected
- `tool-safety-level`: P2P auto-approve respects SafetyLevel; dangerous tools require explicit approval

## Impact

- **Files modified**: `internal/p2p/protocol/handler.go`, `internal/approval/composite.go`, `internal/approval/grant.go`, `internal/p2p/firewall/firewall.go`, `internal/app/app.go`, `internal/app/tools.go`
- **API change**: `firewall.AddRule()` now returns `error` (previously void) — callers must handle the error
- **Behavior change**: P2P requests denied by default if no approval handler is set; headless environments must configure a P2P-compatible approval provider or disable P2P
- **New tests**: handler_test.go (7 cases), firewall_test.go (6 cases), approval_test.go (3 new cases), grant_test.go (4 new cases)
