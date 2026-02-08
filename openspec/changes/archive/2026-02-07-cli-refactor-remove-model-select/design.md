# Design: Remove Model Selection

## Overview
We will simplify the onboarding wizard by removing the explicit model selection step. The wizard will automatically select a sensible default model based on the chosen provider.

## Changes

### 1. Wizard State Machine
- Remove `StepModel` from the `WizardStep` enum.
- Update `handleEnter` transition:
    - FROM: `StepAPIKey` -> `StepModel` -> `StepChannel`
    - TO: `StepAPIKey` -> `StepChannel`

### 2. Default Model Logic
- When the user selects a provider (e.g., "openai"), the wizard will look up the `DefaultModel` from `ProviderMetadata` (e.g., "gpt-4o").
- This default model will be set in the `WizardConfig` automatically.

### 3. UI Updates
- Remove `viewModel` method and its associated view logic.
- Update `maxCursor` logic to remove `StepModel` case.

## Benefits
- **Reduced Maintenance**: No need to update the CLI code every time a new model version is released.
- **Improved UX**: Fewer steps for the user to get started.
- **Fail-Safe**: Defaults are guaranteed to work; power users can still edit `lango.json` for bleeding-edge models.
