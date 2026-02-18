## MODIFIED Requirements

### Requirement: Cron job delivery channel resolution
The cron_add tool handler SHALL resolve delivery channels using the three-tier fallback chain: explicit deliver_to parameter → session auto-detection → cron.defaultDeliverTo config. The cron executor SHALL log a Warn-level message when a job completes with no delivery channel configured.

#### Scenario: Explicit deliver_to provided
- **WHEN** cron_add is called with a non-empty deliver_to array
- **THEN** the system SHALL use the provided channels without fallback

#### Scenario: Auto-detect from Telegram session
- **WHEN** cron_add is called without deliver_to AND the session key starts with "telegram:"
- **THEN** the system SHALL set deliver_to to ["telegram"]

#### Scenario: Config default used
- **WHEN** cron_add is called without deliver_to AND session auto-detection returns empty AND cron.defaultDeliverTo is configured
- **THEN** the system SHALL use the config default channels

#### Scenario: No delivery channel warning
- **WHEN** a cron job executes with empty DeliverTo
- **THEN** the executor SHALL log a Warn-level message including the job name and a configuration hint
