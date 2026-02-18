## ADDED Requirements

### Requirement: Makefile Docker Compose targets
The Makefile SHALL provide targets for managing Docker Compose profiles and containers.

#### Scenario: Start default profile
- **WHEN** running `make docker-up`
- **THEN** the system SHALL execute `docker compose --profile default up -d`

#### Scenario: Start browser profile
- **WHEN** running `make docker-up-browser`
- **THEN** the system SHALL execute `docker compose --profile browser up -d`

#### Scenario: Start browser-sidecar profile
- **WHEN** running `make docker-up-sidecar`
- **THEN** the system SHALL execute `docker compose --profile browser-sidecar up -d`

#### Scenario: Stop all containers
- **WHEN** running `make docker-down`
- **THEN** the system SHALL stop containers across all profiles (default, browser, browser-sidecar)

#### Scenario: Tail logs
- **WHEN** running `make docker-logs`
- **THEN** the system SHALL follow container logs via `docker compose logs -f`

### Requirement: Makefile Docker build variants
The Makefile SHALL provide targets for building Docker images with different configurations.

#### Scenario: Build with latest tag
- **WHEN** running `make docker-build`
- **THEN** the system SHALL tag the image with both the version tag and `latest`

#### Scenario: Build browser variant
- **WHEN** running `make docker-build-browser`
- **THEN** the system SHALL build with `WITH_BROWSER=true` and tag as `lango:browser`

#### Scenario: Push to registry
- **WHEN** running `make docker-push REGISTRY=my.registry.io`
- **THEN** the system SHALL tag and push both version and latest tags to the specified registry
- **WHEN** running `make docker-push` without REGISTRY set
- **THEN** the system SHALL fail with an error message
