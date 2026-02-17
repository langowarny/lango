## MODIFIED Requirements

### Requirement: Background task delivery channel resolution
The bg_submit tool handler SHALL resolve the delivery channel using the three-tier fallback chain: explicit channel parameter → session auto-detection → background.defaultDeliverTo config (first element). The notification system SHALL log a Warn-level message when a task completes with no origin channel.

#### Scenario: Explicit channel provided
- **WHEN** bg_submit is called with a non-empty channel parameter
- **THEN** the system SHALL use the provided channel without fallback

#### Scenario: Auto-detect from Discord session
- **WHEN** bg_submit is called without channel AND the session key starts with "discord:"
- **THEN** the system SHALL set channel to "discord"

#### Scenario: Config default used
- **WHEN** bg_submit is called without channel AND session auto-detection returns empty AND background.defaultDeliverTo is configured
- **THEN** the system SHALL use the first element of the config default

#### Scenario: No origin channel warning
- **WHEN** a background task notification is attempted with empty OriginChannel
- **THEN** the notification system SHALL log a Warn-level message with a configuration hint
