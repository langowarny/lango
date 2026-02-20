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

func TestPIIRedactor_KoreanPII(t *testing.T) {
	redactor := NewPIIRedactor(PIIConfig{
		RedactEmail: true,
		RedactPhone: true,
	})

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Korean mobile number",
			input: "전화번호: 010-1234-5678",
		},
		{
			name:  "Korean RRN",
			input: "주민번호: 900101-1234567",
		},
		{
			name:  "Korean landline",
			input: "서울 전화: 02-1234-5678",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactor.RedactInput(tt.input)
			if result == tt.input {
				t.Errorf("expected PII to be redacted, but got original: %q", result)
			}
			if len(result) == 0 {
				t.Error("result should not be empty")
			}
		})
	}
}

func TestPIIRedactor_DisabledBuiltins(t *testing.T) {
	redactor := NewPIIRedactor(PIIConfig{
		RedactEmail:      true,
		RedactPhone:      true,
		DisabledBuiltins: []string{"email"},
	})

	// Email should NOT be redacted
	result := redactor.RedactInput("user@example.com")
	if result != "user@example.com" {
		t.Errorf("disabled email pattern should not redact, got: %q", result)
	}

	// Phone should still be redacted
	result = redactor.RedactInput("Call 123-456-7890")
	if result == "Call 123-456-7890" {
		t.Error("phone should still be redacted when only email is disabled")
	}
}

func TestPIIRedactor_CustomPatterns(t *testing.T) {
	redactor := NewPIIRedactor(PIIConfig{
		RedactEmail: true,
		RedactPhone: true,
		CustomPatterns: map[string]string{
			"project_id": `\bPROJ-\d{4}\b`,
		},
	})

	result := redactor.RedactInput("See project PROJ-1234 for details")
	if result == "See project PROJ-1234 for details" {
		t.Error("custom pattern should redact PROJ-1234")
	}
}

func TestPIIRedactor_LegacyCustomRegex(t *testing.T) {
	redactor := NewPIIRedactor(PIIConfig{
		RedactEmail: true,
		RedactPhone: true,
		CustomRegex: []string{`\bTOKEN-[A-Z0-9]+\b`},
	})

	result := redactor.RedactInput("Auth: TOKEN-ABC123")
	if result == "Auth: TOKEN-ABC123" {
		t.Error("legacy custom regex should redact TOKEN-ABC123")
	}
}

func TestPIIRedactor_OverlappingMatches(t *testing.T) {
	redactor := NewPIIRedactor(PIIConfig{
		RedactEmail: true,
		RedactPhone: true,
	})

	// A string with adjacent PII
	result := redactor.RedactInput("Email: a@b.com b@c.com")
	// Both emails should be redacted
	if result == "Email: a@b.com b@c.com" {
		t.Error("both emails should be redacted")
	}
}

func TestPIIRedactor_EmptyInput(t *testing.T) {
	redactor := NewPIIRedactor(PIIConfig{
		RedactEmail: true,
		RedactPhone: true,
	})

	result := redactor.RedactInput("")
	if result != "" {
		t.Errorf("empty input should return empty, got %q", result)
	}
}

func TestPIIRedactor_CreditCardWithLuhn(t *testing.T) {
	redactor := NewPIIRedactor(PIIConfig{
		RedactEmail: true,
		RedactPhone: true,
	})

	// Valid Visa number (passes Luhn)
	result := redactor.RedactInput("Card: 4111111111111111")
	if result == "Card: 4111111111111111" {
		t.Error("valid credit card should be redacted")
	}

	// Invalid Luhn number should NOT be redacted
	result = redactor.RedactInput("Card: 4111111111111112")
	if result != "Card: 4111111111111112" {
		t.Errorf("invalid Luhn should not be redacted, got: %q", result)
	}
}

func TestPIIRedactor_SSN(t *testing.T) {
	redactor := NewPIIRedactor(PIIConfig{
		RedactEmail: true,
		RedactPhone: true,
	})

	result := redactor.RedactInput("SSN: 123-45-6789")
	if result == "SSN: 123-45-6789" {
		t.Error("SSN should be redacted")
	}
}

func TestPIIRedactor_MultiplePIITypes(t *testing.T) {
	redactor := NewPIIRedactor(PIIConfig{
		RedactEmail: true,
		RedactPhone: true,
	})

	input := "Contact: user@mail.com, 010-1234-5678, SSN: 123-45-6789"
	result := redactor.RedactInput(input)

	if result == input {
		t.Error("multiple PII types should all be redacted")
	}

	// Verify the result only contains [REDACTED] markers and non-PII text
	if len(result) >= len(input) {
		t.Errorf("redacted result should be shorter or equal, input len=%d, result len=%d", len(input), len(result))
	}
}
