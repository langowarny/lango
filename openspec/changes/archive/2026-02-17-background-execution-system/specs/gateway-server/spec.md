## ADDED Requirements

### Requirement: Turn completion callbacks
The gateway server SHALL support registering turn completion callbacks via OnTurnComplete() that fire after each agent turn.

#### Scenario: Register turn callback
- **WHEN** OnTurnComplete() is called with a callback function
- **THEN** the callback SHALL be appended to the server's turnCallbacks slice

#### Scenario: Fire turn callbacks after agent turn
- **WHEN** an agent turn completes in handleChatMessage (regardless of error)
- **THEN** all registered turn callbacks SHALL be invoked with the session key

#### Scenario: Multiple turn callbacks
- **WHEN** both MemoryBuffer.Trigger and AnalysisBuffer.Trigger are registered as callbacks
- **THEN** both SHALL fire after each agent turn
