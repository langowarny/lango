# Proposal: Remove Model Selection from Onboard Wizard

## Problem
The current onboarding wizard forces users to select an AI model from a hardcoded list. This list is difficult to maintain and quickly becomes outdated as providers release new models (e.g., Claude 3.6, Gemini 2.0 Flash). Maintaining this list in code creates unnecessary toil and friction.

## Solution
Remove the model selection step entirely from the `lango onboard` wizard. Instead, automatically configure a sensible default model based on the selected provider.

## Goals
- Simplify the onboarding flow by removing a step.
- Eliminate the maintenance burden of keeping model lists up-to-date in the CLI code.
- Ensure users still get a working configuration out of the box.

## Non-Goals
- Supporting exhaustive model selection in the CLI (users can edit `lango.json` manually if they need a specific non-default model).
