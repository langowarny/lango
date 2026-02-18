## MODIFIED Requirements

### Requirement: Context Assembly with Observations
The system SHALL include observations and reflections in the context sent to the LLM.

#### Scenario: Context ordering
- **WHEN** assembling the augmented system prompt
- **THEN** the system SHALL order context as: base prompt, knowledge RAG, reflections, observations, recent messages

#### Scenario: No observations available
- **WHEN** a session has no observations or reflections
- **THEN** the context assembly SHALL behave identically to the current system (no change)

#### Scenario: Observation formatting
- **WHEN** observations are included in the prompt
- **THEN** they SHALL be formatted under a "## Conversation Memory" section
- **AND** reflections SHALL appear before observations within that section

#### Scenario: Session key resolved from request context
- **WHEN** `GenerateContent` is called with a context containing a session key (set by gateway or channel adapter)
- **THEN** the adapter SHALL extract the session key via `session.SessionKeyFromContext(ctx)`
- **AND** use it for memory retrieval, RAG session filtering, and runtime context updates

#### Scenario: No session key in context skips memory
- **WHEN** `GenerateContent` is called with a context that has no session key
- **THEN** the adapter SHALL skip memory retrieval entirely
- **AND** the LLM response SHALL not include a Conversation Memory section

#### Scenario: Session key propagated to RuntimeContextAdapter
- **WHEN** a session key is resolved from context
- **THEN** the adapter SHALL call `RuntimeContextAdapter.SetSession(sessionKey)`
- **AND** the runtime context SHALL reflect the correct session key and derived channel type

#### Scenario: Session key propagated to RAG retrieval
- **WHEN** a session key is resolved from context and RAG is enabled
- **THEN** the adapter SHALL pass the session key to RAG/GraphRAG retrieval options
- **AND** results SHALL be filtered by the session key
