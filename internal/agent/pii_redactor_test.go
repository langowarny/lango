package agent

import (
	"context"
	"testing"
)

type mockRuntime struct {
	CapturedInput string
}

func (m *mockRuntime) Run(ctx context.Context, sessionKey string, input string, events chan<- StreamEvent) error {
	m.CapturedInput = input
	close(events)
	return nil
}

func (m *mockRuntime) RegisterTool(tool *Tool) error     { return nil }
func (m *mockRuntime) GetTool(name string) (*Tool, bool) { return nil, false }
func (m *mockRuntime) ListTools() []*Tool                { return nil }
func (m *mockRuntime) ExecuteTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	return nil, nil
}

func TestPIIRedactor(t *testing.T) {
	cfg := PIIConfig{
		RedactEmail: true,
		RedactPhone: true,
	}

	redactor := NewPIIRedactor(cfg)
	mock := &mockRuntime{}
	wrapped := redactor(mock)

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
			events := make(chan StreamEvent)
			// Run is async/blocking depending on implementation, here mock is sync but redactor calls Next.Run
			go func() {
				for range events {
				}
			}()

			err := wrapped.Run(context.Background(), "sess", tt.input, events)
			if err != nil {
				t.Fatalf("Run failed: %v", err)
			}

			if mock.CapturedInput != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, mock.CapturedInput)
			}
		})
	}
}
