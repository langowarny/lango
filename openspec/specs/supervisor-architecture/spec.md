## ADDED Requirements

### Requirement: Secret Isolation
The Agent Runtime environment SHALL NOT contain any sensitive secrets (API Keys, Bot Tokens).

#### Scenario: Runtime Initialization
- **WHEN** the Agent Runtime is initialized
- **THEN** it SHALL NOT receive API keys in its configuration
- **AND** it SHALL NOT inherit environment variables containing secrets

### Requirement: Supervisor Role
A Supervisor component SHALL be responsible for managing the lifecycle of the Agent and holding all secrets.

#### Scenario: Bootstrapping
- **WHEN** the application starts
- **THEN** the Supervisor SHALL be initialized first with full configuration
- **AND** the Supervisor SHALL initialize the Agent Runtime

### Requirement: Provider Proxy
The Agent Runtime SHALL use a proxy mechanism to request AI generation from the Supervisor.

#### Scenario: Generation Request
- **WHEN** the Agent needs to generate text or call tools
- **THEN** it SHALL call a Provider interface
- **AND** this interface SHALL forward the request to the Supervisor
- **AND** the Supervisor SHALL execute the request using the real Provider Client (with keys)

### Requirement: Privileged Tool Execution
Sensitive tools (such as `exec`) SHALL be executed by the Supervisor to enforce security policies.

#### Scenario: Exec Tool Usage
- **WHEN** the Agent invokes the `exec` tool
- **THEN** the Runtime SHALL forward the execution request to the Supervisor
- **AND** the Supervisor SHALL validate the command and environment
- **AND** the Supervisor SHALL execute the command
