## 1. Documentation & Example

- [x] 1.1 Update `lango.example.json` to reflect the new structure: remove `agent.apiKey` and populate `providers` with examples (including OAuth).

## 2. Configuration Structs

- [x] 2.1 Update `internal/config/types.go` to remove `APIKey` from `AgentConfig` struct.
- [x] 2.2 Verify `ProviderConfig` struct is sufficient (already done in previous change, but double-check).

## 3. Logic Updates

- [x] 3.1 Update `internal/supervisor/supervisor.go` to remove fallback logic that checks `agent.apiKey`.
- [x] 3.2 Update `Supervisor.initializeProviders` to ensure it only initializes providers defined in `cfg.Providers`.
- [x] 3.3 Ensure `agent.provider` references a valid key in `cfg.Providers` (or handle failure gracefully).

## 4. Verification

- [ ] 4.1 Run tests to ensure no regressions.
- [ ] 4.2 Manually verify `lango` fails if `lango.json` still uses `agent.apiKey` (due to unmarshalling error or validation).
