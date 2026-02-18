## ADDED Requirements

### Requirement: Directory and file structure management
The system SHALL allow the agent to explore and manage the directory structure.

#### Scenario: List directory with metadata
- **WHEN** a directory listing is requested
- **THEN** the system SHALL return a list of files and subdirectories with sizes and modification times

#### Scenario: Create directory
- **WHEN** a directory creation is requested
- **THEN** the system SHALL create the directory and any necessary parents

#### Scenario: Delete file or directory
- **WHEN** deletion is requested for a path
- **THEN** the system SHALL remove the target and its contents if it is a directory
