# Tasks: Remove Model Selection from Onboard Wizard

- [x] 1. Remove `StepModel` from `wizard.go` state machine.
- [x] 2. Update `handleEnter` to set default model from `ProviderMetadata` and skip to `StepChannel`.
- [x] 3. Remove `viewModel` and relevant view logic.
- [x] 4. Verify `onboard` flow.
