## Why

golangci-lint v2.4.0 CI 실행 시 90개 이슈(errcheck:50, staticcheck:28, unused:11, ineffassign:1)가 발생하여 CI가 통과하지 못했다. `.golangci.yml` 설정 파일이 없어 기본 설정으로 실행되었고, ent 자동생성 코드도 lint 대상에 포함되어 불필요한 이슈가 대량 보고되었다.

## What Changes

- Add `.golangci.yml` v2 configuration with `generated: strict` exclusion and `std-error-handling` preset
- Fix ~50 errcheck violations: unchecked `defer Close()`, `json.Encode`, `tx.Rollback`, `fmt.Scanln`, etc.
- Fix ~28 staticcheck issues: QF1012 (WriteString+Sprintf→Fprintf), S1009 (redundant nil check), S1011 (append spread), SA1012 (nil context), SA9003 (empty branches), QF1003 (if/else→switch), ST1005 (error string case), S1017 (redundant HasSuffix before TrimSuffix)
- Remove 11 unused declarations: functions, struct fields, variables, imports
- Fix 1 ineffassign: dead assignment removal

## Capabilities

### New Capabilities
- `lint-configuration`: golangci-lint v2 configuration (`.golangci.yml`) with generated code exclusion and standard presets

### Modified Capabilities

(No spec-level behavior changes - all modifications are code quality improvements that don't alter functionality)

## Impact

- 20+ files modified across `internal/`, `cmd/lango/`
- No API or behavioral changes - purely code quality improvements
- CI pipeline will pass cleanly with zero lint issues
- New `.golangci.yml` establishes project-wide linting standards
