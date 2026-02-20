## ADDED Requirements

### Requirement: Extended PII pattern detection
PIIRedactor SHALL use a PIIDetector interface with 13 builtin regex patterns across contact, identity, financial, and network categories. PIIConfig SHALL support legacy fields (RedactEmail, RedactPhone, CustomRegex) and new fields (DisabledBuiltins, CustomPatterns, PresidioEnabled, PresidioURL, PresidioThreshold, PresidioLanguage).

#### Scenario: Korean mobile number redaction
- **WHEN** a user prompt contains "전화번호: 010-1234-5678"
- **THEN** the phone number SHALL be replaced with [REDACTED]

#### Scenario: Korean RRN redaction
- **WHEN** a user prompt contains "주민번호: 900101-1234567"
- **THEN** the RRN SHALL be replaced with [REDACTED]

#### Scenario: Disabled builtins
- **WHEN** PIIConfig has DisabledBuiltins=["email"]
- **THEN** PIIRedactor SHALL not detect email addresses

#### Scenario: Custom named patterns
- **WHEN** PIIConfig has CustomPatterns={"proj_id": "\\bPROJ-\\d{4}\\b"}
- **THEN** PIIRedactor SHALL detect matching text

#### Scenario: Presidio enabled
- **WHEN** PIIConfig has PresidioEnabled=true and PresidioURL set
- **THEN** PIIRedactor SHALL create a CompositeDetector with both RegexDetector and PresidioDetector

### Requirement: Position-based redaction
PIIRedactor.RedactInput SHALL use match position offsets to replace detected PII, merging overlapping matches into single [REDACTED] markers.

#### Scenario: Non-overlapping matches
- **WHEN** text contains email and phone at separate positions
- **THEN** each SHALL be replaced with [REDACTED] independently

#### Scenario: Overlapping matches
- **WHEN** two patterns match overlapping text regions
- **THEN** they SHALL be merged into a single [REDACTED] replacement

#### Scenario: No matches
- **WHEN** text contains no PII
- **THEN** RedactInput SHALL return the original text unchanged
