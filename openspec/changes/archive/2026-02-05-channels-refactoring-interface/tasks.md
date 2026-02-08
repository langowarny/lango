## 1. Discord Refactoring

- [x] 1.1 Define `Session` interface and `DiscordSession` adapter in `internal/channels/discord/session.go`.
- [x] 1.2 Update `Channel` struct in `internal/channels/discord/discord.go` to use `Session` interface.
- [x] 1.3 Refactor `discord.go` methods to use interface methods (including `GetState`).
- [x] 1.4 Update or creation `discord_test.go` using a mock implementation of `Session`.

## 2. Slack Refactoring

- [x] 2.1 Define `Client` and `Socket` interfaces and adapters in `internal/channels/slack/client.go`.
- [x] 2.2 Update `Channel` struct in `internal/channels/slack/slack.go` to use newly defined interfaces.
- [x] 2.3 Refactor `slack.go` to use interface methods (including `Events()`).
- [x] 2.4 Update `slack_test.go` using mock implementations.

## 3. Telegram Refactoring

- [x] 3.1 Define `BotAPI` interface and adapter in `internal/channels/telegram/bot.go`.
- [x] 3.2 Update `Channel` struct in `internal/channels/telegram/telegram.go` to use `BotAPI` interface.
- [x] 3.3 Refactor `telegram.go` to use interface methods (including `GetSelf`).
- [x] 3.4 Update `telegram_test.go` using mock implementation.
