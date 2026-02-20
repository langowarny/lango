## Why

Verify that archive works correctly with the fixed spec structure by adding a minor enhancement to the PII redaction system.

## What Changes

- Add a helper method to list all enabled builtin pattern names for diagnostic purposes.

## Capabilities

### New Capabilities

### Modified Capabilities
- `ai-privacy-interceptor`: Add diagnostic helper for listing active patterns

## Impact

- `internal/agent/pii_pattern.go` — one new helper function
