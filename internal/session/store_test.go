package session

import (
	"os"
	"testing"
)

func TestSQLiteStore(t *testing.T) {
	// Create temp database
	tmpFile, err := os.CreateTemp("", "sessions_test_*.db")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	store, err := NewSQLiteStore(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	// Test Create
	session := &Session{
		Key:         "test-session-1",
		AgentID:     "agent-1",
		ChannelType: "telegram",
		ChannelID:   "12345",
		History:     []Message{},
		Metadata:    map[string]string{"user": "test"},
		Model:       "claude-sonnet-4-20250514",
	}

	if err := store.Create(session); err != nil {
		t.Errorf("failed to create session: %v", err)
	}

	// Test Get
	retrieved, err := store.Get("test-session-1")
	if err != nil {
		t.Errorf("failed to get session: %v", err)
	}
	if retrieved.Key != session.Key {
		t.Errorf("expected key %s, got %s", session.Key, retrieved.Key)
	}
	if retrieved.Model != session.Model {
		t.Errorf("expected model %s, got %s", session.Model, retrieved.Model)
	}

	// Test Get non-existent
	_, err = store.Get("non-existent")
	if err == nil {
		t.Error("expected error for non-existent session")
	}

	// Test AppendMessage
	msg := Message{
		Role:    "user",
		Content: "Hello",
	}
	if err := store.AppendMessage("test-session-1", msg); err != nil {
		t.Errorf("failed to append message: %v", err)
	}

	retrieved, _ = store.Get("test-session-1")
	if len(retrieved.History) != 1 {
		t.Errorf("expected 1 message, got %d", len(retrieved.History))
	}
	if retrieved.History[0].Content != "Hello" {
		t.Errorf("expected content 'Hello', got %s", retrieved.History[0].Content)
	}

	// Test Update
	session.Model = "gpt-4"
	if err := store.Update(session); err != nil {
		t.Errorf("failed to update session: %v", err)
	}

	// Test Delete
	if err := store.Delete("test-session-1"); err != nil {
		t.Errorf("failed to delete session: %v", err)
	}

	_, err = store.Get("test-session-1")
	if err == nil {
		t.Error("expected error after deletion")
	}
}
