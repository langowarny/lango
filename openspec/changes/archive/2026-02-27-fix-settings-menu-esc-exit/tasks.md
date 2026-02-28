## 1. Editor Esc Behavior Fix

- [x] 1.1 Change Esc at StepMenu in `editor.go` from `tea.Quit` to `e.step = StepWelcome` with `nil` cmd return
- [x] 1.2 Update help bar text in `menu.go` from `"Quit"` to `"Back"` for the Esc key entry

## 2. Tests

- [x] 2.1 Create `editor_test.go` with test: Esc at StepWelcome triggers quit
- [x] 2.2 Add test: Esc at StepMenu navigates to StepWelcome (no quit)
- [x] 2.3 Add test: Esc at StepMenu while searching stays at StepMenu (search cancelled)
- [x] 2.4 Add test: Ctrl+C at all steps triggers quit with Cancelled flag

## 3. Verification

- [x] 3.1 Run `go build ./...` — no compilation errors
- [x] 3.2 Run `go test ./internal/cli/settings/...` — all tests pass
