## ADDED Requirements

### Requirement: Presidio analyzer Docker service
The docker-compose.yml SHALL include a presidio-analyzer service using the mcr.microsoft.com/presidio-analyzer:latest image, exposed on port 5002, under the "presidio" profile.

#### Scenario: Profile-based activation
- **WHEN** user runs `docker compose --profile presidio up`
- **THEN** the presidio-analyzer container SHALL start alongside the main lango service

#### Scenario: Default compose up
- **WHEN** user runs `docker compose up` without the presidio profile
- **THEN** the presidio-analyzer container SHALL NOT start
