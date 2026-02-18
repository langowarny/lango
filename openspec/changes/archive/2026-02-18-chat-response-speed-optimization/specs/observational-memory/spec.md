## ADDED Requirements

### Requirement: Recent reflections retrieval
The memory Store SHALL provide a `ListRecentReflections` method that returns the N most recent reflections ordered by created_at ascending (chronological) for a given session.

#### Scenario: Retrieve limited recent reflections
- **WHEN** ListRecentReflections is called with limit N
- **THEN** the N most recent reflections SHALL be returned in chronological order (oldest first)

### Requirement: Recent observations retrieval
The memory Store SHALL provide a `ListRecentObservations` method that returns the N most recent observations ordered by created_at ascending (chronological) for a given session.

#### Scenario: Retrieve limited recent observations
- **WHEN** ListRecentObservations is called with limit N
- **THEN** the N most recent observations SHALL be returned in chronological order (oldest first)

### Requirement: Configurable memory context limits
The ObservationalMemoryConfig SHALL support `MaxReflectionsInContext` (default: 5) and `MaxObservationsInContext` (default: 20) fields to limit the number of reflections and observations injected into the LLM context. Zero means unlimited.

#### Scenario: Default limits applied
- **WHEN** maxReflectionsInContext and maxObservationsInContext are not configured (zero)
- **THEN** defaults of 5 reflections and 20 observations SHALL be applied by the wiring layer

#### Scenario: Custom limits applied
- **WHEN** maxReflectionsInContext is set to 3 and maxObservationsInContext is set to 10
- **THEN** only the 3 most recent reflections and 10 most recent observations SHALL be injected into context
