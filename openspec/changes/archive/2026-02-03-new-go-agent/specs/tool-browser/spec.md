## ADDED Requirements

### Requirement: Browser session management
The system SHALL manage Chrome/Chromium browser sessions using the CDP protocol via rod.

#### Scenario: Launch browser session
- **WHEN** a browser action is requested without active session
- **THEN** the system SHALL launch a headless browser instance

#### Scenario: Reuse existing session
- **WHEN** a browser action is requested with active session
- **THEN** the existing session SHALL be reused

#### Scenario: Session cleanup
- **WHEN** the session expires or is explicitly closed
- **THEN** browser resources SHALL be released

### Requirement: Page navigation
The system SHALL navigate browser tabs to specified URLs.

#### Scenario: Navigate to URL
- **WHEN** navigation to a URL is requested
- **THEN** the page SHALL load the specified URL

#### Scenario: Wait for page load
- **WHEN** navigation completes
- **THEN** the system SHALL wait for DOM ready state

### Requirement: Screenshot capture
The system SHALL capture screenshots of web pages or specific elements.

#### Scenario: Full page screenshot
- **WHEN** a full page screenshot is requested
- **THEN** the entire scrollable page SHALL be captured

#### Scenario: Element screenshot
- **WHEN** a screenshot of a specific element is requested
- **THEN** only that element SHALL be captured

### Requirement: DOM interaction
The system SHALL support clicking, typing, and extracting content from web pages.

#### Scenario: Click element
- **WHEN** a click on a selector is requested
- **THEN** the element SHALL be scrolled into view and clicked

#### Scenario: Type into input
- **WHEN** text input is requested
- **THEN** the text SHALL be typed into the focused element

#### Scenario: Get element text
- **WHEN** text extraction is requested
- **THEN** the text content of matching elements SHALL be returned

### Requirement: JavaScript execution
The system SHALL execute JavaScript code in the page context.

#### Scenario: Execute script
- **WHEN** JavaScript code is provided
- **THEN** the script SHALL be executed and the result returned
