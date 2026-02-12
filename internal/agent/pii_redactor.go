package agent

import (
	"regexp"
)

// PIIConfig defines configuration for PII redaction
type PIIConfig struct {
	RedactEmail bool
	RedactPhone bool
	CustomRegex []string
}

// PIIRedactor redacts PII from input strings.
type PIIRedactor struct {
	config  PIIConfig
	regexes []*regexp.Regexp
}

// NewPIIRedactor creates a new PIIRedactor.
func NewPIIRedactor(cfg PIIConfig) *PIIRedactor {
	var regexes []*regexp.Regexp

	if cfg.RedactEmail {
		regexes = append(regexes, regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`))
	}
	if cfg.RedactPhone {
		regexes = append(regexes, regexp.MustCompile(`\b\d{3}-\d{3}-\d{4}\b`))
	}

	for _, pattern := range cfg.CustomRegex {
		if r, err := regexp.Compile(pattern); err == nil {
			regexes = append(regexes, r)
		}
	}

	return &PIIRedactor{
		config:  cfg,
		regexes: regexes,
	}
}

// RedactInput applies PII redaction patterns to an input string.
func (r *PIIRedactor) RedactInput(input string) string {
	redacted := input
	for _, re := range r.regexes {
		redacted = re.ReplaceAllString(redacted, "[REDACTED]")
	}
	return redacted
}
