## MODIFIED Requirements

### Requirement: Set secret value
The `secrets set` command SHALL support storing a secret value either interactively (via passphrase prompt) or non-interactively (via `--value-hex` flag). When `--value-hex` is provided, the command SHALL hex-decode the input (stripping an optional `0x` prefix) and store the raw bytes. When `--value-hex` is not provided, the command SHALL require an interactive terminal and prompt for the value.

#### Scenario: Interactive secret storage
- **WHEN** user runs `lango security secrets set mykey` in an interactive terminal without `--value-hex`
- **THEN** the command SHALL prompt for the secret value and store it encrypted

#### Scenario: Non-interactive hex secret storage
- **WHEN** user runs `lango security secrets set wallet.privatekey --value-hex 0xac0974...` in a non-interactive environment
- **THEN** the command SHALL hex-decode the value (stripping `0x` prefix), store the raw bytes encrypted, and print success

#### Scenario: Non-interactive without value-hex flag
- **WHEN** user runs `lango security secrets set mykey` in a non-interactive terminal without `--value-hex`
- **THEN** the command SHALL return an error suggesting `--value-hex` for non-interactive use

#### Scenario: Invalid hex value
- **WHEN** user runs `lango security secrets set mykey --value-hex "not-hex"`
- **THEN** the command SHALL return a hex decode error
