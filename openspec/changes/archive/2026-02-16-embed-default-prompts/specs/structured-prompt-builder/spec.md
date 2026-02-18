## MODIFIED Requirements

### Requirement: Default sections included
The system SHALL provide a `DefaultBuilder()` that returns a Builder with four built-in sections: Identity (priority 100), Safety (priority 200), Conversation Rules (priority 300), and Tool Usage (priority 400). Section content SHALL be read from embedded `.md` files in the `prompts` package via `prompts.FS.ReadFile()`. If an embedded file read fails, the system SHALL use a minimal fallback string for that section.

#### Scenario: Default builder includes conversation rules
- **WHEN** DefaultBuilder().Build() is called
- **THEN** the output SHALL contain conversation rules instructing the LLM to focus on the current question and not repeat previous answers

#### Scenario: Default builder section order
- **WHEN** DefaultBuilder().Build() is called
- **THEN** Identity SHALL appear before Safety, Safety before Conversation Rules, Conversation Rules before Tool Usage

#### Scenario: Default builder uses embedded content
- **WHEN** DefaultBuilder().Build() is called with a correctly built binary
- **THEN** the Identity section SHALL contain content from AGENTS.md describing Lango's five tool categories
- **AND** the Safety section SHALL contain content from SAFETY.md with security rules
- **AND** the Conversation Rules section SHALL contain content from CONVERSATION_RULES.md
- **AND** the Tool Usage section SHALL contain content from TOOL_USAGE.md with per-tool guidelines
