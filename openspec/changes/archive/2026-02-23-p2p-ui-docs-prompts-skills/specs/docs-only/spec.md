## MODIFIED Requirements

### Requirement: Documentation accuracy
Documentation, prompts, and CLI help text SHALL accurately reflect all implemented features including P2P REST API endpoints, CLI flags, and example projects.

#### Scenario: P2P REST API documented
- **WHEN** a user reads the HTTP API documentation
- **THEN** the P2P REST endpoints (`/api/p2p/status`, `/api/p2p/peers`, `/api/p2p/identity`) SHALL be documented with request/response examples

#### Scenario: Secrets --value-hex documented
- **WHEN** a user reads the secrets set CLI documentation
- **THEN** the `--value-hex` flag SHALL be documented with non-interactive usage examples

#### Scenario: P2P trading example discoverable
- **WHEN** a user reads the README
- **THEN** the `examples/p2p-trading/` directory SHALL be referenced in an Examples section
