package app

import "testing"

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
