## ADDED Requirements

### Requirement: Presidio connectivity check
The OutputScanningCheck SHALL verify Presidio connectivity when Presidio is enabled in config by checking the /health endpoint.

#### Scenario: Presidio enabled and reachable
- **WHEN** Presidio is enabled and the /health endpoint returns HTTP 200
- **THEN** the check SHALL return StatusPass

#### Scenario: Presidio enabled but unreachable
- **WHEN** Presidio is enabled but the endpoint is not reachable
- **THEN** the check SHALL return StatusWarn with a message suggesting the docker compose command

#### Scenario: Presidio not enabled
- **WHEN** Presidio is not enabled in config
- **THEN** the check SHALL not attempt Presidio connectivity verification
