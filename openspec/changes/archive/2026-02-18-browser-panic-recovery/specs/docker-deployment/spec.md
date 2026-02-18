## MODIFIED Requirements

### Requirement: Docker Compose Orchestration
The system SHALL provide a docker-compose.yml with deployment profiles for different browser configurations. The README documentation SHALL describe the import→delete configuration pattern instead of read-only JSON mounting.

#### Scenario: Service definition
- **WHEN** running any docker compose profile
- **THEN** the lango service SHALL expose port 18789
- **AND** volumes SHALL persist data to lango-data volume

#### Scenario: Slim deployment (default profile)
- **WHEN** running `docker compose up` or `docker compose --profile default up`
- **THEN** the lango service SHALL start without Chromium
- **AND** the image SHALL be approximately 200MB

#### Scenario: Built-in browser deployment (browser profile)
- **WHEN** running `docker compose --profile browser up`
- **THEN** the lango-browser service SHALL start with Chromium included in the image

#### Scenario: Sidecar browser deployment (browser-sidecar profile)
- **WHEN** running `docker compose --profile browser-sidecar up`
- **THEN** the lango-sidecar service SHALL start without Chromium
- **AND** a separate `chromedp/headless-shell` container SHALL run alongside
- **AND** lango SHALL connect to Chrome via `ROD_BROWSER_WS=ws://chrome:9222`
- **AND** the Chrome container SHALL have a memory limit of 512MB
- **AND** the Chrome container SHALL have a healthcheck that verifies CDP availability at `http://localhost:9222/json/version`
- **AND** the lango-sidecar service SHALL depend on the Chrome container being healthy (`service_healthy`)

#### Scenario: Optional prompts volume mount
- **WHEN** docker-compose.yml is inspected
- **THEN** it SHALL contain a commented-out volume mount for `./prompts:/usr/share/lango/prompts` to allow runtime prompt customization
- **AND** the default behavior (using embedded prompts) SHALL be unchanged when the comment is not removed

#### Scenario: Configuration via import
- **WHEN** docker-compose starts the lango service
- **THEN** the recommended configuration method is `lango config import` with auto-deletion of the source file
- **AND** `LANGO_PASSPHRASE` environment variable SHALL be used for non-interactive passphrase entry
