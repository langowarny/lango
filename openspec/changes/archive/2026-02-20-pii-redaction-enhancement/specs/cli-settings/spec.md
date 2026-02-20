## ADDED Requirements

### Requirement: Security form PII pattern fields
The Security configuration form SHALL include fields for managing PII patterns: disabled builtin patterns (comma-separated text), custom patterns (name:regex comma-separated text), Presidio enabled (bool), Presidio URL (text), and Presidio language (text).

#### Scenario: Disabled patterns field
- **WHEN** the Security form is created
- **THEN** it SHALL contain field with key "interceptor_pii_disabled"

#### Scenario: Custom patterns field
- **WHEN** the Security form is created with custom patterns {"a": "\\d+"}
- **THEN** it SHALL contain field with key "interceptor_pii_custom" showing "a:\\d+" format

#### Scenario: Presidio fields
- **WHEN** the Security form is created
- **THEN** it SHALL contain fields "presidio_enabled", "presidio_url", "presidio_language"

### Requirement: State update for PII fields
The ConfigState.UpdateConfigFromForm SHALL map the new PII form keys to their corresponding config fields.

#### Scenario: Update disabled patterns
- **WHEN** form field "interceptor_pii_disabled" has value "passport,ipv4"
- **THEN** config PIIDisabledPatterns SHALL be ["passport", "ipv4"]

#### Scenario: Update custom patterns
- **WHEN** form field "interceptor_pii_custom" has value "my_id:\\bID-\\d+\\b"
- **THEN** config PIICustomPatterns SHALL contain {"my_id": "\\bID-\\d+\\b"}

#### Scenario: Update Presidio enabled
- **WHEN** form field "presidio_enabled" is checked
- **THEN** config Presidio.Enabled SHALL be true
