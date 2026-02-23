## ADDED Requirements

### Requirement: ToolApprovalFunc callback type
The protocol package SHALL define a `ToolApprovalFunc` callback type with signature `func(ctx context.Context, peerDID, toolName string, params map[string]interface{}) (bool, error)` that asks the local owner for approval before executing a remote tool invocation.

#### Scenario: Approval function defined
- **WHEN** the protocol package is compiled
- **THEN** ToolApprovalFunc type SHALL be available for use by callers

### Requirement: Handler owner approval for free tool invocations
Handler.handleToolInvoke SHALL check the approvalFn callback after firewall ACL check and before tool execution. If the approval function returns false, the handler SHALL return a "denied" response with error "tool invocation denied by owner".

#### Scenario: Approval granted for free tool
- **WHEN** a remote peer invokes a free tool AND approvalFn returns true
- **THEN** the tool SHALL execute normally

#### Scenario: Approval denied for free tool
- **WHEN** a remote peer invokes a free tool AND approvalFn returns false
- **THEN** the handler SHALL return status "denied" with error "tool invocation denied by owner"

#### Scenario: No approval function configured
- **WHEN** a remote peer invokes a tool AND approvalFn is nil
- **THEN** the tool SHALL execute without approval check (backward compatible)

### Requirement: Handler owner approval for paid tool invocations
Handler.handleToolInvokePaid SHALL check the approvalFn callback after payment verification and before tool execution.

#### Scenario: Approval granted for paid tool
- **WHEN** a remote peer invokes a paid tool with valid payment AND approvalFn returns true
- **THEN** the tool SHALL execute normally

#### Scenario: Approval denied for paid tool
- **WHEN** a remote peer invokes a paid tool with valid payment AND approvalFn returns false
- **THEN** the handler SHALL return status "denied" with error "tool invocation denied by owner"

### Requirement: SetApprovalFunc setter
Handler SHALL expose a `SetApprovalFunc(fn ToolApprovalFunc)` method to set the owner approval callback.

#### Scenario: SetApprovalFunc wires callback
- **WHEN** SetApprovalFunc is called with a non-nil function
- **THEN** subsequent tool invocations SHALL use the provided function for approval checks
