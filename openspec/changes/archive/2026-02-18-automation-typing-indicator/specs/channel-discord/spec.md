## ADDED Requirements

### Requirement: Public StartTyping with context support
The Discord channel SHALL expose a public `StartTyping(ctx context.Context, channelID string) func()` method that calls `ChannelTyping` and refreshes every 8 seconds until the stop function is called or the context is cancelled.

#### Scenario: Context cancellation stops typing
- **WHEN** `StartTyping` is called and the context is subsequently cancelled
- **THEN** the typing indicator goroutine SHALL exit without requiring the stop function to be called

#### Scenario: Stop function is idempotent
- **WHEN** the returned stop function is called multiple times
- **THEN** no panic SHALL occur (protected by `sync.Once`)

#### Scenario: Initial typing failure is non-blocking
- **WHEN** the initial `ChannelTyping` request fails
- **THEN** the error SHALL be logged at Warn level and a valid stop function SHALL still be returned
