## MODIFIED Requirements

### Requirement: Docker config example includes Presidio service fields
The Docker deployment example config.json SHALL include the `presidio` block within `security.interceptor` so users deploying with `--profile presidio` have the correct default configuration.

#### Scenario: User deploys with Presidio profile
- **WHEN** a user runs `docker compose --profile presidio up` with the example config
- **THEN** the config.json already contains `presidio.enabled: false`, `presidio.url: "http://localhost:5002"`, `presidio.scoreThreshold: 0.7`, and `presidio.language: "en"`
- **THEN** the user only needs to set `presidio.enabled: true` to activate Presidio detection
