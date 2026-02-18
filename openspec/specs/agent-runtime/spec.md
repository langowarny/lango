## Requirements

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

### Requirement: Model Fallback execution
The system SHALL execute model fallbacks when the primary provider fails.

#### Scenario: Primary provider failure
- **WHEN** the primary provider fails with a retryable error (e.g., rate limit, overload)
- **THEN** the system SHALL attempt to use the next configured fallback provider
- **AND** SHALL log the fallback event

#### Scenario: All providers fail
- **WHEN** all configured providers (primary + fallbacks) fail
- **THEN** the system SHALL return an error to the user
