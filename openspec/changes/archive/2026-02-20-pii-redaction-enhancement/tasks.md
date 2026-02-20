## 1. Core Types & Pattern Catalog

- [x] 1.1 Create `internal/agent/pii_pattern.go` with PIICategory, PIIPatternDef, PIIMatch types and 13 builtin patterns
- [x] 1.2 Implement Luhn validation function for credit card post-match verification
- [x] 1.3 Create `internal/agent/pii_pattern_test.go` with regex validity, uniqueness, and per-pattern matching tests

## 2. PIIDetector Interface & Implementations

- [x] 2.1 Create `internal/agent/pii_detector.go` with PIIDetector interface, RegexDetector, and CompositeDetector
- [x] 2.2 Create `internal/agent/pii_detector_test.go` with detection, disable, custom pattern, and composite tests

## 3. Presidio Integration

- [x] 3.1 Create `internal/agent/pii_presidio.go` with PresidioDetector HTTP client and functional options
- [x] 3.2 Create `internal/agent/pii_presidio_test.go` with mock server tests for success, error, timeout, health check

## 4. PIIRedactor Refactor

- [x] 4.1 Refactor `internal/agent/pii_redactor.go` to use PIIDetector interface with position-based redaction
- [x] 4.2 Update `internal/agent/pii_redactor_test.go` with Korean PII, disabled builtins, custom patterns, Luhn tests

## 5. Config Changes

- [x] 5.1 Add PIIDisabledPatterns, PIICustomPatterns, PresidioConfig to InterceptorConfig in `internal/config/types.go`
- [x] 5.2 Add Presidio defaults to `internal/config/loader.go`

## 6. Wiring

- [x] 6.1 Update `internal/app/wiring.go` to pass new config fields to PIIRedactor

## 7. Settings TUI

- [x] 7.1 Add PII pattern management and Presidio fields to NewSecurityForm in `internal/cli/settings/forms_impl.go`
- [x] 7.2 Add formatCustomPatterns/ParseCustomPatterns helpers
- [x] 7.3 Add state update cases for new form keys in `internal/cli/tuicore/state_update.go`
- [x] 7.4 Update `internal/cli/settings/forms_impl_test.go` with new field keys and custom pattern parsing tests

## 8. Doctor Health Check

- [x] 8.1 Add Presidio connectivity check to OutputScanningCheck in `internal/cli/doctor/checks/output_scanning.go`

## 9. Docker Compose

- [x] 9.1 Add presidio-analyzer service with profile to `docker-compose.yml`
