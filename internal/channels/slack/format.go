package slack

import (
	"regexp"
	"strings"
)

// Package-level compiled regexes for Slack mrkdwn conversion.
var (
	boldRegex    = regexp.MustCompile(`\*\*(.+?)\*\*`)
	strikeRegex  = regexp.MustCompile(`~~(.+?)~~`)
	linkRegex    = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	headingRegex = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
)

// FormatMrkdwn converts standard Markdown to Slack mrkdwn format.
//
// Conversions:
//   - **bold** → *bold*
//   - ~~strike~~ → ~strike~
//   - [text](url) → <url|text>
//   - # Heading → *Heading*
//
// Content inside code blocks (```) is preserved without transformation.
func FormatMrkdwn(text string) string {
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

		result.WriteString(convertMrkdwnLine(line))
	}

	return result.String()
}

// convertMrkdwnLine transforms a single line from standard Markdown to Slack mrkdwn.
func convertMrkdwnLine(line string) string {
	// Headings → bold
	if m := headingRegex.FindStringSubmatch(line); m != nil {
		return "*" + m[2] + "*"
	}

	// Bold: **text** → *text*
	line = boldRegex.ReplaceAllString(line, "*$1*")

	// Strikethrough: ~~text~~ → ~text~
	line = strikeRegex.ReplaceAllString(line, "~$1~")

	// Links: [text](url) → <url|text>
	line = linkRegex.ReplaceAllString(line, "<$2|$1>")

	return line
}
