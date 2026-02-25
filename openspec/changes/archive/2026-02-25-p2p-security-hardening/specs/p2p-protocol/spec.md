## MODIFIED Requirements

### Requirement: Tool invocation approval check
The protocol handler SHALL deny tool invocation requests when no approval handler (`approvalFn`) is configured. The handler MUST return a response with status "denied" and error message "no approval handler configured for remote tool invocation". This applies to both free (`tool_invoke`) and paid (`tool_invoke_paid`) request types.

#### Scenario: No approval handler configured for tool_invoke
- **WHEN** a remote peer sends a `tool_invoke` request and `approvalFn` is nil
- **THEN** the handler SHALL return status "denied" with error "no approval handler configured for remote tool invocation"

#### Scenario: No approval handler configured for tool_invoke_paid
- **WHEN** a remote peer sends a `tool_invoke_paid` request and `approvalFn` is nil
- **THEN** the handler SHALL return status "denied" with error "no approval handler configured for remote tool invocation"

#### Scenario: Approval handler configured and approves
- **WHEN** a remote peer sends a `tool_invoke` request and `approvalFn` returns (true, nil)
- **THEN** the handler SHALL proceed to execute the tool and return status "ok"

#### Scenario: Approval handler configured and denies
- **WHEN** a remote peer sends a `tool_invoke` request and `approvalFn` returns (false, nil)
- **THEN** the handler SHALL return status "denied" with error "tool invocation denied by owner"

#### Scenario: Approval handler returns error
- **WHEN** a remote peer sends a `tool_invoke` request and `approvalFn` returns an error
- **THEN** the handler SHALL return status "error" with the approval error message
