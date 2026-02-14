## 1. Remove Provider OAuth Code

- [x] 1.1 Delete `internal/cli/auth/` directory (OAuth login CLI command)
- [x] 1.2 Remove `auth` import and `auth.NewCommand()` from `cmd/lango/main.go`
- [x] 1.3 Remove OAuth imports, `getAccessToken()`, and `saveToken()` from `internal/supervisor/supervisor.go`
- [x] 1.4 Replace OAuth fallback with warning log when API key is empty in supervisor
- [x] 1.5 Remove `ClientID`, `ClientSecret`, `Scopes` fields from `ProviderConfig` in `internal/config/types.go`
- [x] 1.6 Remove OAuth field expansion from `substituteEnvVars()` in `internal/config/loader.go`
- [x] 1.7 Remove `clientSecret` registration from `registerConfigSecrets()` in `internal/app/app.go`

## 2. API Key Security Check

- [x] 2.1 Create `internal/cli/doctor/checks/apikey_security.go` with `APIKeySecurityCheck`
- [x] 2.2 Register `APIKeySecurityCheck` in `AllChecks()` in `internal/cli/doctor/checks/checks.go`
- [x] 2.3 Create `internal/cli/doctor/checks/apikey_security_test.go` with table-driven tests

## 3. OpenSpec Documentation

- [x] 3.1 Update `openspec/specs/oauth-login/spec.md` with REMOVED status
- [x] 3.2 Update `openspec/specs/agent-provider-config/spec.md` to remove OAuth scenarios

## 4. Verification

- [x] 4.1 Verify `go build ./...` succeeds
- [x] 4.2 Verify `go vet ./...` passes
- [x] 4.3 Verify `go test ./internal/supervisor/...` passes
- [x] 4.4 Verify `go test ./internal/config/...` passes
- [x] 4.5 Verify `go test ./internal/cli/doctor/...` passes
- [x] 4.6 Run `go mod tidy` to clean dependencies
