## Requirements

### Requirement: Secret Isolation
The agent runtime SHALL have no direct access to API keys. The Supervisor SHALL hold provider credentials and proxy generation requests. This requirement is unchanged.

#### Scenario: Runtime Initialization
- **WHEN** the agent runtime is initialized
- **THEN** it SHALL receive a `ProviderProxy` that delegates to the Supervisor
- **THEN** API keys SHALL never be passed to the agent or tool execution environment

### Requirement: Supervisor Role
The Supervisor SHALL manage provider registry initialization and privileged exec tool execution. It SHALL NOT manage RPC crypto provider lifecycle. Crypto provider creation SHALL be handled optionally by the application layer only when security is explicitly configured.

#### Scenario: Bootstrapping
- **WHEN** the Supervisor is created
- **THEN** it SHALL initialize the provider registry and exec tool
- **THEN** it SHALL NOT initialize any crypto or security components

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
