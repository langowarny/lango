## Context

The ADK runner calls `ModelAdapter.GenerateContent(ctx, req, stream)` with `stream=false` by default. The current implementation ignores this flag and always streams provider events through, yielding partial text events (`Partial=true`) followed by an empty done event (`Partial=false, Parts=[]`). The ADK runner only commits non-partial events to the session, resulting in empty assistant messages being stored. On the next turn, the LLM sees unanswered questions in history and re-answers them.

## Goals / Non-Goals

**Goals:**
- Ensure non-streaming mode yields a single complete `LLMResponse` with accumulated text
- Ensure streaming mode's final done event includes accumulated full text for session storage
- Prevent double-counting of text in `runAndCollectOnce` when both partial and non-partial events carry text

**Non-Goals:**
- Changing the provider interface or adding a non-streaming provider path
- Modifying the ADK runner's session commit logic
- Adding new streaming capabilities or SSE endpoints

## Decisions

1. **Accumulate internally in non-streaming mode**: When `stream=false`, consume all provider `StreamEvent`s internally using `strings.Builder`, then yield one `LLMResponse{Partial: false, TurnComplete: true}` with all accumulated text and tool calls. This matches the ADK runner's expectation for non-streaming responses.

2. **Enrich the done event in streaming mode**: When `stream=true`, continue yielding partial text events for real-time UI display. In the final done event, include the fully accumulated text so the ADK runner stores a complete assistant message in the session.

3. **Guard against double-counting in `runAndCollectOnce`**: Track whether partial events were seen. If partial events were collected, skip the non-partial final event's text (it's a duplicate). If no partial events were seen (non-streaming mode), collect from the final event.

## Risks / Trade-offs

- **[Risk]** Streaming mode's done event now carries full text, increasing its payload size → This is negligible compared to total stream size, and necessary for correct session storage.
- **[Risk]** If a provider emits no done event, non-streaming mode would hang → Mitigated by the fact that all existing providers (OpenAI, Anthropic, Gemini) always emit a done event to signal stream completion.
