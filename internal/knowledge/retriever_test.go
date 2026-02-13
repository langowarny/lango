package knowledge

import (
	"context"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/langowarny/lango/internal/ent/enttest"
	_ "github.com/mattn/go-sqlite3"
)

func newTestRetriever(t *testing.T) (*ContextRetriever, *Store) {
	t.Helper()
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	t.Cleanup(func() { client.Close() })
	logger := zap.NewNop().Sugar()
	store := NewStore(client, logger, 20, 10, 5)
	retriever := NewContextRetriever(store, 5, logger)
	return retriever, store
}

func TestExtractKeywords(t *testing.T) {
	tests := []struct {
		give string
		want []string
	}{
		{
			give: "how to handle errors in Go",
			want: []string{"handle", "errors", "go"},
		},
		{
			give: "the a is are was",
			want: nil,
		},
		{
			give: "",
			want: nil,
		},
		{
			give: "hello, world! testing...",
			want: []string{"hello", "world", "testing"},
		},
		{
			give: "a b c x",
			want: nil,
		},
		{
			give: "Deploy PRODUCTION server",
			want: []string{"deploy", "production", "server"},
		},
		{
			give: "(important) [note] {value}",
			want: []string{"important", "note", "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			got := extractKeywords(tt.give)
			if tt.want == nil {
				if got != nil {
					t.Errorf("want nil, got %v", got)
				}
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("want %d keywords %v, got %d keywords %v", len(tt.want), tt.want, len(got), got)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("keyword[%d]: want %q, got %q", i, tt.want[i], got[i])
				}
			}
		})
	}
}

func TestContextRetriever_Retrieve(t *testing.T) {
	retriever, store := newTestRetriever(t)
	ctx := context.Background()

	// Seed data
	if err := store.SaveKnowledge(ctx, "s1", KnowledgeEntry{
		Key:      "deployment-rule",
		Category: "rule",
		Content:  "Always deploy with rollback plan",
	}); err != nil {
		t.Fatalf("SaveKnowledge: %v", err)
	}

	if err := store.SaveLearning(ctx, "s1", LearningEntry{
		Trigger:      "deployment failure",
		ErrorPattern: "deploy timeout",
		Fix:          "Increase deployment timeout to 5m",
		Category:     "timeout",
	}); err != nil {
		t.Fatalf("SaveLearning: %v", err)
	}

	if err := store.SaveSkill(ctx, SkillEntry{
		Name:        "deploy-canary",
		Description: "Canary deployment workflow",
		Type:        "composite",
		Definition:  map[string]interface{}{"steps": []interface{}{"build", "deploy"}},
	}); err != nil {
		t.Fatalf("SaveSkill: %v", err)
	}
	if err := store.ActivateSkill(ctx, "deploy-canary"); err != nil {
		t.Fatalf("ActivateSkill: %v", err)
	}

	if err := store.SaveExternalRef(ctx, "deploy-guide", "url", "https://example.com/deploy", "Deployment best practices guide"); err != nil {
		t.Fatalf("SaveExternalRef: %v", err)
	}

	t.Run("multi-layer retrieval", func(t *testing.T) {
		result, err := retriever.Retrieve(ctx, RetrievalRequest{
			Query:      "rollback",
			SessionKey: "s1",
			Layers: []ContextLayer{
				LayerUserKnowledge,
				LayerAgentLearnings,
				LayerSkillPatterns,
				LayerExternalKnowledge,
			},
		})
		if err != nil {
			t.Fatalf("Retrieve: %v", err)
		}
		if result.TotalItems == 0 {
			t.Fatal("expected at least one item across layers")
		}

		// Check knowledge layer
		if items, ok := result.Items[LayerUserKnowledge]; ok {
			found := false
			for _, item := range items {
				if item.Key == "deployment-rule" {
					found = true
				}
			}
			if !found {
				t.Error("expected deployment-rule in user knowledge layer")
			}
		}
	})

	t.Run("single layer", func(t *testing.T) {
		result, err := retriever.Retrieve(ctx, RetrievalRequest{
			Query:  "deploy canary",
			Layers: []ContextLayer{LayerSkillPatterns},
		})
		if err != nil {
			t.Fatalf("Retrieve: %v", err)
		}
		if _, ok := result.Items[LayerUserKnowledge]; ok {
			t.Error("should not have user knowledge when only skill layer requested")
		}
	})
}

func TestContextRetriever_Retrieve_DefaultLayers(t *testing.T) {
	retriever, store := newTestRetriever(t)
	ctx := context.Background()

	if err := store.SaveKnowledge(ctx, "s1", KnowledgeEntry{
		Key:      "default-layer-test",
		Category: "fact",
		Content:  "Testing default layer selection",
	}); err != nil {
		t.Fatalf("SaveKnowledge: %v", err)
	}

	result, err := retriever.Retrieve(ctx, RetrievalRequest{
		Query: "testing default layer",
	})
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}

	// With nil Layers, all 4 default layers should be checked
	// We just verify no error and the result is not nil
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Items == nil {
		t.Fatal("expected non-nil Items map")
	}
}

func TestContextRetriever_Retrieve_StopWordsOnly(t *testing.T) {
	retriever, _ := newTestRetriever(t)
	ctx := context.Background()

	result, err := retriever.Retrieve(ctx, RetrievalRequest{
		Query: "the a is",
	})
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if result.TotalItems != 0 {
		t.Errorf("want 0 total items for stop-words-only query, got %d", result.TotalItems)
	}
}

func TestAssemblePrompt(t *testing.T) {
	retriever, _ := newTestRetriever(t)
	basePrompt := "You are a helpful assistant."

	t.Run("nil result returns base", func(t *testing.T) {
		got := retriever.AssemblePrompt(basePrompt, nil)
		if got != basePrompt {
			t.Errorf("want base prompt, got %q", got)
		}
	})

	t.Run("empty items returns base", func(t *testing.T) {
		result := &RetrievalResult{
			Items:      map[ContextLayer][]ContextItem{},
			TotalItems: 0,
		}
		got := retriever.AssemblePrompt(basePrompt, result)
		if got != basePrompt {
			t.Errorf("want base prompt, got %q", got)
		}
	})

	t.Run("4-layer result formats correctly", func(t *testing.T) {
		result := &RetrievalResult{
			Items: map[ContextLayer][]ContextItem{
				LayerUserKnowledge: {
					{Layer: LayerUserKnowledge, Key: "rule-1", Content: "Always test", Category: "rule"},
				},
				LayerAgentLearnings: {
					{Layer: LayerAgentLearnings, Key: "learning-1", Content: "When 'timeout' occurs: increase timeout"},
				},
				LayerSkillPatterns: {
					{Layer: LayerSkillPatterns, Key: "deploy", Content: "Deploy service"},
				},
				LayerExternalKnowledge: {
					{Layer: LayerExternalKnowledge, Key: "docs", Content: "API reference", Source: "https://example.com"},
				},
			},
			TotalItems: 4,
		}
		got := retriever.AssemblePrompt(basePrompt, result)

		if !strings.HasPrefix(got, basePrompt) {
			t.Error("assembled prompt should start with base prompt")
		}
		if !strings.Contains(got, "## User Knowledge") {
			t.Error("expected User Knowledge section")
		}
		if !strings.Contains(got, "## Known Solutions") {
			t.Error("expected Known Solutions section")
		}
		if !strings.Contains(got, "## Available Skills") {
			t.Error("expected Available Skills section")
		}
		if !strings.Contains(got, "## External References") {
			t.Error("expected External References section")
		}
		if !strings.Contains(got, "rule-1") {
			t.Error("expected rule-1 in output")
		}
		if !strings.Contains(got, "deploy") {
			t.Error("expected deploy in output")
		}
		if !strings.Contains(got, "https://example.com") {
			t.Error("expected source URL in output")
		}
	})
}
