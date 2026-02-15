## ADDED Requirements

### Requirement: Headless auto-approve provider
The system SHALL provide a `HeadlessProvider` that implements the approval `Provider` interface and automatically approves all tool execution requests. `HeadlessProvider` SHALL be used only as a TTY fallback in the `CompositeProvider` and SHALL NOT be prefix-matched by session key (`CanHandle` returns false).

#### Scenario: Auto-approve in headless environment
- **WHEN** `HeadlessAutoApprove` is enabled in config
- **AND** no channel-specific provider or companion matches
- **THEN** the system SHALL auto-approve the tool execution
- **AND** the system SHALL log a WARN-level message including tool name, session key, and request ID

#### Scenario: Default disabled
- **WHEN** `HeadlessAutoApprove` is not set or is false
- **THEN** the system SHALL use `TTYProvider` as the fallback (existing behavior)
- **AND** headless auto-approve SHALL NOT be available

### Requirement: HeadlessAutoApprove config field
The `InterceptorConfig` SHALL include a `HeadlessAutoApprove` boolean field (default: false) that controls whether `HeadlessProvider` is wired as the TTY fallback.

#### Scenario: Config field recognition
- **WHEN** `security.interceptor.headlessAutoApprove: true` is set in configuration
- **THEN** the system SHALL wire `HeadlessProvider` as the TTY fallback in `CompositeProvider`
