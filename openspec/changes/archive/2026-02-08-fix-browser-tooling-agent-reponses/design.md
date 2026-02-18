## Context

The Lango agent currently operates without a system prompt and with limited browser tooling (navigation only). This design introduces full browser control and system-level guidance to improve agent effectiveness, especially for light models like LFM.

## Goals / Non-Goals

**Goals:**
- Enable the agent to read web content and capture screenshots.
- Provide immediate feedback (context) after navigation.
- Implement system prompt support in the agent runtime.
- Guide the agent via a default system prompt to use its tools.

**Non-Goals:**
- Changing the underlying browser library (rod).
- Implementing multi-tab support (currently one page/session).
- Adding complex JS interaction tools (evolved interactions like hovering/scrolling).

## Decisions

### 1. Expanded Browser Toolset
We will register `browser_read` and `browser_screenshot` handlers in `internal/app/app.go`. These will use the existing `browser.Tool` methods (`GetText` and `Screenshot`).

### 2. Contextual Navigation
`browser_navigate` will be modified to return a JSON string containing the page title and the first 1000 characters of the page text. This ensures the model doesn't "fly blind" after a successful navigation.

### 3. Agent Runtime System Prompt
- Add `SystemPrompt string` to `agent.Config`.
- In `Runtime.Run`, if `sess.History` is empty (new session), prepend a message with `role: system` and `content: r.config.SystemPrompt`.
- This ensures the model is aware of its capabilities (Web, Crypto, Secrets) from the start.

### 4. Default System Prompt
The default prompt will be defined in `internal/app/app.go` (or loaded from config) and will state:
"You are Antigravity, a powerful AI assistant. You have access to tools for web navigation (browser), secure secrets management (secrets), and cryptographic operations (crypto). Use them when appropriate."

## Risks / Trade-offs

- **Token Usage**: Returning content from `browser_navigate` and adding a system prompt increases token consumption. We will limit the snippet size to mitigate this.
- **Model Sensitivity**: Small models might overvalue the system prompt. We need to ensure the prompt is concise.
