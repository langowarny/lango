package workflow

import (
	"fmt"
	"regexp"
	"strings"
)

// placeholderRe matches {{stepID.result}} patterns in prompt templates.
// Step IDs may contain letters, digits, hyphens, and underscores.
var placeholderRe = regexp.MustCompile(`\{\{([a-zA-Z0-9_-]+)\.result\}\}`)

// RenderPrompt substitutes {{stepID.result}} placeholders in a prompt template
// with actual results from previous steps. It returns an error if a referenced
// step has no result available.
func RenderPrompt(tmpl string, results map[string]string) (string, error) {
	var missingKeys []string

	rendered := placeholderRe.ReplaceAllStringFunc(tmpl, func(match string) string {
		sub := placeholderRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		stepID := sub[1]
		val, ok := results[stepID]
		if !ok {
			missingKeys = append(missingKeys, stepID)
			return match
		}
		return val
	})

	if len(missingKeys) > 0 {
		return "", fmt.Errorf("missing results for steps: %s", strings.Join(missingKeys, ", "))
	}

	return rendered, nil
}
