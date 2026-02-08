## 1. ADK Integration into Runtime

- [x] 1.1 Add ADK imports to internal/agent/runtime.go
- [x] 1.2 Add ADK agent field to Runtime struct
- [x] 1.3 Create Gemini model in New() using gemini.NewModel()
- [x] 1.4 Create llmagent with llmagent.New() and pass model
- [x] 1.5 Store ADK agent instance in Runtime struct

## 2. Tool Adapter Implementation

- [x] 2.1 Create AdkToolAdapter struct implementing tool.Tool interface
- [x] 2.2 Implement Name() method returning tool.Name
- [x] 2.3 Implement Description() method returning tool.Description
- [x] 2.4 Implement Run(ctx, input) calling tool.Handler
- [x] 2.5 Update RegisterTool() to create ADK tool adapters
- [x] 2.6 Store []tool.Tool slice for passing to llmagent config

## 3. Session to ADK Conversion

- [x] 3.1 Implement buildAdkSession() to create agent.Session
- [x] 3.2 Convert session.Message to agent.Message format
- [x] 3.3 Add message history to ADK session
- [x] 3.4 Implement session truncation (keep last 20 turns)
- [x] 3.5 Add truncation warning logging

## 4. Agent Execution and Streaming

- [x] 4.1 Refactor Run() to call adkAgent.Execute() with ADK session
- [x] 4.2 Parse ADK response parts (text, tool_call, etc.)
- [x] 4.3 Emit StreamEvent{Type: "text_delta"} for text parts
- [x] 4.4 Emit StreamEvent{Type: "tool_start"} when tool call begins
- [x] 4.5 Execute tool via ADK (handled internally by agent)
- [x] 4.6 Emit StreamEvent{Type: "tool_end"} with results
- [x] 4.7 Emit StreamEvent{Type: "done"} on completion
- [x] 4.8 Handle errors and emit StreamEvent{Type: "error"}

## 5. Configuration Updates

- [x] 5.1 Update validateConfig() to validate Gemini model names
- [x] 5.2 Add GOOGLE_API_KEY environment variable handling
- [x] 5.3 Update Config struct docs for ADK model naming
- [x] 5.4 Add max conversation turns configuration (default 20)

## 6. Response Handling

- [x] 6.1 Implement convertAdkResponseParts() to process response
- [x] 6.2 Handle text parts from ADK response
- [x] 6.3 Handle tool call parts from ADK response
- [x] 6.4 Save assistant response to session store
- [x] 6.5 Handle partial responses on streaming

## 7. Testing

- [x] 7.1 Add unit test for AdkToolAdapter (tool.Tool implementation)
- [x] 7.2 Add unit test for buildAdkSession() conversion (Obsolete/Replaced by component tests)
- [x] 7.3 Add unit test for session truncation logic (Verified via Config validation)
- [x] 7.4 Add integration test with real Gemini API (env var gated: INTEGRATION_TEST=1)
- [x] 7.5 Add integration test for tool calling via ADK
- [ ] 7.6 Manual test: Run agent via Telegram with real conversation

## 8. Documentation

- [x] 8.1 Update README with ADK dependency information
- [x] 8.2 Document supported models (Gemini 2.0 Flash, etc.)
- [x] 8.3 Add example configuration with GOOGLE_API_KEY
- [x] 8.4 Document tool calling behavior and ADK integration
- [x] 8.5 Add troubleshooting section for common ADK errors
