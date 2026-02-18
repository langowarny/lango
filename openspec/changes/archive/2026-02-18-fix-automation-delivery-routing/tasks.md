## 1. Session Detection

- [x] 1.1 Update `detectChannelFromContext` in `internal/app/tools.go` to return `channel:targetID` format by splitting session key into 3 parts and joining first two
- [x] 1.2 Update tool parameter descriptions for `deliver_to` (cron_add) and `channel` (bg_submit) with format examples

## 2. Delivery Target Parsing and Routing

- [x] 2.1 Add `parseDeliveryTarget` helper in `internal/app/sender.go` to split `channel:id` into channel name and target ID
- [x] 2.2 Rewrite `SendMessage` to use parsed target ID for Telegram (with allowlist fallback), Discord, and Slack routing
- [x] 2.3 Add `strconv` import for Telegram chat ID parsing

## 3. Prompt Hints

- [x] 3.1 Update automation prompt section in `internal/app/wiring.go` with `channel:id` format hints for cron, background, and workflow descriptions

## 4. Verification

- [x] 4.1 Run `go build ./...` to verify compilation
- [x] 4.2 Run `go test ./internal/app/...` to verify tests pass
