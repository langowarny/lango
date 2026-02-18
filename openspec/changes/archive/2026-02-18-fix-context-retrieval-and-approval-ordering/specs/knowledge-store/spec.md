## MODIFIED Requirements

### Requirement: Knowledge search uses per-keyword predicates
`SearchKnowledge` SHALL split the query into individual keywords and create separate `ContentContains`/`KeyContains` LIKE predicates for each keyword, combined with OR logic. The system SHALL NOT use a single concatenated query string as a LIKE pattern.

#### Scenario: Multi-keyword search
- **WHEN** `SearchKnowledge` is called with query "deploy server config"
- **THEN** the SQL query uses `(content LIKE '%deploy%' OR key LIKE '%deploy%') OR (content LIKE '%server%' OR key LIKE '%server%') OR (content LIKE '%config%' OR key LIKE '%config%')`

#### Scenario: Single keyword search
- **WHEN** `SearchKnowledge` is called with query "deploy"
- **THEN** the SQL query uses `content LIKE '%deploy%' OR key LIKE '%deploy%'`

### Requirement: Learning search uses per-keyword predicates
`SearchLearnings` and `SearchLearningEntities` SHALL split the query into individual keywords and create separate `ErrorPatternContains`/`TriggerContains` LIKE predicates for each keyword, combined with OR logic.

#### Scenario: Multi-keyword learning search
- **WHEN** `SearchLearnings` is called with error pattern "connection timeout retry"
- **THEN** each keyword generates individual LIKE predicates combined with OR

### Requirement: External ref search uses per-keyword predicates
`SearchExternalRefs` SHALL split the query into individual keywords and create separate `NameContains`/`SummaryContains` LIKE predicates for each keyword, combined with OR logic.

#### Scenario: Multi-keyword external ref search
- **WHEN** `SearchExternalRefs` is called with query "api documentation"
- **THEN** each keyword generates individual LIKE predicates combined with OR
