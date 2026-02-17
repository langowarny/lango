## MODIFIED Requirements

### Requirement: Session channel detection helper
The system SHALL provide a `detectChannelFromContext` helper that extracts the delivery target from the session key stored in context. The helper SHALL return a `channel:targetID` string (e.g. `telegram:123456789`) when the session key starts with a known channel prefix, or an empty string otherwise.

#### Scenario: Known channel prefix with target ID
- **WHEN** the session key is "telegram:123456789:987654321"
- **THEN** the helper SHALL return "telegram:123456789"

#### Scenario: Discord session key
- **WHEN** the session key is "discord:ch_abc:user_def"
- **THEN** the helper SHALL return "discord:ch_abc"

#### Scenario: Slack session key
- **WHEN** the session key is "slack:C1234567:U9876543"
- **THEN** the helper SHALL return "slack:C1234567"

#### Scenario: Unknown prefix
- **WHEN** the session key is "cron:jobname:123" or "api:user1"
- **THEN** the helper SHALL return ""

#### Scenario: Session key with fewer than 2 parts
- **WHEN** the session key has no colon separator
- **THEN** the helper SHALL return ""

## ADDED Requirements

### Requirement: Delivery target parsing
The system SHALL provide a `parseDeliveryTarget` helper that splits a delivery target string into channel name and optional target ID. The channel name SHALL be normalized to lowercase.

#### Scenario: Full target with ID
- **WHEN** the target is "telegram:123456789"
- **THEN** the parser SHALL return channelName="telegram" and targetID="123456789"

#### Scenario: Bare channel name
- **WHEN** the target is "telegram"
- **THEN** the parser SHALL return channelName="telegram" and targetID=""

#### Scenario: Negative chat ID
- **WHEN** the target is "telegram:-100123456"
- **THEN** the parser SHALL return channelName="telegram" and targetID="-100123456"

### Requirement: Telegram delivery with target ID routing
The system SHALL use the target ID from a parsed delivery target as the Telegram chat ID for message delivery. When no target ID is provided, the system SHALL fall back to the first allowlisted chat ID. When neither target ID nor allowlist is available, the system SHALL return an error.

#### Scenario: Target ID provided
- **WHEN** SendMessage receives target "telegram:123456789"
- **THEN** the system SHALL send the message to chat ID 123456789

#### Scenario: Bare channel name with allowlist
- **WHEN** SendMessage receives target "telegram" AND the allowlist contains at least one chat ID
- **THEN** the system SHALL send the message to the first allowlisted chat ID

#### Scenario: No target ID and empty allowlist
- **WHEN** SendMessage receives target "telegram" AND the allowlist is empty
- **THEN** the system SHALL return an error indicating a chat ID is required

### Requirement: Discord delivery with target ID routing
The system SHALL use the target ID from a parsed delivery target as the Discord channel ID for message delivery. When no target ID is provided, the system SHALL return an error.

#### Scenario: Target ID provided
- **WHEN** SendMessage receives target "discord:ch_abc"
- **THEN** the system SHALL send the message to Discord channel "ch_abc"

#### Scenario: No target ID
- **WHEN** SendMessage receives target "discord"
- **THEN** the system SHALL return an error indicating a channel ID is required

### Requirement: Slack delivery with target ID routing
The system SHALL use the target ID from a parsed delivery target as the Slack channel ID for message delivery. When no target ID is provided, the system SHALL return an error.

#### Scenario: Target ID provided
- **WHEN** SendMessage receives target "slack:C1234567"
- **THEN** the system SHALL send the message to Slack channel "C1234567"

#### Scenario: No target ID
- **WHEN** SendMessage receives target "slack"
- **THEN** the system SHALL return an error indicating a channel ID is required
