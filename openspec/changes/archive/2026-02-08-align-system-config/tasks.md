# Align System Config Tasks

## CLI Doctor
-   [x] Update `internal/cli/doctor/checks/security.go` to handle encrypted databases gracefully
-   [x] Implement specific check for "out of memory" / "encrypted" errors when opening session DB
-   [x] Verify `doctor` command output with encrypted database (manual test)

## Config & Docs
-   [x] Update `lango.example.json` with `providers` map
-   [x] Update `lango.example.json` with expanded `security` section
-   [x] Update `README.md` to document `providers` configuration
-   [x] Update `README.md` to document TUI security capabilities
