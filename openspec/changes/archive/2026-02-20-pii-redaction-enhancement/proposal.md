## Why

The current PII redaction system only supports 2 regex patterns (email and US phone), making it ineffective for real-world personal data protection. Korean PII (resident registration numbers, Korean phone numbers) is not detected at all, users cannot manage patterns through the Settings TUI, and there is no integration with advanced NER-based detection services like Microsoft Presidio.

## What Changes

- Expand builtin PII patterns from 2 to 13 (Korean mobile/landline, Korean RRN, US SSN, credit cards with Luhn validation, IBAN, passport, IPv4, etc.)
- Introduce a `PIIDetector` interface with `RegexDetector`, `CompositeDetector`, and `PresidioDetector` implementations
- Refactor `PIIRedactor` to use the `PIIDetector` interface instead of managing regexes directly
- Add match-position-based redaction to handle overlapping patterns correctly
- Add config fields for disabling builtin patterns, adding custom named patterns, and Presidio integration
- Add PII pattern management fields to the Settings TUI Security form
- Add Presidio health check to the doctor system
- Add Presidio analyzer service to Docker Compose (via profile)

## Capabilities

### New Capabilities
- `pii-pattern-catalog`: Builtin PII pattern definitions with categories, labels, enable/disable defaults, and optional post-match validation
- `pii-detector-interface`: PIIDetector interface with RegexDetector, CompositeDetector, and PresidioDetector implementations
- `presidio-integration`: Microsoft Presidio HTTP client for NER-based PII detection with graceful degradation

### Modified Capabilities


## Impact

- **Core**: `internal/agent/` — 3 new files, 1 refactored file
- **Config**: `internal/config/types.go`, `internal/config/loader.go` — new fields and defaults
- **Wiring**: `internal/app/wiring.go` — passes new config fields to PIIRedactor
- **TUI**: `internal/cli/settings/forms_impl.go`, `internal/cli/tuicore/state_update.go` — new form fields and state mapping
- **Doctor**: `internal/cli/doctor/checks/output_scanning.go` — Presidio connectivity check
- **Docker**: `docker-compose.yml` — new Presidio service
- **Dependencies**: No new Go module dependencies (Presidio uses standard `net/http`)
- **Backward Compatibility**: Fully backward compatible — legacy `RedactEmail`/`RedactPhone` bools and `CustomRegex` still work
