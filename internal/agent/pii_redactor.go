package agent

import (
	"context"
	"regexp"
)

// PIIRedactor implements AgentRuntime to redact PII from inputs.
type PIIRedactor struct {
	BaseRuntimeMiddleware
	Config  PIIConfig
	regexes []*regexp.Regexp
}

// NewPIIRedactor creates a new PIIRedactor middleware.
func NewPIIRedactor(cfg PIIConfig) RuntimeMiddleware {
	var regexes []*regexp.Regexp

	// Default patterns
	if cfg.RedactEmail {
		// Simple email regex
		regexes = append(regexes, regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`))
	}
	if cfg.RedactPhone {
		// Simple phone regex (US context mostly)
		regexes = append(regexes, regexp.MustCompile(`\b\d{3}-\d{3}-\d{4}\b`))
	}

	for _, pattern := range cfg.CustomRegex {
		if r, err := regexp.Compile(pattern); err == nil {
			regexes = append(regexes, r)
		}
	}

	return func(next AgentRuntime) AgentRuntime {
		return &PIIRedactor{
			BaseRuntimeMiddleware: BaseRuntimeMiddleware{Next: next},
			Config:                cfg,
			regexes:               regexes,
		}
	}
}

// Run intercepts the execution to redact input.
func (mw *PIIRedactor) Run(ctx context.Context, sessionKey string, input string, events chan<- StreamEvent) error {
	redactedInput := input
	for _, r := range mw.regexes {
		redactedInput = r.ReplaceAllString(redactedInput, "[REDACTED]")
	}

	// Log redaction if changes occurred
	if redactedInput != input {
		// In a real system, we might want to log that redaction happened
		// logger.Debugw("redacted input", "original_len", len(input), "redacted_len", len(redactedInput))
	}

	return mw.Next.Run(ctx, sessionKey, redactedInput, events)
}
