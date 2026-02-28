## MODIFIED Requirements

### Requirement: Identity command output
The `lango p2p identity` command SHALL display `keyStorage` information (either "secrets-store" or "file") instead of the raw `keyDir` filesystem path.

#### Scenario: Identity with encrypted storage
- **WHEN** the user runs `lango p2p identity` and SecretsStore is available
- **THEN** the output SHALL show `Key Storage: secrets-store` instead of a directory path

#### Scenario: Identity with file storage
- **WHEN** the user runs `lango p2p identity` and SecretsStore is not available
- **THEN** the output SHALL show `Key Storage: file`

#### Scenario: JSON output reflects key storage
- **WHEN** the user runs `lango p2p identity --json`
- **THEN** the JSON SHALL contain `"keyStorage": "secrets-store"` or `"keyStorage": "file"` instead of `"keyDir"`
