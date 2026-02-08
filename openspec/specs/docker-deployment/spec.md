## ADDED Requirements

### Requirement: Docker Container Configuration
The system SHALL provide a Dockerfile optimized for production deployment.

#### Scenario: Multi-stage build
- **WHEN** building the Docker image
- **THEN** the system SHALL use a multi-stage build
- **AND** the builder stage SHALL compile with CGO_ENABLED=0
- **AND** the runtime stage SHALL use debian:bookworm-slim

#### Scenario: Browser tool support
- **WHEN** Docker image is built
- **THEN** the runtime image SHALL include Chromium browser
- **AND** the ROD_BROWSER environment variable SHALL be set

#### Scenario: Non-root execution
- **WHEN** the container starts
- **THEN** the lango process SHALL run as non-root user (UID 1000)

#### Scenario: Health check
- **WHEN** the container is running
- **THEN** Docker SHALL perform health checks via HTTP endpoint
- **AND** unhealthy containers SHALL be marked for restart

### Requirement: Docker Compose Orchestration
The system SHALL provide a docker-compose.yml for simplified deployment.

#### Scenario: Service definition
- **WHEN** running `docker-compose up`
- **THEN** the lango service SHALL start on port 18789
- **AND** volumes SHALL persist data to lango-data volume

#### Scenario: Configuration mounting
- **WHEN** docker-compose starts the lango service
- **THEN** the system SHALL mount lango.json from host to /app/lango.json
- **AND** the config file SHALL be mounted read-only

#### Scenario: Environment variable injection
- **WHEN** docker-compose starts the lango service
- **THEN** environment variables (ANTHROPIC_API_KEY, DISCORD_BOT_TOKEN, TELEGRAM_BOT_TOKEN, SLACK_BOT_TOKEN, SLACK_APP_TOKEN) SHALL be passed to the container

### Requirement: Data Persistence
The system SHALL persist data across container restarts.

#### Scenario: SQLite database persistence
- **WHEN** the container restarts
- **THEN** the SQLite database at /data SHALL be preserved
- **AND** session history SHALL not be lost

#### Scenario: Volume mount
- **WHEN** docker-compose is used
- **THEN** a named volume (lango-data) SHALL be mounted at /data
