## ADDED Requirements

### Requirement: Conversation analyzer extracts knowledge from turns
The system SHALL analyze conversation turns using an LLM to extract facts, patterns, corrections, and preferences as knowledge and learning entries.

#### Scenario: Extract facts from conversation
- **WHEN** conversation analyzer processes a batch of messages containing factual statements
- **THEN** it SHALL produce KnowledgeEntry items with category matching the extracted type (fact, pattern, correction, preference)

#### Scenario: Extract corrections from user feedback
- **WHEN** a user corrects the agent's behavior in conversation
- **THEN** the analyzer SHALL produce a LearningEntry with category "user_correction" and confidence >= 0.7

### Requirement: Analysis triggers on turn count or token threshold
The system SHALL trigger conversation analysis when either the configured turn threshold (default 10) or token threshold (default 2000) is exceeded since the last analysis for that session.

#### Scenario: Turn threshold triggers analysis
- **WHEN** 10 new turns accumulate since last analysis
- **THEN** analysis SHALL be triggered for that session

#### Scenario: Token threshold triggers analysis
- **WHEN** 2000 tokens of new content accumulate since last analysis (even if fewer than 10 turns)
- **THEN** analysis SHALL be triggered for that session

### Requirement: Session learner extracts high-confidence knowledge at session end
The system SHALL analyze the full session at session end/pause to produce high-confidence knowledge entries.

#### Scenario: Skip short sessions
- **WHEN** a session ends with fewer than 4 turns
- **THEN** no session-end analysis SHALL occur

#### Scenario: Sample long sessions
- **WHEN** a session ends with more than 20 turns
- **THEN** the system SHALL sample first 3 + every 5th + last 5 messages for analysis

#### Scenario: Only store high-confidence results
- **WHEN** session analysis produces results
- **THEN** only entries with confidence == "high" SHALL be stored

### Requirement: Analysis buffer provides async processing
The system SHALL process conversation analysis asynchronously via a buffer with Start/Trigger/TriggerSessionEnd/Stop lifecycle.

#### Scenario: Buffer start and stop
- **WHEN** the application starts and knowledge+observational memory are enabled
- **THEN** the analysis buffer SHALL start a background goroutine and stop cleanly on shutdown

#### Scenario: Queue full drops gracefully
- **WHEN** the analysis queue is full
- **THEN** the request SHALL be dropped with a warn-level log message
