## MODIFIED Requirements

### Requirement: SAFETY prompt reflects PII detection scope
The SAFETY.md prompt SHALL enumerate specific PII categories (email, phone numbers, national IDs, financial account numbers) and mention 13 builtin patterns. The prompt SHALL reference Presidio NER-based detection as additional coverage when enabled.

#### Scenario: Agent responds to PII-related user question
- **WHEN** the agent processes SAFETY.md prompt during system prompt assembly
- **THEN** the agent understands it protects 13 builtin PII pattern categories
- **THEN** the agent can accurately inform users about PII protection coverage including Presidio
