## MODIFIED Requirements

### Requirement: Approval callback unblocks agent before UI update
`HandleCallback` SHALL send the approval result to the waiting agent's channel BEFORE calling `editApprovalMessage` to update the Telegram message. This ensures the agent pipeline is not blocked by Telegram API latency.

#### Scenario: Approval granted
- **WHEN** a user clicks "Approve" on an inline keyboard
- **THEN** the approval result (true) is sent to the agent's channel immediately, and THEN the Telegram message is edited to show "Approved" status

#### Scenario: Approval denied
- **WHEN** a user clicks "Deny" on an inline keyboard
- **THEN** the denial result (false) is sent to the agent's channel immediately, and THEN the Telegram message is edited to show "Denied" status

#### Scenario: Multiple consecutive approvals
- **WHEN** 4 tools require consecutive approval and all are approved
- **THEN** the agent processes each approval without cumulative Telegram API latency between them
