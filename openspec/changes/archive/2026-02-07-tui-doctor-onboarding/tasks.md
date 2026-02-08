## 1. Project Setup

- [x] 1.1 Add bubbletea and lipgloss dependencies to go.mod
- [x] 1.2 Create internal/cli/tui/ directory structure for shared TUI components
- [x] 1.3 Implement shared TUI styles in internal/cli/tui/styles.go

## 2. Doctor Command - Core

- [x] 2.1 Create internal/cli/doctor/ package structure
- [x] 2.2 Define Check interface and Result types in checks/checks.go
- [x] 2.3 Implement configuration file check (checks/config.go)
- [x] 2.4 Implement API key verification check (checks/apikey.go)
- [x] 2.5 Implement channel token validation check (checks/channels.go)
- [x] 2.6 Implement session database check (checks/database.go)
- [x] 2.7 Implement server port availability check (checks/network.go)

## 3. Doctor Command - Output

- [x] 3.1 Implement TUI output renderer (output/tui.go)
- [x] 3.2 Implement JSON output renderer (output/json.go)
- [x] 3.3 Implement check result summary display

## 4. Doctor Command - Fix Mode

- [x] 4.1 Add --fix flag support to doctor command
- [x] 4.2 Implement auto-repair for missing database directory
- [x] 4.3 Implement auto-repair for missing config file

## 5. Doctor Command - Integration

- [x] 5.1 Add doctor subcommand to cmd/lango/main.go
- [x] 5.2 Add --json flag for non-interactive output
- [x] 5.3 Write unit tests for doctor checks

## 6. Onboard Command - Core

- [x] 6.1 Create internal/cli/onboard/ package structure
- [x] 6.2 Implement welcome screen with mode selection
- [x] 6.3 Implement API key input step with validation
- [x] 6.4 Implement model selection step with dropdown
- [x] 6.5 Implement channel selection and setup step

## 7. Onboard Command - Configuration

- [x] 7.1 Implement configuration file generation
- [x] 7.2 Implement existing config detection and handling
- [x] 7.3 Display environment variable hints after save

## 8. Onboard Command - Integration

- [x] 8.1 Add onboard subcommand to cmd/lango/main.go
- [x] 8.2 Implement post-setup doctor verification prompt
- [x] 8.3 Write unit tests for onboard flow

## 9. Testing & Documentation

- [x] 9.1 Add integration tests for doctor command
- [x] 9.2 Add integration tests for onboard command
- [x] 9.3 Update README.md with new commands
