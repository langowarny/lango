## ADDED Requirements

### Requirement: Background task typing indicator
The `Notification` struct SHALL accept a `TypingIndicator` in addition to `ChannelNotifier`. The `NewNotification` constructor SHALL accept `(notifier ChannelNotifier, typing TypingIndicator, logger)`.

#### Scenario: Typing indicator during task execution
- **WHEN** a background task is executed with an origin channel set
- **THEN** the manager SHALL call `notify.StartTyping(ctx, originChannel)` before `runner.Run()` and call the returned stop function after `runner.Run()` completes

#### Scenario: Typing indicator with nil typing
- **WHEN** `Notification.StartTyping` is called but no `TypingIndicator` was provided
- **THEN** the method SHALL return a no-op stop function

#### Scenario: Typing indicator with empty channel
- **WHEN** `Notification.StartTyping` is called with an empty channel
- **THEN** the method SHALL return a no-op stop function
