package onboard

import (
	"testing"

	"github.com/langoai/lango/internal/cli/tuicore"
	"github.com/langoai/lango/internal/config"
)

func fieldByKey(form *tuicore.FormModel, key string) *tuicore.Field {
	for _, f := range form.Fields {
		if f.Key == key {
			return f
		}
	}
	return nil
}

func TestNewProviderStepForm(t *testing.T) {
	cfg := config.DefaultConfig()
	form := NewProviderStepForm(cfg)

	wantKeys := []string{"type", "id", "apikey", "baseurl"}

	if len(form.Fields) != len(wantKeys) {
		t.Fatalf("expected %d fields, got %d", len(wantKeys), len(form.Fields))
	}

	for _, key := range wantKeys {
		if f := fieldByKey(form, key); f == nil {
			t.Errorf("missing field %q", key)
		}
	}

	if f := fieldByKey(form, "type"); f.Type != tuicore.InputSelect {
		t.Errorf("type: want InputSelect, got %d", f.Type)
	}
	if f := fieldByKey(form, "apikey"); f.Type != tuicore.InputPassword {
		t.Errorf("apikey: want InputPassword, got %d", f.Type)
	}
}

func TestNewAgentStepForm(t *testing.T) {
	cfg := config.DefaultConfig()
	form := NewAgentStepForm(cfg)

	wantKeys := []string{"provider", "model", "maxtokens", "temp"}

	if len(form.Fields) != len(wantKeys) {
		t.Fatalf("expected %d fields, got %d", len(wantKeys), len(form.Fields))
	}

	for _, key := range wantKeys {
		if f := fieldByKey(form, key); f == nil {
			t.Errorf("missing field %q", key)
		}
	}

	if f := fieldByKey(form, "provider"); f.Type != tuicore.InputSelect {
		t.Errorf("provider: want InputSelect, got %d", f.Type)
	}
}

func TestNewChannelStepForm_Telegram(t *testing.T) {
	cfg := config.DefaultConfig()
	form := NewChannelStepForm("telegram", cfg)
	if form == nil {
		t.Fatal("form should not be nil for telegram")
	}
	if len(form.Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(form.Fields))
	}
	if f := fieldByKey(form, "telegram_token"); f == nil {
		t.Error("missing telegram_token field")
	}
}

func TestNewChannelStepForm_Discord(t *testing.T) {
	cfg := config.DefaultConfig()
	form := NewChannelStepForm("discord", cfg)
	if form == nil {
		t.Fatal("form should not be nil for discord")
	}
	if len(form.Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(form.Fields))
	}
	if f := fieldByKey(form, "discord_token"); f == nil {
		t.Error("missing discord_token field")
	}
}

func TestNewChannelStepForm_Slack(t *testing.T) {
	cfg := config.DefaultConfig()
	form := NewChannelStepForm("slack", cfg)
	if form == nil {
		t.Fatal("form should not be nil for slack")
	}
	if len(form.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(form.Fields))
	}
	if f := fieldByKey(form, "slack_token"); f == nil {
		t.Error("missing slack_token field")
	}
	if f := fieldByKey(form, "slack_app_token"); f == nil {
		t.Error("missing slack_app_token field")
	}
}

func TestNewChannelStepForm_Skip(t *testing.T) {
	cfg := config.DefaultConfig()
	form := NewChannelStepForm("skip", cfg)
	if form != nil {
		t.Error("form should be nil for skip")
	}
}

func TestNewSecurityStepForm(t *testing.T) {
	cfg := config.DefaultConfig()
	form := NewSecurityStepForm(cfg)

	wantKeys := []string{"interceptor_enabled", "interceptor_pii", "interceptor_policy"}

	if len(form.Fields) != len(wantKeys) {
		t.Fatalf("expected %d fields, got %d", len(wantKeys), len(form.Fields))
	}

	for _, key := range wantKeys {
		if f := fieldByKey(form, key); f == nil {
			t.Errorf("missing field %q", key)
		}
	}
}

func TestSuggestModel(t *testing.T) {
	tests := []struct {
		give string
		want string
	}{
		{give: "anthropic", want: "claude-sonnet-4-5-20250929"},
		{give: "openai", want: "gpt-4o"},
		{give: "gemini", want: "gemini-2.0-flash"},
		{give: "ollama", want: "llama3.1"},
		{give: "github", want: "gpt-4o"},
		{give: "unknown", want: "claude-sonnet-4-5-20250929"},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			got := suggestModel(tt.give)
			if got != tt.want {
				t.Errorf("suggestModel(%q): want %q, got %q", tt.give, tt.want, got)
			}
		})
	}
}

// TestAllFormsHaveDescriptions verifies every field across all forms has a non-empty Description.
func TestAllFormsHaveDescriptions(t *testing.T) {
	cfg := config.DefaultConfig()

	forms := map[string]*tuicore.FormModel{
		"provider": NewProviderStepForm(cfg),
		"agent":    NewAgentStepForm(cfg),
		"telegram": NewChannelStepForm("telegram", cfg),
		"discord":  NewChannelStepForm("discord", cfg),
		"slack":    NewChannelStepForm("slack", cfg),
		"security": NewSecurityStepForm(cfg),
	}

	for name, form := range forms {
		for _, f := range form.Fields {
			if f.Description == "" {
				t.Errorf("form %q field %q has empty Description", name, f.Key)
			}
		}
	}
}

// TestProviderOptionsIncludeGitHub verifies "github" is in provider options.
func TestProviderOptionsIncludeGitHub(t *testing.T) {
	cfg := config.DefaultConfig()

	// Check Provider Step form
	form := NewProviderStepForm(cfg)
	typeField := fieldByKey(form, "type")
	if typeField == nil {
		t.Fatal("missing type field")
	}
	found := false
	for _, opt := range typeField.Options {
		if opt == "github" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("provider type options missing 'github': %v", typeField.Options)
	}

	// Check fallback provider options list
	opts := buildProviderOptions(cfg)
	found = false
	for _, opt := range opts {
		if opt == "github" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("buildProviderOptions fallback missing 'github': %v", opts)
	}
}

// TestTemperatureValidator verifies temperature range validation.
func TestTemperatureValidator(t *testing.T) {
	cfg := config.DefaultConfig()
	form := NewAgentStepForm(cfg)
	tempField := fieldByKey(form, "temp")
	if tempField == nil {
		t.Fatal("missing temp field")
	}
	if tempField.Validate == nil {
		t.Fatal("temp field has no Validate function")
	}

	tests := []struct {
		give    string
		wantErr bool
	}{
		{give: "0.0", wantErr: false},
		{give: "1.5", wantErr: false},
		{give: "2.0", wantErr: false},
		{give: "2.1", wantErr: true},
		{give: "-0.1", wantErr: true},
		{give: "abc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			err := tempField.Validate(tt.give)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate(%q): wantErr=%v, got %v", tt.give, tt.wantErr, err)
			}
		})
	}
}

// TestMaxTokensValidator verifies max tokens positive integer validation.
func TestMaxTokensValidator(t *testing.T) {
	cfg := config.DefaultConfig()
	form := NewAgentStepForm(cfg)
	mtField := fieldByKey(form, "maxtokens")
	if mtField == nil {
		t.Fatal("missing maxtokens field")
	}
	if mtField.Validate == nil {
		t.Fatal("maxtokens field has no Validate function")
	}

	tests := []struct {
		give    string
		wantErr bool
	}{
		{give: "4096", wantErr: false},
		{give: "1", wantErr: false},
		{give: "0", wantErr: true},
		{give: "-1", wantErr: true},
		{give: "abc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			err := mtField.Validate(tt.give)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate(%q): wantErr=%v, got %v", tt.give, tt.wantErr, err)
			}
		})
	}
}

// TestSecurityConditionalVisibility verifies interceptor sub-fields are conditionally visible.
func TestSecurityConditionalVisibility(t *testing.T) {
	cfg := config.DefaultConfig()
	form := NewSecurityStepForm(cfg)

	enabledField := fieldByKey(form, "interceptor_enabled")
	piiField := fieldByKey(form, "interceptor_pii")
	policyField := fieldByKey(form, "interceptor_policy")

	if enabledField == nil || piiField == nil || policyField == nil {
		t.Fatal("missing security fields")
	}

	// When interceptor is disabled, sub-fields should be hidden
	enabledField.Checked = false
	if piiField.IsVisible() {
		t.Error("interceptor_pii should be hidden when interceptor is disabled")
	}
	if policyField.IsVisible() {
		t.Error("interceptor_policy should be hidden when interceptor is disabled")
	}

	// Count visible fields when disabled
	visibleCount := 0
	for _, f := range form.Fields {
		if f.IsVisible() {
			visibleCount++
		}
	}
	if visibleCount != 1 {
		t.Errorf("expected 1 visible field when interceptor disabled, got %d", visibleCount)
	}

	// When interceptor is enabled, all fields should be visible
	enabledField.Checked = true
	if !piiField.IsVisible() {
		t.Error("interceptor_pii should be visible when interceptor is enabled")
	}
	if !policyField.IsVisible() {
		t.Error("interceptor_policy should be visible when interceptor is enabled")
	}

	visibleCount = 0
	for _, f := range form.Fields {
		if f.IsVisible() {
			visibleCount++
		}
	}
	if visibleCount != 3 {
		t.Errorf("expected 3 visible fields when interceptor enabled, got %d", visibleCount)
	}
}
