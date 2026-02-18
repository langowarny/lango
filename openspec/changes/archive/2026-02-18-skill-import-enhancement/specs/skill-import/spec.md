## MODIFIED Requirements

### Requirement: Import skills from GitHub repository
The `ImportFromRepo` method SHALL prefer git clone when git is available, falling back to the GitHub HTTP API when git is not installed or when clone fails.

#### Scenario: Git available — uses git clone
- **WHEN** `git` is available on PATH and `ImportFromRepo` is called
- **THEN** the system SHALL clone the repository with `git clone --depth=1 --branch <branch>` to a temp directory
- **AND** discover skill directories by scanning for subdirectories containing SKILL.md
- **AND** import each skill with its resource directories
- **AND** clean up the temp directory after import

#### Scenario: Git not available — uses HTTP API
- **WHEN** `git` is not available on PATH
- **THEN** the system SHALL use the GitHub Contents API to discover and fetch skills
- **AND** fetch resource files from recognized resource directories via API

#### Scenario: Git clone fails — falls back to HTTP API
- **WHEN** `git` is available but `git clone` fails (network error, invalid repo)
- **THEN** the system SHALL log a warning and fall back to the GitHub HTTP API method

### Requirement: Single skill import with resources
The `ImportSingleWithResources` method SHALL import a single named skill from a GitHub repository, including its resource directories.

#### Scenario: Single skill with resources via git
- **WHEN** `ImportSingleWithResources` is called with git available
- **THEN** the system SHALL clone the repo, read the specific skill's SKILL.md, and copy resource directories

#### Scenario: Single skill with resources via HTTP
- **WHEN** `ImportSingleWithResources` is called without git
- **THEN** the system SHALL fetch SKILL.md via API and fetch resource files from recognized directories

## ADDED Requirements

### Requirement: Exec guard redirects skill-related commands
`blockLangoExec` SHALL redirect skill-related `git clone`, `curl`, and `wget` commands to the `import_skill` tool with a guidance message.

#### Scenario: Git clone with skill keyword
- **WHEN** a command starts with "git clone" and contains "skill"
- **THEN** the system SHALL return a message redirecting to `import_skill`

#### Scenario: Curl with skill keyword
- **WHEN** a command starts with "curl " and contains "skill"
- **THEN** the system SHALL return a message redirecting to `import_skill`

#### Scenario: Git clone without skill keyword
- **WHEN** a command starts with "git clone" but does NOT contain "skill"
- **THEN** the system SHALL NOT block the command
