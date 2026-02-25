## 1. Default-Deny on Nil ApprovalFn

- [x] 1.1 Modify `handleToolInvoke()` in handler.go to return "denied" when `approvalFn` is nil
- [x] 1.2 Modify `handleToolInvokePaid()` in handler.go to return "denied" when `approvalFn` is nil
- [x] 1.3 Add handler_test.go with nil approval, approved, denied, and error test cases

## 2. P2P Approval Fallback Isolation

- [x] 2.1 Add `p2pFallback` field and `SetP2PFallback()` method to CompositeProvider
- [x] 2.2 Update `RequestApproval()` to route `"p2p:..."` sessions to P2P fallback instead of TTY fallback
- [x] 2.3 Wire `composite.SetP2PFallback(&approval.TTYProvider{})` in app.go when P2P is enabled
- [x] 2.4 Add tests for P2P session blocking HeadlessProvider, P2P fallback routing, non-P2P TTY routing

## 3. SafetyLevel Enforcement for P2P Auto-Approve

- [x] 3.1 Move approval func inside handler block to access `toolIndex`
- [x] 3.2 Add SafetyLevel check before price-based auto-approve in P2P approvalFn
- [x] 3.3 Record grants on approval success to prevent double-approval prompting

## 4. Firewall Wildcard Rule Validation

- [x] 4.1 Add `ValidateRule()` function to firewall.go
- [x] 4.2 Change `AddRule()` return type to `error` and call `ValidateRule()`
- [x] 4.3 Update `New()` to warn on overly permissive initial rules (backward compat)
- [x] 4.4 Update `p2p_firewall_add` tool handler in tools.go to handle `AddRule` error
- [x] 4.5 Add firewall_test.go with ValidateRule and AddRule test cases

## 5. Grant TTL Support

- [x] 5.1 Replace `map[string]struct{}` with `map[string]grantEntry{grantedAt}` in GrantStore
- [x] 5.2 Add `SetTTL()` and TTL expiration check in `IsGranted()`
- [x] 5.3 Add `CleanExpired()` method
- [x] 5.4 Wire `grantStore.SetTTL(time.Hour)` in app.go when P2P is enabled
- [x] 5.5 Add grant_test.go tests for TTL expired, TTL zero, CleanExpired

## 6. Build & Test Verification

- [x] 6.1 Verify `go build ./...` passes
- [x] 6.2 Verify `go test ./internal/approval/... ./internal/p2p/... ./internal/app/...` passes
