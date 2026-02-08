## ADDED Requirements

### Requirement: File reading
The system SHALL read file contents with encoding detection and size limits.

#### Scenario: Read text file
- **WHEN** reading a text file
- **THEN** the content SHALL be returned as a string with detected encoding

#### Scenario: Read binary file
- **WHEN** reading a binary file
- **THEN** the content SHALL be returned as base64 or indicate binary nature

#### Scenario: File size limit
- **WHEN** a file exceeds the configured size limit
- **THEN** an error SHALL be returned with the actual size

### Requirement: File writing
The system SHALL write content to files with atomic write support.

#### Scenario: Write new file
- **WHEN** writing to a non-existent path
- **THEN** the file SHALL be created with specified content

#### Scenario: Overwrite existing file
- **WHEN** writing to an existing file
- **THEN** the content SHALL replace the previous content

#### Scenario: Create parent directories
- **WHEN** parent directories do not exist
- **THEN** they SHALL be created automatically

### Requirement: File editing
The system SHALL support surgical edits to existing files.

#### Scenario: Line range replacement
- **WHEN** editing with a line range and replacement content
- **THEN** only the specified lines SHALL be replaced

#### Scenario: Search and replace
- **WHEN** editing with a search pattern and replacement
- **THEN** matching content SHALL be replaced

### Requirement: Directory operations
The system SHALL support listing and navigating directories.

#### Scenario: List directory contents
- **WHEN** listing a directory
- **THEN** files and subdirectories SHALL be returned with metadata

#### Scenario: Recursive listing
- **WHEN** listing with recursive option
- **THEN** all nested contents SHALL be included up to depth limit

### Requirement: Path safety
The system SHALL validate file paths to prevent directory traversal attacks.

#### Scenario: Path traversal attempt
- **WHEN** a path contains ".." to escape allowed directory
- **THEN** the operation SHALL be rejected with an error
