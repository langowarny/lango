## REMOVED Requirements

### Requirement: Security Tools Registration
**Reason**: Security tools (secrets, crypto) are no longer registered by the agent runtime. Security is handled optionally at the application layer.
**Migration**: Tool registration moved to `internal/app/tools.go`. Security tools will be re-added in Phase 2 when security signer is configured.

### Requirement: Sensitive Tool Configuration
**Reason**: Approval-based tool gating removed from MVP. No sensitive tools are registered.
**Migration**: Re-introduce when security tools return in Phase 2.

## MODIFIED Requirements

### Requirement: Multi-Provider Support
The agent runtime SHALL delegate all LLM interactions to the ADK agent (`internal/adk/agent.go`). The legacy `Runtime.Run()` execution loop SHALL be removed. Type definitions (`Tool`, `ToolHandler`, `StreamEvent`, `Config`, `ParameterDef`, `AdkToolAdapter`) SHALL be retained for use by tool registration code.

#### Scenario: Provider initialization
- **WHEN** the application creates an ADK agent
- **THEN** it SHALL use `supervisor.NewProviderProxy()` and `adk.NewModelAdapter()` to bridge the provider to ADK
- **THEN** the `adk.Agent.Run()` method SHALL be the only execution path for generating responses

#### Scenario: Legacy runtime removal
- **WHEN** code references `agent.Runtime`
- **THEN** only type definitions and tool registration methods (`RegisterTool`, `GetTool`, `ListTools`, `ExecuteTool`) SHALL be available
- **THEN** the `Run()` method SHALL NOT exist
