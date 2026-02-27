## 1. Settings Help Update

- [x] 1.1 Replace Long description in `internal/cli/settings/settings.go` with group-based category listing and `/` search mention

## 2. Doctor Help Update

- [x] 2.1 Replace Long description in `internal/cli/doctor/doctor.go` with all 14 checks and `--fix`/`--json` flag guidance

## 3. Onboard Help Update

- [x] 3.1 Replace Long description in `internal/cli/onboard/onboard.go` with GitHub provider, auto-fetch models, and approval policy

## 4. Verification

- [x] 4.1 Run `go build ./...` and `go test ./...` to verify no regressions
- [x] 4.2 Verify `lango settings --help`, `lango doctor --help`, and `lango onboard --help` output matches specs
