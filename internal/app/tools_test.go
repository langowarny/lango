package app

import (
	"strings"
	"testing"
)

func TestBlockLangoExec_SkillGuards(t *testing.T) {
	tests := []struct {
		give    string
		wantMsg bool
	}{
		{
			give:    "git clone https://github.com/owner/skill-repo",
			wantMsg: true,
		},
		{
			give:    "Git Clone https://github.com/owner/skills",
			wantMsg: true,
		},
		{
			give:    "curl https://example.com/skill.md",
			wantMsg: true,
		},
		{
			give:    "wget https://example.com/skills/SKILL.md",
			wantMsg: true,
		},
		{
			give:    "git clone https://github.com/owner/unrelated-repo",
			wantMsg: false,
		},
		{
			give:    "curl https://example.com/api/data",
			wantMsg: false,
		},
		{
			give:    "ls -la",
			wantMsg: false,
		},
		{
			give:    "lango cron list",
			wantMsg: true,
		},
	}

	auto := map[string]bool{"cron": true, "background": true}
	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			msg := blockLangoExec(tt.give, auto)
			gotMsg := msg != ""
			if gotMsg != tt.wantMsg {
				t.Errorf("blockLangoExec(%q) returned msg=%q, wantMsg=%v", tt.give, msg, tt.wantMsg)
			}
		})
	}
}

func TestBlockLangoExec_AllSubcommands(t *testing.T) {
	auto := map[string]bool{"cron": true, "background": true, "workflow": true}

	tests := []struct {
		give        string
		wantBlocked bool
		wantContain string // substring expected in the message
	}{
		// Phase 1: subcommands with in-process tool equivalents
		{give: "lango cron list", wantBlocked: true, wantContain: "cron_"},
		{give: "lango bg submit", wantBlocked: true, wantContain: "bg_"},
		{give: "lango background list", wantBlocked: true, wantContain: "bg_"},
		{give: "lango workflow run", wantBlocked: true, wantContain: "workflow_"},
		{give: "lango graph query", wantBlocked: true, wantContain: "graph_"},
		{give: "lango memory list", wantBlocked: true, wantContain: "memory_"},
		{give: "lango p2p status", wantBlocked: true, wantContain: "p2p_"},
		{give: "lango security keyring status", wantBlocked: true, wantContain: "crypto_"},
		{give: "lango payment send", wantBlocked: true, wantContain: "payment_"},

		// Phase 2: catch-all for subcommands without in-process equivalents
		{give: "lango config list", wantBlocked: true, wantContain: "passphrase"},
		{give: "lango doctor", wantBlocked: true, wantContain: "passphrase"},
		{give: "lango serve", wantBlocked: true, wantContain: "passphrase"},
		{give: "lango settings", wantBlocked: true, wantContain: "passphrase"},
		{give: "lango onboard", wantBlocked: true, wantContain: "passphrase"},
		{give: "lango agent list", wantBlocked: true, wantContain: "passphrase"},
		{give: "lango", wantBlocked: true, wantContain: "passphrase"},
		{give: "LANGO SECURITY DB-MIGRATE", wantBlocked: true, wantContain: "crypto_"},

		// Not blocked: non-lango commands
		{give: "ls -la", wantBlocked: false},
		{give: "go build ./...", wantBlocked: false},
		{give: "echo lango", wantBlocked: false},
		{give: "cat lango.yaml", wantBlocked: false},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			msg := blockLangoExec(tt.give, auto)
			gotBlocked := msg != ""
			if gotBlocked != tt.wantBlocked {
				t.Errorf("blockLangoExec(%q): blocked=%v, want %v (msg=%q)",
					tt.give, gotBlocked, tt.wantBlocked, msg)
			}
			if tt.wantContain != "" && !strings.Contains(msg, tt.wantContain) {
				t.Errorf("blockLangoExec(%q): message %q does not contain %q",
					tt.give, msg, tt.wantContain)
			}
		})
	}
}

func TestBlockLangoExec_DisabledFeature(t *testing.T) {
	// When automation features are disabled, cron/bg/workflow guards
	// should still block but suggest enabling the feature.
	auto := map[string]bool{}

	msg := blockLangoExec("lango cron list", auto)
	if msg == "" {
		t.Fatal("expected blocked message for disabled cron")
	}
	if !strings.Contains(msg, "Enable the") {
		t.Errorf("expected 'Enable the' suggestion, got: %s", msg)
	}

	// Non-automation guards (graph, memory, etc.) should always block
	// regardless of automation flags.
	msg = blockLangoExec("lango graph query", auto)
	if msg == "" {
		t.Fatal("expected blocked message for graph")
	}
	if strings.Contains(msg, "Enable the") {
		t.Errorf("graph guard should not suggest enabling a feature, got: %s", msg)
	}
}
