package agent

import (
	"context"
	"testing"
)

func TestAdkToolAdapter(t *testing.T) {
	// Setup
	called := false
	handler := func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
		called = true
		if val, ok := params["test"]; ok {
			return val, nil
		}
		return "executed", nil
	}

	tool := &Tool{
		Name:        "test_tool",
		Description: "A test tool",
		Handler:     handler,
	}

	adapter := &AdkToolAdapter{tool: tool}

	// Test Name
	if adapter.Name() != "test_tool" {
		t.Errorf("expected name 'test_tool', got '%s'", adapter.Name())
	}

	// Test Description
	if adapter.Description() != "A test tool" {
		t.Errorf("expected description 'A test tool', got '%s'", adapter.Description())
	}

	// Test IsLongRunning
	if adapter.IsLongRunning() {
		t.Error("expected IsLongRunning to be false")
	}

	// Test Run
	ctx := context.Background()
	input := map[string]interface{}{
		"test": "value",
	}

	result, err := adapter.Run(ctx, input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("handler was not called")
	}
	if result != "value" {
		t.Errorf("expected result 'value', got '%v'", result)
	}

	// Test Run with nil input (should be safe)
	called = false
	_, err = adapter.Run(ctx, nil)
	if err != nil {
		t.Errorf("unexpected error with nil input: %v", err)
	}
	if !called {
		t.Error("handler was not called on nil input")
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Provider:             "gemini",
				Model:                "gemini-2.0-flash-exp",
				MaxConversationTurns: 20,
			},
			wantErr: false,
		},
		{
			name: "missing provider",
			config: Config{
				Model: "gemini-2.0-flash-exp",
			},
			wantErr: true,
		},
		{
			name: "missing model",
			config: Config{
				Provider: "gemini",
			},
			wantErr: true,
		},
		// 		{
		// 			name: "missing api key",
		// 			config: Config{
		// 				Provider: "gemini",
		// 				Model:    "gemini-2.0-flash-exp",
		// 			},
		// 			wantErr: true,
		// 		},
		{
			name: "negative max conversation turns",
			config: Config{
				Provider:             "gemini",
				Model:                "gemini-2.0-flash-exp",
				MaxConversationTurns: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateConfig(tt.config); (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
