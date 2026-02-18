## 1. Core Fix

- [x] 1.1 Add `LoadedSkills()` method to `Registry` in `internal/skill/registry.go` that returns only `r.loaded` (not `baseTools`)
- [x] 1.2 Change `internal/app/app.go` line 114 from `kc.registry.AllTools()` to `kc.registry.LoadedSkills()`

## 2. Tests

- [x] 2.1 Add `TestRegistry_LoadedSkills` test verifying empty result before load and only dynamic skills after activation
- [x] 2.2 Verify existing `TestRegistry_LoadSkills_AllTools` still passes (AllTools unchanged)

## 3. Verification

- [x] 3.1 Run `go build ./...` — build succeeds
- [x] 3.2 Run `go test ./...` — all tests pass
