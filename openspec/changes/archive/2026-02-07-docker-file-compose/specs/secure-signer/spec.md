## MODIFIED Requirements

### Requirement: Composite Provider Strategy
The system SHALL use a composite provider that tries companion first, then falls back to local.

#### Scenario: Companion available
- **WHEN** companion is connected
- **THEN** the system SHALL delegate crypto operations to companion via RPCProvider

#### Scenario: Companion unavailable with fallback
- **WHEN** companion is not connected
- **AND** local fallback is configured
- **AND** terminal is interactive (TTY available)
- **THEN** the system SHALL use local provider

#### Scenario: No providers available
- **WHEN** companion is not connected
- **AND** local fallback is not configured
- **THEN** the system SHALL return an error "no crypto provider available"

#### Scenario: Docker environment detection
- **WHEN** the system detects it is running in a Docker container (/.dockerenv exists OR cgroup contains "docker")
- **AND** no companion is connected
- **THEN** the system SHALL log error "Docker environment requires RPC Provider. Please connect Companion app."
- **AND** SHALL NOT attempt to use LocalCryptoProvider
