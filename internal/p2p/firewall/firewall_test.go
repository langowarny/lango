package firewall

import (
	"testing"

	"go.uber.org/zap"
)

func TestValidateRule_AllowWildcardPeerAndTools(t *testing.T) {
	tests := []struct {
		give    ACLRule
		wantErr bool
	}{
		{
			give:    ACLRule{PeerDID: "*", Action: "allow"},
			wantErr: true, // wildcard peer + empty tools (= all)
		},
		{
			give:    ACLRule{PeerDID: "*", Action: "allow", Tools: []string{"*"}},
			wantErr: true, // wildcard peer + wildcard tool
		},
		{
			give:    ACLRule{PeerDID: "*", Action: "allow", Tools: []string{"echo", "*"}},
			wantErr: true, // wildcard tool mixed in
		},
		{
			give:    ACLRule{PeerDID: "*", Action: "deny"},
			wantErr: false, // deny rules always safe
		},
		{
			give:    ACLRule{PeerDID: "*", Action: "deny", Tools: []string{"*"}},
			wantErr: false, // deny rules always safe
		},
		{
			give:    ACLRule{PeerDID: "did:key:specific", Action: "allow", Tools: []string{"*"}},
			wantErr: false, // specific peer OK
		},
		{
			give:    ACLRule{PeerDID: "*", Action: "allow", Tools: []string{"echo"}},
			wantErr: false, // specific tool OK
		},
		{
			give:    ACLRule{PeerDID: "did:key:abc", Action: "allow"},
			wantErr: false, // specific peer, all tools
		},
	}

	for _, tt := range tests {
		t.Run(tt.give.PeerDID+"/"+tt.give.Action, func(t *testing.T) {
			err := ValidateRule(tt.give)
			if tt.wantErr && err == nil {
				t.Error("expected error for overly permissive rule")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestAddRule_RejectsOverlyPermissive(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	fw := New(nil, logger.Sugar())

	err := fw.AddRule(ACLRule{PeerDID: "*", Action: "allow", Tools: []string{"*"}})
	if err == nil {
		t.Error("expected AddRule to reject wildcard allow rule")
	}

	// Verify the rule was NOT added.
	rules := fw.Rules()
	if len(rules) != 0 {
		t.Errorf("expected no rules, got %d", len(rules))
	}
}

func TestAddRule_AcceptsValidRule(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	fw := New(nil, logger.Sugar())

	err := fw.AddRule(ACLRule{PeerDID: "did:key:peer-1", Action: "allow", Tools: []string{"echo"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rules := fw.Rules()
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].PeerDID != "did:key:peer-1" {
		t.Errorf("unexpected peer DID: %s", rules[0].PeerDID)
	}
}

func TestAddRule_AcceptsDenyWildcard(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	fw := New(nil, logger.Sugar())

	err := fw.AddRule(ACLRule{PeerDID: "*", Action: "deny", Tools: []string{"*"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rules := fw.Rules()
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
}

func TestNew_WarnsOnOverlyPermissiveInitialRules(t *testing.T) {
	// Should not panic — just logs a warning for backward compatibility.
	logger, _ := zap.NewDevelopment()
	fw := New([]ACLRule{
		{PeerDID: "*", Action: "allow"},
	}, logger.Sugar())

	// Rule is still loaded (backward compat).
	rules := fw.Rules()
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule (backward compat), got %d", len(rules))
	}
}
