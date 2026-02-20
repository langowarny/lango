## MODIFIED Requirements

### Requirement: Example config includes PII and Presidio fields
The example `config.json` SHALL include `piiDisabledPatterns` (empty array), `piiCustomPatterns` (empty object), and a `presidio` block with `enabled`, `url`, `scoreThreshold`, and `language` fields within the `security.interceptor` section.

#### Scenario: Docker headless user imports example config
- **WHEN** a user copies config.json for Docker headless deployment
- **THEN** the interceptor block contains `piiDisabledPatterns`, `piiCustomPatterns`, and `presidio` fields with sensible defaults
- **THEN** `presidio.enabled` defaults to `false` and `presidio.url` defaults to `http://localhost:5002`
