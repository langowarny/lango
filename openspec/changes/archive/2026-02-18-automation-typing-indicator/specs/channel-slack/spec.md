## ADDED Requirements

### Requirement: Public StartTyping with placeholder pattern
The Slack channel SHALL expose a public `StartTyping(channelID string) func()` method that posts a `_Processing..._` placeholder message. The returned stop function SHALL delete the placeholder message.

#### Scenario: Successful placeholder lifecycle
- **WHEN** `StartTyping` is called and the stop function is subsequently called
- **THEN** the placeholder message SHALL be posted on start and deleted on stop

#### Scenario: Post failure returns no-op
- **WHEN** posting the placeholder message fails
- **THEN** the error SHALL be logged at Warn level and a no-op stop function SHALL be returned

#### Scenario: Stop function is idempotent
- **WHEN** the returned stop function is called multiple times
- **THEN** no panic SHALL occur (protected by `sync.Once`)

### Requirement: Slack Client interface includes DeleteMessage
The Slack `Client` interface SHALL include a `DeleteMessage(channelID, messageTimestamp string) (string, string, error)` method for placeholder cleanup.

#### Scenario: Mock clients implement DeleteMessage
- **WHEN** test mock clients implement the `Client` interface
- **THEN** they SHALL include a `DeleteMessage` method
