## 1. Channel Adapter Public StartTyping

- [x] 1.1 Add `StartTyping(ctx, chatID) func()` to `internal/channels/telegram/telegram.go` with ctx.Done() monitoring and sync.Once stop
- [x] 1.2 Add `StartTyping(ctx, channelID) func()` to `internal/channels/discord/discord.go` with ctx.Done() monitoring and sync.Once stop
- [x] 1.3 Add `DeleteMessage` to Slack `Client` interface in `internal/channels/slack/client.go`
- [x] 1.4 Add `StartTyping(channelID) func()` to `internal/channels/slack/slack.go` with placeholder post/delete pattern and sync.Once stop
- [x] 1.5 Add `DeleteMessage` to mock clients in Slack test files (`slack_test.go`)

## 2. Channel Sender Typing Dispatch

- [x] 2.1 Add `StartTyping(ctx, channel) (func(), error)` to `channelSender` in `internal/app/sender.go`

## 3. Cron Typing Integration

- [x] 3.1 Add `TypingIndicator` interface to `internal/cron/delivery.go`
- [x] 3.2 Add `typing` field to `Delivery` struct and update `NewDelivery` constructor signature
- [x] 3.3 Add `Delivery.StartTyping(ctx, targets) func()` method for multi-target typing aggregation
- [x] 3.4 Call `delivery.StartTyping` / `stopTyping()` around `runner.Run()` in `internal/cron/executor.go`

## 4. Background Typing Integration

- [x] 4.1 Add `TypingIndicator` interface to `internal/background/notification.go`
- [x] 4.2 Add `typing` field to `Notification` struct and update `NewNotification` constructor signature
- [x] 4.3 Add `Notification.StartTyping(ctx, channel) func()` method
- [x] 4.4 Call `notify.StartTyping` / `stopTyping()` around `runner.Run()` in `internal/background/manager.go`

## 5. Wiring

- [x] 5.1 Update `initCron` in `internal/app/wiring.go` to pass `sender` as typing indicator to `NewDelivery`
- [x] 5.2 Update `initBackground` in `internal/app/wiring.go` to pass `sender` as typing indicator to `NewNotification`

## 6. Verification

- [x] 6.1 Run `go build ./...` to verify compilation
- [x] 6.2 Run tests for all modified packages
