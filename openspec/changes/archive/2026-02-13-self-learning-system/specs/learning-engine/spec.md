## ADDED Requirements

### Requirement: Tool Result Observation
The system SHALL observe every tool execution result to detect error patterns and track successes.

#### Scenario: Tool execution error
- **WHEN** `OnToolResult` is called with a non-nil error
- **THEN** the system SHALL extract the error pattern, categorize it, and store a learning entry

#### Scenario: Tool execution success
- **WHEN** `OnToolResult` is called with a nil error
- **THEN** the system SHALL search for related learnings by tool name and boost their confidence

#### Scenario: Skip duplicate high-confidence learnings
- **WHEN** an error occurs and a matching learning with confidence > 0.5 already exists
- **THEN** the system SHALL skip creating a new learning entry

### Requirement: Error Pattern Extraction
The system SHALL normalize error messages into reusable patterns.

#### Scenario: Strip dynamic identifiers
- **WHEN** extracting an error pattern
- **THEN** the system SHALL replace UUIDs, timestamps, file paths, and port numbers with placeholders

### Requirement: Error Categorization
The system SHALL categorize errors into predefined categories.

#### Scenario: Timeout errors
- **WHEN** the error is a context deadline exceeded or contains "timeout"
- **THEN** the system SHALL categorize it as "timeout"
- **AND** SHALL use `errors.Is` for proper wrapped error detection of `context.DeadlineExceeded`

#### Scenario: Permission errors
- **WHEN** the error contains "permission denied", "access denied", or "forbidden"
- **THEN** the system SHALL categorize it as "permission"

#### Scenario: Provider errors
- **WHEN** the error contains "api", "model", "provider", or "rate limit"
- **THEN** the system SHALL categorize it as "provider_error"

#### Scenario: Tool errors
- **WHEN** a tool name is provided and the error does not match other categories
- **THEN** the system SHALL categorize it as "tool_error"

### Requirement: Parameter Summarization
The system SHALL summarize tool parameters before storing them in learnings.

#### Scenario: Truncate long strings
- **WHEN** a parameter value is a string longer than 200 characters
- **THEN** the system SHALL truncate it to 200 characters with "..." suffix

#### Scenario: Summarize arrays
- **WHEN** a parameter value is an array
- **THEN** the system SHALL replace it with a count summary (e.g., "[5 items]")
