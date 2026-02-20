## MODIFIED Requirements

### Requirement: README documents PII detection capabilities
The README AI Privacy Interceptor section SHALL describe all 13 builtin PII detection patterns organized by category (Contact, Identity, Financial, Network). The section SHALL document pattern customization via `piiDisabledPatterns` and `piiCustomPatterns`. The section SHALL document optional Presidio NER-based detection integration.

#### Scenario: User reads AI Privacy Interceptor section
- **WHEN** a user reads the AI Privacy Interceptor section in README.md
- **THEN** they see the 4 pattern categories with specific pattern names listed
- **THEN** they see how to customize patterns (disable builtin, add custom)
- **THEN** they see how to enable Presidio integration with Docker Compose

### Requirement: README configuration table includes PII fields
The README configuration reference table SHALL include rows for `piiDisabledPatterns`, `piiCustomPatterns`, `presidio.enabled`, `presidio.url`, `presidio.scoreThreshold`, and `presidio.language` with correct types and defaults.

#### Scenario: User looks up PII config fields
- **WHEN** a user searches the configuration reference table for PII settings
- **THEN** they find 6 new rows after `piiRegexPatterns` with type, default, and description for each field
