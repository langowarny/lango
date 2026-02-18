## ADDED Requirements

### Requirement: Cron job typing indicator
The `Delivery` struct SHALL accept a `TypingIndicator` in addition to `ChannelSender`. The `NewDelivery` constructor SHALL accept `(sender ChannelSender, typing TypingIndicator, logger)`.

#### Scenario: Typing indicator during job execution
- **WHEN** a cron job is executed with delivery targets configured
- **THEN** the executor SHALL call `delivery.StartTyping(ctx, targets)` before `runner.Run()` and call the returned stop function after `runner.Run()` completes

#### Scenario: Typing indicator with nil typing
- **WHEN** `Delivery.StartTyping` is called but no `TypingIndicator` was provided
- **THEN** the method SHALL return a no-op stop function

#### Scenario: Typing indicator with empty targets
- **WHEN** `Delivery.StartTyping` is called with empty targets
- **THEN** the method SHALL return a no-op stop function

### Requirement: Multi-target typing aggregation
`Delivery.StartTyping` SHALL start typing on all provided targets and return a single stop function that stops all of them.

#### Scenario: Multiple delivery targets
- **WHEN** `StartTyping` is called with targets `["telegram:123", "discord:456"]`
- **THEN** typing indicators SHALL start on both channels and the returned stop function SHALL stop both
