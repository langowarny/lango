## Purpose

External skill import system — fetches SKILL.md files from GitHub repositories and arbitrary URLs, enabling the agent to expand its capabilities with community-published instruction skills.

## Requirements

### Requirement: GitHub URL parsing
The system SHALL parse GitHub repository URLs into owner, repo, branch, and path components. Supported formats: `https://github.com/owner/repo`, `https://github.com/owner/repo/tree/branch`, and `https://github.com/owner/repo/tree/branch/path`.

#### Scenario: Basic repo URL
- **WHEN** the URL is `https://github.com/kepano/obsidian-skills`
- **THEN** the system SHALL parse owner=`kepano`, repo=`obsidian-skills`, branch=`main`, path=``

#### Scenario: URL with branch
- **WHEN** the URL is `https://github.com/kepano/obsidian-skills/tree/develop`
- **THEN** the system SHALL parse branch=`develop`

#### Scenario: URL with branch and path
- **WHEN** the URL is `https://github.com/kepano/obsidian-skills/tree/main/skills`
- **THEN** the system SHALL parse branch=`main`, path=`skills`

#### Scenario: Invalid URL
- **WHEN** the URL does not contain at least owner and repo segments
- **THEN** the system SHALL return an error

### Requirement: GitHub URL detection
The system SHALL determine whether a URL points to github.com.

#### Scenario: GitHub URL
- **WHEN** the URL contains `github.com`
- **THEN** `IsGitHubURL` SHALL return true

#### Scenario: Non-GitHub URL
- **WHEN** the URL does not contain `github.com`
- **THEN** `IsGitHubURL` SHALL return false

### Requirement: Skill discovery from GitHub repo
The system SHALL list subdirectories at the given path in a GitHub repository via the Contents API. Each subdirectory is assumed to contain a SKILL.md file.

#### Scenario: Directory with skills
- **WHEN** the GitHub Contents API returns a mix of directories and files
- **THEN** the system SHALL return only directory names as skill candidates

### Requirement: Single skill fetch from GitHub
The system SHALL fetch a SKILL.md file from `{path}/{skillName}/SKILL.md` via the GitHub Contents API and decode the base64-encoded content.

#### Scenario: Successful fetch
- **WHEN** the Contents API returns a base64-encoded file response
- **THEN** the system SHALL decode the content and return raw bytes

#### Scenario: Fetch failure
- **WHEN** the Contents API returns a non-2xx status
- **THEN** the system SHALL return an error

### Requirement: URL-based skill fetch
The system SHALL fetch raw SKILL.md content from any URL via HTTP GET.

#### Scenario: Successful fetch
- **WHEN** the URL returns 200 OK with markdown content
- **THEN** the system SHALL return the raw bytes

### Requirement: Bulk import from repository
The system SHALL discover all skills in a repo, fetch and parse each SKILL.md, and save them to the store. The `ImportFromRepo` method SHALL prefer git clone when git is available, falling back to the GitHub HTTP API when git is not installed or when clone fails. Skills that already exist SHALL be skipped. Failed skills SHALL be recorded as errors without stopping the import.

#### Scenario: Full import
- **WHEN** a repo contains skills `A`, `B`, and `C` where `A` already exists
- **THEN** the system SHALL import `B` and `C`, skip `A`, and report the result

#### Scenario: Partial failure
- **WHEN** skill `B` fails to parse
- **THEN** the system SHALL continue importing remaining skills and record `B` as an error

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

### Requirement: Single skill import
The system SHALL parse raw SKILL.md content, set the source URL, and save the skill entry to the store.

#### Scenario: Import from raw content
- **WHEN** valid SKILL.md bytes and a source URL are provided
- **THEN** the system SHALL parse, set source, save, and return the entry

### Requirement: import_skill agent tool
The system SHALL expose an `import_skill` tool with `url` (required) and `skill_name` (optional) parameters. It SHALL support three modes: GitHub bulk import, GitHub single import, and direct URL import.

#### Scenario: GitHub bulk import
- **WHEN** `url` is a GitHub URL and `skill_name` is empty
- **THEN** the system SHALL discover and import all skills from the repository

#### Scenario: GitHub single import
- **WHEN** `url` is a GitHub URL and `skill_name` is specified
- **THEN** the system SHALL import only the named skill

#### Scenario: Direct URL import
- **WHEN** `url` is not a GitHub URL
- **THEN** the system SHALL fetch the URL directly and parse it as a SKILL.md file

#### Scenario: Tool reload after import
- **WHEN** import completes successfully
- **THEN** the system SHALL reload skills so imported skills appear as tools immediately

### Requirement: Single skill import with resources
The `ImportSingleWithResources` method SHALL import a single named skill from a GitHub repository, including its resource directories.

#### Scenario: Single skill with resources via git
- **WHEN** `ImportSingleWithResources` is called with git available
- **THEN** the system SHALL clone the repo, read the specific skill's SKILL.md, and copy resource directories

#### Scenario: Single skill with resources via HTTP
- **WHEN** `ImportSingleWithResources` is called without git
- **THEN** the system SHALL fetch SKILL.md via API and fetch resource files from recognized directories

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

### Requirement: AllowImport config flag
The `SkillConfig` SHALL include an `AllowImport` boolean field that defaults to `true`.

#### Scenario: Default value
- **WHEN** no `allowImport` is specified in config
- **THEN** it SHALL default to `true`
