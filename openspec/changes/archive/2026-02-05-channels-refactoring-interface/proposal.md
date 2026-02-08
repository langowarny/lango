# Proposal: Refactor Channels for Testability

## Summary
Refactor `internal/channels` packages (`discord`, `slack`, `telegram`) to use interfaces for their respective external clients (`discordgo`, `slack-go`, `tgbotapi`). This enables standard mocking in tests.

## Problem Statement
Currently, the channel implementations depend directly on concrete structs from external libraries (`*discordgo.Session`, `*slack.Client`, `*tgbotapi.BotAPI`). This makes unit testing difficult because:
- We cannot easily mock the external dependencies.
- We have to rely on partial struct mocks or integration tests that require valid tokens or complex configuration.

## Proposed Solution
Introduce internal interface definitions for the subset of methods used by each channel implementation.
- **Discord**: Define `Session` interface and `DiscordSession` adapter.
- **Slack**: Define `Client` and `Socket` interfaces with adapters.
- **Telegram**: Define `BotAPI` interface and adapter.
- **Update Channels**: Modify `Channel` structs to use these interfaces.

## Capabilities

### New Capabilities
- `testability-interfaces`: internal interfaces for external channel libraries.

### Modified Capabilities
- `channel-discord`: Implementation refactored to use `Session` interface.
- `channel-slack`: Implementation refactored to use `Client` and `Socket` interfaces.
- `channel-telegram`: Implementation refactored to use `BotAPI` interface.

## Impact
- **Codebase**: low impact, internal refactoring of channel packages.
- **Dependencies**: No new external dependencies.
- **Testing**: significantly improved testability for channel logic.
