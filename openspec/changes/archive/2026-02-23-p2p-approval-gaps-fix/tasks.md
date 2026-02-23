## 1. SpendingLimiter autoApproveBelow Wiring

- [x] 1.1 Add `autoApproveBelow *big.Int` field to `EntSpendingLimiter` struct
- [x] 1.2 Add `autoApproveBelow string` parameter to `NewEntSpendingLimiter` constructor
- [x] 1.3 Add `IsAutoApprovable(ctx, amount) (bool, error)` to `SpendingLimiter` interface
- [x] 1.4 Implement `IsAutoApprovable` on `EntSpendingLimiter`
- [x] 1.5 Update `NewEntSpendingLimiter` call in `internal/app/wiring.go` with `cfg.Payment.Limits.AutoApproveBelow`
- [x] 1.6 Update `NewEntSpendingLimiter` call in `internal/cli/payment/payment.go` with `cfg.Payment.Limits.AutoApproveBelow`
- [x] 1.7 Add unit tests for `IsAutoApprovable` and `NewEntSpendingLimiter` autoApproveBelow parsing

## 2. Inbound P2P Tool Approval Layer

- [x] 2.1 Define `ToolApprovalFunc` callback type in `internal/p2p/protocol/handler.go`
- [x] 2.2 Add `approvalFn ToolApprovalFunc` field to `Handler` struct
- [x] 2.3 Add `SetApprovalFunc(fn ToolApprovalFunc)` setter method
- [x] 2.4 Insert approval check in `handleToolInvoke` after firewall, before executor
- [x] 2.5 Insert approval check in `handleToolInvokePaid` after payment verification, before executor
- [x] 2.6 Add `pricingFn` field to `p2pComponents` struct and populate in `initP2P` return

## 3. Outbound Payment Auto-Approve Integration

- [x] 3.1 Add `limiter wallet.SpendingLimiter` parameter to `wrapWithApproval` function
- [x] 3.2 Add amount-based auto-approve logic for `p2p_pay` and `payment_send` tools
- [x] 3.3 Add `p2p_pay` case to `buildApprovalSummary`
- [x] 3.4 Add `wallet` import to `internal/app/tools.go`

## 4. Application Wiring

- [x] 4.1 Pass `pc.limiter` to `wrapWithApproval` in `app.go` (nil when payment disabled)
- [x] 4.2 Wire `SetApprovalFunc` on P2P handler in `app.go` with pricingFn + limiter + composite approval
- [x] 4.3 Add `wallet` and `time` imports to `app.go`

## 5. Verification

- [x] 5.1 `go build ./...` passes
- [x] 5.2 `go test ./internal/wallet/...` passes (IsAutoApprovable tests)
- [x] 5.3 `go test ./internal/p2p/...` passes
