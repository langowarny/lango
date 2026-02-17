## MODIFIED Requirements

### Requirement: Model adapter streams LLM responses
The `ModelAdapter.GenerateContent` method SHALL respect the `stream` parameter to control how `LLMResponse` events are yielded to the ADK runner.

When `stream` is `false`, the adapter SHALL consume all provider `StreamEvent`s internally, accumulate text and tool call parts, and yield exactly one `LLMResponse` with `Partial=false`, `TurnComplete=true`, and `Content.Parts` containing the full accumulated text and tool calls.

When `stream` is `true`, the adapter SHALL yield partial `LLMResponse` events for each text delta (`Partial=true`), and the final done event SHALL include the fully accumulated text in `Content.Parts` with `Partial=false` and `TurnComplete=true`.

#### Scenario: Non-streaming mode accumulates complete response
- **WHEN** `GenerateContent` is called with `stream=false` and the provider emits text deltas "Hello " and "world" followed by a done event
- **THEN** the adapter yields exactly one `LLMResponse` with `Partial=false`, `TurnComplete=true`, and `Content.Parts[0].Text` equal to "Hello world"

#### Scenario: Non-streaming mode accumulates tool calls
- **WHEN** `GenerateContent` is called with `stream=false` and the provider emits a tool call event followed by a done event
- **THEN** the adapter yields exactly one `LLMResponse` with `Partial=false`, `TurnComplete=true`, and `Content.Parts` containing the `FunctionCall`

#### Scenario: Streaming mode yields partial events and complete final
- **WHEN** `GenerateContent` is called with `stream=true` and the provider emits text deltas "Hello " and "world" followed by a done event
- **THEN** the adapter yields two partial `LLMResponse` events (one per delta) with `Partial=true`, followed by one final `LLMResponse` with `Partial=false`, `TurnComplete=true`, and `Content.Parts[0].Text` equal to "Hello world"

#### Scenario: Provider error propagates correctly
- **WHEN** the provider emits a `StreamEventError` event in either streaming or non-streaming mode
- **THEN** the adapter yields the error immediately and stops iteration

### Requirement: Agent text collection avoids duplication
The `runAndCollectOnce` method SHALL collect text from either partial events or the final non-partial event, but never both, to prevent duplicate text in the response.

#### Scenario: Streaming mode collects from partial events only
- **WHEN** `runAndCollectOnce` processes events that include partial text events followed by a non-partial final event containing the same accumulated text
- **THEN** the method returns text collected only from partial events, not from the final event

#### Scenario: Non-streaming mode collects from final event
- **WHEN** `runAndCollectOnce` processes events that contain no partial events and one non-partial final event with text
- **THEN** the method returns text from the non-partial final event
