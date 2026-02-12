## MODIFIED Requirements

### Requirement: Supervisor Role
The Supervisor SHALL manage provider registry initialization and privileged exec tool execution. It SHALL NOT manage RPC crypto provider lifecycle. Crypto provider creation SHALL be handled optionally by the application layer only when security is explicitly configured.

#### Scenario: Bootstrapping
- **WHEN** the Supervisor is created
- **THEN** it SHALL initialize the provider registry and exec tool
- **THEN** it SHALL NOT initialize any crypto or security components

### Requirement: Secret Isolation
The agent runtime SHALL have no direct access to API keys. The Supervisor SHALL hold provider credentials and proxy generation requests. This requirement is unchanged.

#### Scenario: Runtime Initialization
- **WHEN** the agent runtime is initialized
- **THEN** it SHALL receive a `ProviderProxy` that delegates to the Supervisor
- **THEN** API keys SHALL never be passed to the agent or tool execution environment
