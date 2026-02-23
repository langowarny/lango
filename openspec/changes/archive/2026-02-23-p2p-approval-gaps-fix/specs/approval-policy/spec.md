## ADDED Requirements

### Requirement: Amount-based auto-approve for payment tools
wrapWithApproval SHALL accept an optional SpendingLimiter parameter (nil allowed). When non-nil and the tool is a payment tool (`p2p_pay` or `payment_send`), it SHALL check the amount parameter against `IsAutoApprovable` before requesting interactive approval.

#### Scenario: Auto-approve small payment
- **WHEN** tool is `p2p_pay` with amount "0.05" AND limiter.IsAutoApprovable returns true
- **THEN** the tool SHALL execute without interactive approval

#### Scenario: Require approval for large payment
- **WHEN** tool is `p2p_pay` with amount "5.00" AND limiter.IsAutoApprovable returns false
- **THEN** the tool SHALL request interactive approval via the approval provider

#### Scenario: No limiter provided
- **WHEN** limiter is nil
- **THEN** wrapWithApproval SHALL behave as before (no amount-based auto-approve)

#### Scenario: Non-payment tool unaffected
- **WHEN** tool is `exec` AND limiter is non-nil
- **THEN** wrapWithApproval SHALL ignore the limiter and follow normal approval policy

### Requirement: P2P payment approval summary
buildApprovalSummary SHALL return a human-readable summary for `p2p_pay` tool invocations including amount, peer DID (truncated), and memo.

#### Scenario: p2p_pay approval summary
- **WHEN** buildApprovalSummary is called with toolName "p2p_pay" and params containing amount, peer_did, and memo
- **THEN** it SHALL return a string containing the amount, truncated peer DID, and memo
