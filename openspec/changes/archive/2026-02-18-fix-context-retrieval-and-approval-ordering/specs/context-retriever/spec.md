## MODIFIED Requirements

### Requirement: Keyword extraction from query
The `extractKeywords` function SHALL extract at most 5 keywords from a query string. Each keyword SHALL be at most 50 characters long. Keywords SHALL contain only alphanumeric characters, hyphens, and underscores after sanitization. Keywords shorter than 2 characters after sanitization SHALL be excluded.

#### Scenario: Normal query with multiple keywords
- **WHEN** a query "how to handle errors in Go deployment configuration" is processed
- **THEN** at most 5 non-stop-word keywords are returned, each sanitized to alphanumeric/hyphen/underscore characters

#### Scenario: Query with only stop words
- **WHEN** a query "the a is are was" is processed
- **THEN** nil is returned (no keywords extracted)

#### Scenario: Query with special characters
- **WHEN** a query contains characters like `@`, `#`, `$`, `(`, `)`, etc.
- **THEN** those characters are stripped from keywords during sanitization

#### Scenario: Very long keyword
- **WHEN** a single word exceeds 50 characters
- **THEN** it is truncated to 50 characters
