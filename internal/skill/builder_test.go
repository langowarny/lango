package skill

import (
	"testing"
)

func TestBuildCompositeSkill(t *testing.T) {
	t.Run("basic fields and steps conversion", func(t *testing.T) {
		steps := []SkillStep{
			{Tool: "read", Params: map[string]interface{}{"path": "/tmp"}},
			{Tool: "write", Params: map[string]interface{}{"path": "/out"}},
		}
		got := BuildCompositeSkill("my-skill", "does things", steps, nil)

		if got.Name != "my-skill" {
			t.Errorf("Name = %q, want %q", got.Name, "my-skill")
		}
		if got.Description != "does things" {
			t.Errorf("Description = %q, want %q", got.Description, "does things")
		}
		if got.Type != "composite" {
			t.Errorf("Type = %q, want %q", got.Type, "composite")
		}
		if !got.RequiresApproval {
			t.Error("RequiresApproval = false, want true")
		}

		stepDefs, ok := got.Definition["steps"].([]interface{})
		if !ok {
			t.Fatalf("Definition[\"steps\"] is %T, want []interface{}", got.Definition["steps"])
		}
		if len(stepDefs) != 2 {
			t.Fatalf("len(steps) = %d, want 2", len(stepDefs))
		}

		first, ok := stepDefs[0].(map[string]interface{})
		if !ok {
			t.Fatalf("stepDefs[0] is %T, want map[string]interface{}", stepDefs[0])
		}
		if first["tool"] != "read" {
			t.Errorf("stepDefs[0][\"tool\"] = %v, want %q", first["tool"], "read")
		}
	})

	t.Run("nil params leaves Parameters nil", func(t *testing.T) {
		got := BuildCompositeSkill("s", "d", nil, nil)
		if got.Parameters != nil {
			t.Errorf("Parameters = %v, want nil", got.Parameters)
		}
	})

	t.Run("non-nil params sets Parameters", func(t *testing.T) {
		params := map[string]interface{}{"key": "value"}
		got := BuildCompositeSkill("s", "d", nil, params)
		if got.Parameters == nil {
			t.Fatal("Parameters is nil, want non-nil")
		}
		if got.Parameters["key"] != "value" {
			t.Errorf("Parameters[\"key\"] = %v, want %q", got.Parameters["key"], "value")
		}
	})
}

func TestBuildScriptSkill(t *testing.T) {
	tests := []struct {
		give       string
		giveScript string
		giveParams map[string]interface{}
	}{
		{
			give:       "with params",
			giveScript: "echo hello",
			giveParams: map[string]interface{}{"env": "prod"},
		},
		{
			give:       "without params",
			giveScript: "ls -la",
			giveParams: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			got := BuildScriptSkill("run", "runs script", tt.giveScript, tt.giveParams)

			if got.Type != "script" {
				t.Errorf("Type = %q, want %q", got.Type, "script")
			}
			if !got.RequiresApproval {
				t.Error("RequiresApproval = false, want true")
			}

			script, ok := got.Definition["script"].(string)
			if !ok {
				t.Fatalf("Definition[\"script\"] is %T, want string", got.Definition["script"])
			}
			if script != tt.giveScript {
				t.Errorf("Definition[\"script\"] = %q, want %q", script, tt.giveScript)
			}

			if tt.giveParams != nil && got.Parameters == nil {
				t.Error("Parameters is nil, want non-nil")
			}
			if tt.giveParams == nil && got.Parameters != nil {
				t.Errorf("Parameters = %v, want nil", got.Parameters)
			}
		})
	}
}

func TestBuildTemplateSkill(t *testing.T) {
	tests := []struct {
		give         string
		giveTemplate string
		giveParams   map[string]interface{}
	}{
		{
			give:         "with params",
			giveTemplate: "Hello {{.Name}}",
			giveParams:   map[string]interface{}{"name": "string"},
		},
		{
			give:         "without params",
			giveTemplate: "static template",
			giveParams:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			got := BuildTemplateSkill("tmpl", "renders template", tt.giveTemplate, tt.giveParams)

			if got.Type != "template" {
				t.Errorf("Type = %q, want %q", got.Type, "template")
			}
			if !got.RequiresApproval {
				t.Error("RequiresApproval = false, want true")
			}

			tmpl, ok := got.Definition["template"].(string)
			if !ok {
				t.Fatalf("Definition[\"template\"] is %T, want string", got.Definition["template"])
			}
			if tmpl != tt.giveTemplate {
				t.Errorf("Definition[\"template\"] = %q, want %q", tmpl, tt.giveTemplate)
			}

			if tt.giveParams != nil && got.Parameters == nil {
				t.Error("Parameters is nil, want non-nil")
			}
			if tt.giveParams == nil && got.Parameters != nil {
				t.Errorf("Parameters = %v, want nil", got.Parameters)
			}
		})
	}
}
