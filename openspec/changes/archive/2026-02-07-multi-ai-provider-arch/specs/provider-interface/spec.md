## ADDED Requirements

### Requirement: Provider Interface Definition
The system SHALL define a `Provider` interface that all LLM provider implementations MUST implement.

#### Scenario: Interface methods defined
- **WHEN** a new LLM provider is implemented
- **THEN** it SHALL implement `ID() string` returning the provider identifier
- **AND** SHALL implement `Generate(ctx, params) iter.Seq2[StreamEvent, error]` for streaming responses
- **AND** SHALL implement `ListModels(ctx) ([]ModelInfo, error)` for model discovery

### Requirement: Streaming Response Support
The system SHALL support streaming LLM responses via Go iterators.

#### Scenario: Text streaming
- **WHEN** the provider generates a response
- **THEN** it SHALL yield `StreamEvent` with `Type: "text_delta"` for each text chunk

#### Scenario: Tool call streaming
- **WHEN** the provider generates a tool call
- **THEN** it SHALL yield `StreamEvent` with `Type: "tool_call"` containing the tool call details

#### Scenario: Stream completion
- **WHEN** the response generation completes
- **THEN** it SHALL yield `StreamEvent` with `Type: "done"`

#### Scenario: Error during streaming
- **WHEN** an error occurs during generation
- **THEN** it SHALL yield `StreamEvent` with `Type: "error"` and the error details

### Requirement: Generation Parameters
The system SHALL accept standard generation parameters across all providers.

#### Scenario: Common parameters supported
- **WHEN** `GenerateParams` is passed to a provider
- **THEN** it SHALL accept `Model`, `Messages`, `Tools`, `Temperature`, and `MaxTokens`

### Requirement: Model Information
The system SHALL provide standardized model metadata.

#### Scenario: ModelInfo structure
- **WHEN** `ListModels` is called
- **THEN** each `ModelInfo` SHALL contain `ID`, `Name`, `ContextWindow`, `SupportsVision`, `SupportsTools`, and `IsReasoning` fields
