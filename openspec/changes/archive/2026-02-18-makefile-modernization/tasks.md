## 1. Variables & Foundation

- [x] 1.1 Update GOBUILD and GOTEST to include CGO_ENABLED=1
- [x] 1.2 Add REGISTRY and COVERAGE_DIR variables
- [x] 1.3 Update .PHONY to list all targets

## 2. Build & Install Targets

- [x] 2.1 Add `install` target to install binary to $GOPATH/bin
- [x] 2.2 Add `run` target to start server from existing binary

## 3. Testing & Coverage

- [x] 3.1 Add `coverage` target generating HTML report to .coverage/

## 4. Code Quality Targets

- [x] 4.1 Add `fmt` target for gofmt formatting
- [x] 4.2 Add `fmt-check` target for CI format violation checks
- [x] 4.3 Add `vet` target for go vet
- [x] 4.4 Update `lint` to auto-install golangci-lint if missing
- [x] 4.5 Add `ci` target chaining fmt-check → vet → lint → test

## 5. Docker Build Improvements

- [x] 5.1 Update `docker-build` to tag with both version and latest
- [x] 5.2 Add `docker-build-browser` target with WITH_BROWSER=true
- [x] 5.3 Add `docker-push` target with REGISTRY validation

## 6. Docker Compose Targets

- [x] 6.1 Add `docker-up` for default profile
- [x] 6.2 Add `docker-up-browser` for browser profile
- [x] 6.3 Add `docker-up-sidecar` for browser-sidecar profile
- [x] 6.4 Add `docker-down` to stop all profile containers
- [x] 6.5 Add `docker-logs` to tail container logs

## 7. Utility & Cleanup

- [x] 7.1 Add `health` target using curl to check server
- [x] 7.2 Update `clean` to also remove .coverage/ directory
- [x] 7.3 Implement comment-based `help` auto-generation
- [x] 7.4 Add bin/ and .coverage/ to .gitignore
