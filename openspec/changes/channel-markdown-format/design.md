## Context

LLM agents produce standard Markdown output, but messaging platforms use different markup dialects. Telegram supports Markdown v1 (limited subset) and MarkdownV2 (requires extensive escaping). Slack uses mrkdwn (similar but different syntax). Discord natively supports standard Markdown.

Currently, `Send()` methods pass raw text to platform APIs. Telegram receives no `ParseMode`, so messages render as plain text. Slack receives unformatted text, causing partially broken formatting.

## Goals / Non-Goals

**Goals:**
- Auto-convert standard Markdown to each platform's native format at the `Send()` layer
- Preserve code blocks without transformation across all platforms
- Provide plain text fallback for Telegram when formatted messages fail API parsing
- Zero changes to handler logic in `internal/app/channels.go`

**Non-Goals:**
- MarkdownV2 support (too fragile with LLM output requiring 18 special character escapes)
- Rich media formatting (images, buttons, cards)
- Block Kit conversion for Slack (existing Block Kit support is untouched)
- Discord formatting changes (native Markdown support, no conversion needed)

## Decisions

### 1. Telegram v1 over MarkdownV2
**Decision**: Use Markdown v1 (`ParseMode: "Markdown"`) instead of MarkdownV2.

**Rationale**: MarkdownV2 requires escaping 18 special characters (`_`, `*`, `[`, `]`, `(`, `)`, `~`, `` ` ``, `>`, `#`, `+`, `-`, `=`, `|`, `{`, `}`, `.`, `!`) outside of code blocks. LLM output contains these characters naturally in math, lists, and prose. V1 is simpler and covers the primary use cases (bold, italic, code, links).

**Alternative**: MarkdownV2 with full escaping — rejected due to fragility and maintenance cost.

### 2. Send-level conversion
**Decision**: Apply format conversion inside each channel's `Send()` method, not in handlers.

**Rationale**: Single point of transformation. Handlers remain platform-agnostic. Any code path calling `Send()` gets formatting automatically.

**Alternative**: Middleware/decorator pattern — rejected as over-engineering for two simple converters.

### 3. Plain text fallback for Telegram
**Decision**: On Telegram API parse error, re-send the original (unconverted) text without ParseMode.

**Rationale**: LLM output may contain edge-case Markdown that breaks v1 parsing. Delivering unformatted text is better than failing silently or showing error messages to users.

### 4. Code block preservation
**Decision**: Track ` ``` ` fences and skip all conversions inside code blocks.

**Rationale**: Code snippets should never have their formatting altered. Both platforms support ` ``` ` natively.

## Risks / Trade-offs

- [Telegram v1 limitations] Bold uses single `*`, conflicting with italic on some edge cases → Acceptable trade-off; v1 has worked for years in Telegram bots
- [Regex-based Slack conversion] Could break on deeply nested or malformed Markdown → Mitigated by simple non-greedy regex patterns and code block skip logic
- [Fallback hides formatting] Users get plain text instead of formatted on fallback → Better than error messages; logged as warning for debugging
