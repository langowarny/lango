## ADDED Requirements

### Requirement: Delivery channel fallback chain
The system SHALL resolve delivery channels using a three-tier fallback chain: (1) explicit parameter, (2) auto-detected session channel, (3) config default. If all tiers produce empty, the system SHALL log a warning and skip delivery without failing the operation.

#### Scenario: Explicit parameter takes precedence
- **WHEN** a cron job, background task, or workflow specifies an explicit delivery channel
- **THEN** the system SHALL use the explicit value and skip auto-detection and config fallback

#### Scenario: Auto-detect from session key
- **WHEN** no explicit delivery channel is provided AND the session key starts with a known channel prefix (telegram, discord, slack)
- **THEN** the system SHALL extract the channel name from the session key prefix and use it as the delivery target

#### Scenario: Config default fallback
- **WHEN** no explicit delivery channel is provided AND session auto-detection returns empty
- **THEN** the system SHALL use the configured `DefaultDeliverTo` value from the respective system config

#### Scenario: No delivery channel resolved
- **WHEN** all three fallback tiers produce no delivery channel
- **THEN** the system SHALL log a warning with the hint to configure a default and continue without delivering results

### Requirement: Session channel detection helper
The system SHALL provide a `detectChannelFromContext` helper that extracts the channel name from the session key stored in context. The helper SHALL return an empty string if the session key is empty or does not start with a known channel prefix.

#### Scenario: Known channel prefix
- **WHEN** the session key is "telegram:12345"
- **THEN** the helper SHALL return "telegram"

#### Scenario: Unknown prefix
- **WHEN** the session key is "api:user1" or "cron:jobname"
- **THEN** the helper SHALL return ""
