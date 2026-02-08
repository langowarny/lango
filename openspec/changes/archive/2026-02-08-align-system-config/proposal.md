# Align System Configuration

## Why
The recent TUI upgrades for `security` and `providers` interact with the system configuration in ways that breach existing documentation and tooling.
1.  **`doctor` Failure**: The `doctor` command crashes or reports errors because it attempts to open the encrypted session database without the necessary passphrase context, due to the new `security` configuration.
2.  **Documentation Gap**: `README.md` and `lango.example.json` do not reflect the new `providers` map or the expanded `security` options, leaving users without guidance on how to configure these features manually.

Aligning these components is critical to ensure a consistent and broken-free user experience.

## What Changes
1.  **Fix `doctor`**: Update `internal/cli/doctor/checks/security.go` to correctly handle the new `security` configuration, specifically ensuring the DB check respects the encryption settings.
2.  **Update Config Example**: Add `providers` and expanded `security` sections to `lango.example.json`.
3.  **Update Documentation**: Update `README.md` to document the new configuration structure and TUI capabilities.

## Capabilities

### New Capabilities
- `docs`: Comprehensive documentation for new Security and Provider features.

### Modified Capabilities
- `cli-doctor`: Enhanced security checks that are compatible with the new encryption scheme.

## Impact
-   **CLI**: `doctor` command will function correctly with encrypted databases.
-   **Docs**: `README.md` and `lango.example.json` will be accurate.
-   **User Experience**: Users will have working diagnostics and correct examples.
