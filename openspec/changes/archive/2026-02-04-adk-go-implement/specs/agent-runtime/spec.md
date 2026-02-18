## MODIFIED Requirements

### Requirement: Agent initialization with ADK
The system SHALL initialize an LLM agent using the ADK framework (`google.golang.org/adk`) with configurable model provider and parameters.

#### Scenario: Successful agent initialization with Gemini
- **WHEN** the application starts with a valid GOOGLE_API_KEY environment variable
- **THEN** the agent runtime SHALL create an ADK llmagent instance using gemini.NewModel() and llmagent.New()
- **AND** the agent SHALL be ready to process messages

#### Scenario: Missing API credentials
- **WHEN** the application starts without required API credentials
- **THEN** the system SHALL return a clear error message indicating which credentials are missing

#### Scenario: Invalid model name
- **WHEN** configuration specifies an invalid or unsupported Gemini model name
- **THEN** the system SHALL return an error from the Gemini SDK with the invalid model name

## ADDED Requirements

### Requirement: Tool conversion to ADK format
The system SHALL automatically convert registered Tool instances to ADK tool.Tool interface implementations.

#### Scenario: Tool adapter creation
- **WHEN** a tool is registered with name, description, and handler
- **THEN** the system SHALL create an ADK tool adapter that implements tool.Tool interface
- **AND** the adapter SHALL expose Name(), Description(), and Run() methods

#### Scenario: Tool execution via ADK
- **WHEN** the LLM agent calls a tool through ADK
- **THEN** the system SHALL invoke the original Tool.Handler function
- **AND** SHALL return results to the ADK agent

### Requirement: ADK session management
The system SHALL convert internal session history to ADK session format for agent execution.

#### Scenario: Session history conversion
- **WHEN** loading conversation history for an existing session
- **THEN** the system SHALL create an ADK agent.Session with all previous messages
- **AND** SHALL preserve message roles (user, assistant, tool)

#### Scenario: Session history truncation
- **WHEN** conversation history exceeds reasonable message count
- **THEN** the system SHALL keep the most recent 20 conversation turns
- **AND** SHALL log a warning about truncation

### Requirement: Response streaming from ADK
The system SHALL execute ADK agent and emit StreamEvent types based on response parts.

#### Scenario: Text response streaming
- **WHEN** the ADK agent returns a response with text parts
- **THEN** the system SHALL emit StreamEvent with type "text_delta" for each text part

#### Scenario: Tool call in response
- **WHEN** the ADK agent response includes tool calls
- **THEN** the system SHALL emit StreamEvent with type "tool_start" when tool execution begins
- **AND** SHALL emit StreamEvent with type "tool_end" when tool execution completes

#### Scenario: Error handling
- **WHEN** the ADK agent execution fails or a tool call errors
- **THEN** the system SHALL emit StreamEvent with type "error" containing the error details
