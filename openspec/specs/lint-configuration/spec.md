## Purpose

Defines linting standards and configuration for the Lango project using golangci-lint v2.

## Requirements

### Requirement: golangci-lint v2 configuration
The project SHALL have a `.golangci.yml` configuration file using version 2 format with the `standard` default linter set.

#### Scenario: Generated code exclusion
- **WHEN** golangci-lint runs on the project
- **THEN** files with `// Code generated` headers (ent auto-generated code) SHALL be excluded via `generated: strict`

#### Scenario: Standard error handling preset
- **WHEN** golangci-lint evaluates error handling patterns
- **THEN** standard patterns (defer Close, fmt.Fprint return values) SHALL be suppressed via `std-error-handling` preset

#### Scenario: Test file errcheck exclusion
- **WHEN** golangci-lint evaluates test files (`_test.go`)
- **THEN** errcheck linter SHALL be disabled for those files

### Requirement: Zero lint issues in CI
The project SHALL pass golangci-lint with zero issues on every CI run.

#### Scenario: Clean lint run
- **WHEN** `golangci-lint run` executes on the codebase
- **THEN** the exit code SHALL be 0 with zero reported issues

### Requirement: Explicit error handling for intentionally ignored errors
All intentionally ignored error return values SHALL use explicit `_ =` assignment to document intent.

#### Scenario: Defer close pattern
- **WHEN** an HTTP response body is closed in a defer
- **THEN** the pattern `defer func() { _ = resp.Body.Close() }()` SHALL be used

#### Scenario: Rollback in error paths
- **WHEN** a database transaction rollback is called in an error/defer path
- **THEN** the pattern `_ = tx.Rollback()` SHALL be used
