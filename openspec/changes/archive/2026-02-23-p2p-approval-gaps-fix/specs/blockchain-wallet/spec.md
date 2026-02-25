## ADDED Requirements

### Requirement: SpendingLimiter auto-approve threshold
SpendingLimiter interface SHALL include `IsAutoApprovable(ctx context.Context, amount *big.Int) (bool, error)` for threshold-based auto-approval decisions. EntSpendingLimiter SHALL implement this method using the `autoApproveBelow` field.

#### Scenario: IsAutoApprovable returns true for amount below threshold
- **WHEN** autoApproveBelow is "0.10" and amount is 0.05 USDC and daily limit is not exceeded
- **THEN** IsAutoApprovable SHALL return (true, nil)

#### Scenario: IsAutoApprovable returns false for amount above threshold
- **WHEN** autoApproveBelow is "0.10" and amount is 0.50 USDC
- **THEN** IsAutoApprovable SHALL return (false, nil)

#### Scenario: IsAutoApprovable returns false when threshold is zero
- **WHEN** autoApproveBelow is "0" or empty
- **THEN** IsAutoApprovable SHALL return (false, nil) regardless of amount

#### Scenario: IsAutoApprovable returns error when daily limit exceeded
- **WHEN** amount is below threshold but daily spending limit would be exceeded
- **THEN** IsAutoApprovable SHALL return (false, error) with the limit error

### Requirement: EntSpendingLimiter autoApproveBelow parameter
NewEntSpendingLimiter SHALL accept an `autoApproveBelow` string parameter (4th argument) representing the USDC amount threshold for auto-approval. Empty string or "0" SHALL disable auto-approval.

#### Scenario: Valid autoApproveBelow value
- **WHEN** NewEntSpendingLimiter is called with autoApproveBelow "0.10"
- **THEN** the limiter SHALL store the parsed threshold as 100000 (smallest USDC units)

#### Scenario: Empty autoApproveBelow disables auto-approval
- **WHEN** NewEntSpendingLimiter is called with autoApproveBelow ""
- **THEN** the limiter SHALL set autoApproveBelow to 0 (disabled)

#### Scenario: Invalid autoApproveBelow returns error
- **WHEN** NewEntSpendingLimiter is called with autoApproveBelow "invalid"
- **THEN** NewEntSpendingLimiter SHALL return an error
