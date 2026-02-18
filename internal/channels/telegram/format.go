package telegram

import (
	"strings"
)

// FormatMarkdown converts standard Markdown to Telegram Markdown v1.
//
// Telegram v1 supports: *bold*, _italic_, `code`, ```pre```, [text](url).
// Standard Markdown features not in v1 (e.g. ~~strike~~) are stripped.
// Content inside code blocks (```) is preserved without transformation.
func FormatMarkdown(text string) string {
	var result strings.Builder
	lines := strings.Split(text, "\n")
	inCodeBlock := false

	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}

		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCodeBlock = !inCodeBlock
			result.WriteString(line)
			continue
		}

		if inCodeBlock {
			result.WriteString(line)
			continue
		}

		result.WriteString(convertLine(line))
	}

	return result.String()
}

// convertLine transforms a single line from standard Markdown to Telegram v1.
func convertLine(line string) string {
	// Headings → bold: "# Heading" → "*Heading*"
	trimmed := strings.TrimSpace(line)
	if len(trimmed) > 0 && trimmed[0] == '#' {
		heading := strings.TrimLeft(trimmed, "#")
		heading = strings.TrimSpace(heading)
		if heading != "" {
			return "*" + heading + "*"
		}
		return line
	}

	// Bold: **text** → *text*
	line = convertBold(line)

	// Strikethrough: ~~text~~ → text (v1 doesn't support it)
	line = convertStrike(line)

	return line
}

// convertBold replaces **text** with *text*, preserving inline code spans.
func convertBold(line string) string {
	var result strings.Builder
	i := 0
	for i < len(line) {
		// Skip inline code
		if line[i] == '`' {
			end := strings.Index(line[i+1:], "`")
			if end >= 0 {
				result.WriteString(line[i : i+end+2])
				i += end + 2
				continue
			}
		}

		// Match **...**
		if i+1 < len(line) && line[i] == '*' && line[i+1] == '*' {
			end := strings.Index(line[i+2:], "**")
			if end >= 0 {
				inner := line[i+2 : i+2+end]
				result.WriteString("*" + inner + "*")
				i += 2 + end + 2
				continue
			}
		}

		result.WriteByte(line[i])
		i++
	}
	return result.String()
}

// convertStrike removes ~~text~~ markers (Telegram v1 doesn't support strikethrough).
func convertStrike(line string) string {
	var result strings.Builder
	i := 0
	for i < len(line) {
		// Skip inline code
		if line[i] == '`' {
			end := strings.Index(line[i+1:], "`")
			if end >= 0 {
				result.WriteString(line[i : i+end+2])
				i += end + 2
				continue
			}
		}

		// Match ~~...~~
		if i+1 < len(line) && line[i] == '~' && line[i+1] == '~' {
			end := strings.Index(line[i+2:], "~~")
			if end >= 0 {
				inner := line[i+2 : i+2+end]
				result.WriteString(inner)
				i += 2 + end + 2
				continue
			}
		}

		result.WriteByte(line[i])
		i++
	}
	return result.String()
}
