## Why

The P2P paid value exchange system has three critical gaps in its user approval (HITL) layer: the `autoApproveBelow` config field is defined and exposed in TUI settings but never wired into the `EntSpendingLimiter`, inbound remote tool invocations bypass owner approval entirely (only checking firewall ACL), and outbound `p2p_pay` payments require manual approval even for trivially small amounts that should be auto-approved.

## What Changes

- Add `autoApproveBelow` field to `EntSpendingLimiter` and expose `IsAutoApprovable()` on the `SpendingLimiter` interface for threshold-based auto-approval decisions.
- Add `ToolApprovalFunc` callback to the P2P protocol handler, inserting an owner approval gate between firewall ACL and tool execution for both free and paid inbound invocations.
- Integrate spending limiter into the `wrapWithApproval` tool wrapper so that outbound payment tools (`p2p_pay`, `payment_send`) auto-approve amounts below the configured threshold.
- Wire all components together in `app.go` and `wiring.go` — approval callback for inbound P2P, limiter for outbound approval wrapper.

## Capabilities

### New Capabilities

### Modified Capabilities
- `approval-policy`: `wrapWithApproval` gains a `SpendingLimiter` parameter for amount-based auto-approval of payment tools.
- `p2p-protocol`: Protocol handler gains `ToolApprovalFunc` callback for owner approval of inbound remote tool invocations.
- `blockchain-wallet`: `SpendingLimiter` interface adds `IsAutoApprovable` method; `EntSpendingLimiter` constructor accepts `autoApproveBelow` parameter.

## Impact

- **Core**: `internal/wallet/spending.go` — interface change (`IsAutoApprovable` added), constructor signature change (new parameter).
- **Core**: `internal/p2p/protocol/handler.go` — new callback type and approval gate in request handling.
- **Application**: `internal/app/tools.go` — `wrapWithApproval` signature change, auto-approve logic for payment tools.
- **Application**: `internal/app/app.go` — approval callback wiring for P2P handler, limiter passed to tool wrapper.
- **Application**: `internal/app/wiring.go` — `NewEntSpendingLimiter` call updated, `pricingFn` exposed on `p2pComponents`.
- **CLI**: `internal/cli/payment/payment.go` — `NewEntSpendingLimiter` call updated.
- **Tests**: `internal/wallet/spending_test.go` — new unit tests for `IsAutoApprovable`.
