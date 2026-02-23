## Context

The P2P paid value exchange system was introduced with `autoApproveBelow` as a config field (default "0.10" USDC) exposed in TUI settings, but the `EntSpendingLimiter` never reads or uses it. Inbound P2P tool invocations pass through firewall ACL only — no owner approval gate exists. Outbound `p2p_pay` uses `SafetyLevelDangerous` which triggers `wrapWithApproval`, but has no awareness of the auto-approve threshold.

## Goals / Non-Goals

**Goals:**
- Wire `autoApproveBelow` from config into `EntSpendingLimiter` so the threshold is enforced.
- Add owner approval gate for inbound P2P tool invocations (both free and paid) via callback pattern.
- Enable amount-based auto-approval for outbound payment tools in `wrapWithApproval`.
- Maintain fail-closed semantics: deny by default when approval provider is unavailable.

**Non-Goals:**
- Changing the approval UI/UX in TUI or Gateway WebSocket.
- Adding new config fields beyond using the existing `autoApproveBelow`.
- Modifying the P2P firewall ACL system or reputation scoring.
- Per-peer approval granularity (future enhancement).

## Decisions

### 1. `IsAutoApprovable` on `SpendingLimiter` interface (not standalone function)

Adding `IsAutoApprovable(ctx, amount) (bool, error)` to the `SpendingLimiter` interface keeps the auto-approve decision co-located with the spending limit check. The method composes threshold check + limit check atomically, preventing race conditions where a threshold-passing amount could exceed daily limits.

**Alternative**: Standalone function taking threshold + limiter — rejected because it splits the decision across two call sites and requires callers to parse the threshold themselves.

### 2. Callback pattern for inbound approval (`ToolApprovalFunc`)

Using `ToolApprovalFunc func(ctx, peerDID, toolName, params) (bool, error)` as a callback on the Handler avoids import cycles between `p2p/protocol` and `approval`/`wallet` packages. The closure is wired in `app.go` where all dependencies are available.

**Alternative**: Direct dependency on approval package — rejected due to import cycle with `app` package.

### 3. Limiter parameter on `wrapWithApproval` (nil-safe)

Adding `limiter wallet.SpendingLimiter` as a parameter (nil allowed) to `wrapWithApproval` keeps the auto-approve logic centralized in the approval wrapper. When nil, behavior is unchanged. When non-nil, payment tools (`p2p_pay`, `payment_send`) extract the amount parameter and check `IsAutoApprovable` before falling through to interactive approval.

**Alternative**: Separate wrapper function — rejected because it would require re-implementing all of `wrapWithApproval`'s grant logic.

### 4. Inbound approval uses pricing function to determine auto-approvability

For inbound paid tool invocations, the approval callback uses `pricingFn` to look up the tool's price, then checks `IsAutoApprovable`. This means the owner auto-approves based on the price they set, not the amount the peer pays. For free tools, the approval always goes to the interactive provider.

## Risks / Trade-offs

- **[Auto-approve threshold too high]** → Users can set it to "0" to disable. Default "0.10" USDC is conservative.
- **[Inbound approval blocks P2P latency]** → TTY/Gateway approval is synchronous, adding latency to remote tool calls. Mitigated by auto-approve for small paid tools.
- **[Breaking change: `NewEntSpendingLimiter` signature]** → Only two call sites (wiring.go, cli/payment.go), both updated. Not a public API.
- **[Breaking change: `SpendingLimiter` interface]** → Adding `IsAutoApprovable` is a breaking interface change. Only one implementation exists (`EntSpendingLimiter`), so impact is contained.
