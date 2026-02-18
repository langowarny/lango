## Requirement: System prompt construction
The system SHALL construct the system prompt using a structured `prompt.Builder` instead of a single string. The `ContextAwareModelAdapter` constructor SHALL accept a `*prompt.Builder` and call `Build()` to produce the base prompt string. Dynamic context injection (knowledge, memory, RAG) SHALL continue to append to the built prompt at runtime.

#### Scenario: Prepend system prompt to new session
- **WHEN** a new agent session is started
- **THEN** the system SHALL prepend a message with `role: system` containing the configured prompt to the conversation history

#### Scenario: Default identity prompt
- **WHEN** no custom system prompt is provided
- **THEN** a default prompt describing the agent's identity and tools SHALL be used

#### Scenario: Builder produces base prompt
- **WHEN** ContextAwareModelAdapter is created with a prompt.Builder
- **THEN** the basePrompt field SHALL equal the builder's Build() output

#### Scenario: Dynamic context still appended
- **WHEN** knowledge retrieval returns context during GenerateContent
- **THEN** the retrieved context SHALL be appended to the builder-produced base prompt
