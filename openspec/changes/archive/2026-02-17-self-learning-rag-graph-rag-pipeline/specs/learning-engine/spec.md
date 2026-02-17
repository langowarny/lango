## MODIFIED Requirements

### Requirement: Confidence propagation uses float64 math
The system SHALL apply fractional confidence boosts when propagating success across similar learnings. BoostLearningConfidence SHALL accept a `confidenceBoost float64` parameter; when > 0, it adds the value directly to confidence and clamps to [0.1, 1.0]. When 0, existing success/occurrence ratio calculation is used.

#### Scenario: Graph engine propagates fractional confidence
- **WHEN** a tool succeeds and similar learnings exist in the graph
- **THEN** each similar learning's confidence SHALL increase by `0.1 * propagationRate` (0.03 for rate 0.3)

#### Scenario: Base engine uses existing ratio calculation
- **WHEN** the base engine boosts confidence on tool success
- **THEN** it SHALL call BoostLearningConfidence with confidenceBoost=0.0, using success/occurrence ratio

#### Scenario: Confidence clamps to valid range
- **WHEN** a confidence boost would result in a value outside [0.1, 1.0]
- **THEN** the value SHALL be clamped to [0.1, 1.0]
