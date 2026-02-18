## MODIFIED Requirements

### Requirement: Post-save next steps output
After saving a profile, the onboard command SHALL display a confirmation message with the profile name and storage location, followed by concise next steps: starting Lango (`lango serve`), running doctor (`lango doctor`), and profile management commands (`lango config list`, `lango config use`). The output SHALL NOT include environment variable export instructions or channel-specific token setup guidance.

#### Scenario: Successful save displays simplified next steps
- **WHEN** user completes the onboard wizard and saves
- **THEN** system prints confirmation with profile name and storage path, followed by numbered next steps for serve and doctor, and profile management hints
- **THEN** system does NOT print any `export` commands or environment variable guidance

### Requirement: No env example file generation
The onboard command SHALL NOT generate a `.lango.env.example` file after saving configuration. The `generateEnvExample()` function SHALL be removed.

#### Scenario: Save does not create .lango.env.example
- **WHEN** user completes the onboard wizard and saves
- **THEN** no `.lango.env.example` file is created in the working directory

### Requirement: Command Long description
The onboard command's Long description SHALL list all current configuration categories: Agent, Server, Channels, Tools, Auth, Security, Session, Knowledge, and Providers. The description SHALL note that all settings are saved in an encrypted profile.

#### Scenario: Long description includes Auth and Session categories
- **WHEN** user runs `lango onboard --help`
- **THEN** the Long description includes Auth (OIDC providers, JWT settings) and Session (Session DB, TTL) as separate categories
- **THEN** Security category lists PII interceptor and Signer only

## REMOVED Requirements

### Requirement: Environment variable export guidance
**Reason**: Encrypted SQLite profiles make direct API key storage safe; environment variable guidance is misleading and creates confusion.
**Migration**: API keys are stored directly in encrypted profiles. Users who prefer environment variables can still use `${ENV_VAR}` references.

### Requirement: .lango.env.example file generation
**Reason**: Exposes key names in plaintext, incompatible with encrypted storage architecture.
**Migration**: No replacement needed; keys are managed through encrypted profiles.
