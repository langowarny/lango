## Context

The current PII redaction system (`internal/agent/pii_redactor.go`) directly compiles and manages regex patterns. It supports only email and US phone patterns, with optional custom regex via `PIIConfig.CustomRegex`. The redactor is wired into the LLM pipeline via `PIIRedactingModelAdapter` in `internal/adk/pii_model.go`, which calls `RedactInput()` on user messages.

The system needs to support Korean PII (RRN, mobile/landline numbers), financial data (credit cards with Luhn validation), and optionally NER-based detection via Microsoft Presidio.

## Goals / Non-Goals

**Goals:**
- Expand builtin PII patterns to 13 (contact, identity, financial, network categories)
- Introduce `PIIDetector` interface for pluggable detection backends
- Support per-pattern enable/disable and custom named patterns via config
- Integrate Microsoft Presidio for NER-based detection with graceful degradation
- Add TUI settings for PII pattern management
- Maintain full backward compatibility with existing config

**Non-Goals:**
- Real-time PII detection in streaming responses (existing batch approach is retained)
- Custom Presidio recognizer configuration (uses Presidio's built-in recognizers)
- PII detection in non-text content (images, audio)
- Anonymization modes beyond `[REDACTED]` replacement

## Decisions

### 1. PIIDetector Interface Pattern
**Decision**: Introduce `PIIDetector` interface with `Detect(text string) []PIIMatch` method.
**Rationale**: Decouples detection logic from redaction, enables pluggable backends (regex, Presidio, future ML-based). Follows the existing adapter pattern used throughout the codebase (e.g., `ContextAwareModelAdapter`).
**Alternative**: Extend existing regex-only approach — rejected because it can't support NER-based detection.

### 2. CompositeDetector for Multi-Backend
**Decision**: Use `CompositeDetector` to chain `RegexDetector` + `PresidioDetector`, deduplicating overlapping matches by score.
**Rationale**: Regex provides fast, reliable detection for known patterns. Presidio adds NER-based detection for context-dependent PII. Composite approach ensures regex always works as fallback.
**Alternative**: Single detector with mixed strategies — rejected for separation of concerns.

### 3. Position-Based Redaction
**Decision**: Replace sequential `ReplaceAllString` with position-based replacement using `PIIMatch` start/end offsets.
**Rationale**: Sequential replacement can corrupt positions when multiple patterns match overlapping or adjacent regions. Position-based approach merges overlapping matches first, then replaces in a single pass.

### 4. Presidio as Optional Docker Service
**Decision**: Presidio runs as a Docker Compose service under the `presidio` profile.
**Rationale**: Presidio requires Python runtime and model downloads. Docker profile keeps it optional — `docker compose --profile presidio up` to enable. Graceful degradation: if Presidio is unreachable, regex detection continues normally.

### 5. Luhn Validation for Credit Cards
**Decision**: Add optional `Validate func(string) bool` to `PIIPatternDef` for post-match validation.
**Rationale**: Credit card regex matches format but cannot verify check digits. Luhn validation eliminates false positives from random digit sequences matching the credit card format.

### 6. Legacy Backward Compatibility
**Decision**: Retain `RedactEmail`, `RedactPhone`, and `CustomRegex` fields in `PIIConfig`.
**Rationale**: Existing configurations and wiring code continue to work. Legacy bools map to enabling/disabling the `email` and `us_phone` builtin patterns.

## Risks / Trade-offs

- **[Risk] Korean landline regex may match non-phone numbers** → Pattern is disabled by default; users opt-in via config
- **[Risk] Presidio adds latency to LLM pipeline** → Presidio runs async; timeout set to 5s with graceful degradation to regex-only
- **[Risk] Credit card regex may conflict with other numeric patterns (SSN, phone)** → CompositeDetector deduplicates overlapping matches by position, preferring higher-score matches
- **[Trade-off] `kr_bank_account` pattern disabled by default** → High false-positive rate due to overlap with phone number formats; users enable explicitly
