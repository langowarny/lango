## ADDED Requirements

### Requirement: Browser binary auto-detection
The system SHALL auto-detect system-installed browser binaries using `launcher.LookPath()` when no explicit `BrowserBin` config is set.

#### Scenario: System browser found
- **WHEN** `BrowserBin` config is empty
- **AND** `launcher.LookPath()` finds a system browser
- **THEN** the system SHALL use the detected binary path

#### Scenario: Explicit browser path
- **WHEN** `BrowserBin` config is set to a non-empty path
- **THEN** the system SHALL use the configured path regardless of LookPath result

#### Scenario: No browser found
- **WHEN** `BrowserBin` config is empty
- **AND** `launcher.LookPath()` does not find a system browser
- **THEN** the system SHALL fall back to go-rod's default browser download behavior

### Requirement: BrowserBin config field
The `BrowserToolConfig` SHALL include a `BrowserBin` string field for specifying an explicit browser binary path.

#### Scenario: Config field recognition
- **WHEN** `tools.browser.browserBin: "/usr/bin/chromium"` is set in configuration
- **THEN** the system SHALL pass the path to `launcher.Bin()`
