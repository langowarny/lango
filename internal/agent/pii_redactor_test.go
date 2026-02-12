package agent

import (
	"testing"
)

func TestPIIRedactor(t *testing.T) {
	cfg := PIIConfig{
		RedactEmail: true,
		RedactPhone: true,
	}

	redactor := NewPIIRedactor(cfg)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No PII",
			input:    "Hello world",
			expected: "Hello world",
		},
		{
			name:     "Email redaction",
			input:    "My email is test@example.com contact me",
			expected: "My email is [REDACTED] contact me",
		},
		{
			name:     "Phone redaction",
			input:    "Call 123-456-7890 now",
			expected: "Call [REDACTED] now",
		},
		{
			name:     "Mixed PII",
			input:    "Email: bob@mail.com, Phone: 555-123-4567",
			expected: "Email: [REDACTED], Phone: [REDACTED]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactor.RedactInput(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
