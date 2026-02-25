## MODIFIED Requirements

### Requirement: P2P auto-approve respects SafetyLevel
The P2P approval function SHALL check the tool's SafetyLevel before applying price-based auto-approval. Dangerous tools (SafetyLevel == Dangerous or unknown/zero) MUST always go through explicit approval, regardless of price. Tools not found in the tool index SHALL be treated as dangerous.

#### Scenario: Dangerous tool bypasses auto-approve
- **WHEN** a P2P remote peer invokes a tool with SafetyLevel "dangerous" and the price is within auto-approve limits
- **THEN** the system SHALL NOT auto-approve and SHALL route to the composite approval provider

#### Scenario: Unknown tool treated as dangerous
- **WHEN** a P2P remote peer invokes a tool not found in the tool index
- **THEN** the system SHALL NOT auto-approve and SHALL route to the composite approval provider

#### Scenario: Safe tool within price limit auto-approves
- **WHEN** a P2P remote peer invokes a tool with SafetyLevel "safe" and the price is within auto-approve limits
- **THEN** the system SHALL auto-approve and record a grant
