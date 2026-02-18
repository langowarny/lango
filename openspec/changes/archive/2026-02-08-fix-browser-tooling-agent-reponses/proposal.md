## Why

The current Lango service has limited browser integration, only exposing navigation without returning content or allowing for reading/screenshotting. Additionally, the agent runtime lacks a system prompt, which leads small models like LFM 2.5 to default to generic refusals regarding "web access" even when tools are available. This change aims to make browser interactions functional and ensure the agent correctly identifies its capabilities.

## What Changes

- **Tool Registration**: `browser_read` and `browser_screenshot` will be registered with the agent.
- **Enhanced Navigation**: `browser_navigate` will be updated to return immediate page metadata (title and a text snippet) to provide instant context to the LLM.
- **System Prompt Support**: The agent runtime will be updated to support a configurable system prompt, allowing the service to explicitly inform the model about its identity and available tools.
- **Default Guidance**: A default system prompt will be added to guide the model on how to use the browser, crypto, and secrets tools effectively.

## Capabilities

### New Capabilities
- `browser-automation`: Provides full browser control including navigation, reading page content, and capturing screenshots.
- `agent-prompting`: Support for system prompts to guide agent behavior and tool usage.

### Modified Capabilities
- `browser-navigate`: Enhanced to return page metadata instead of just a success message.

## Impact

- `internal/app/app.go`: Tool registration and default prompt loading.
- `internal/agent/runtime.go`: Runtime logic for system prompt handling.
- `internal/tools/browser/browser.go`: Potential helper additions for snapshots.
