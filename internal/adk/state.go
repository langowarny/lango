package adk

import (
	"encoding/json"
	"iter"
	"time"

	"github.com/google/uuid"
	"google.golang.org/adk/session"
	"google.golang.org/genai"

	internal "github.com/langowarny/lango/internal/session"
	"google.golang.org/adk/model"
)

// SessionAdapter adapts internal.Session to adk.Session
type SessionAdapter struct {
	sess  *internal.Session
	store internal.Store
}

func NewSessionAdapter(s *internal.Session, store internal.Store) *SessionAdapter {
	return &SessionAdapter{sess: s, store: store}
}

func (s *SessionAdapter) ID() string      { return s.sess.Key }
func (s *SessionAdapter) AppName() string { return "lango" }
func (s *SessionAdapter) UserID() string  { return "user" } // Default header

func (s *SessionAdapter) State() session.State {
	return &StateAdapter{sess: s.sess, store: s.store}
}

func (s *SessionAdapter) Events() session.Events {
	return &EventsAdapter{history: s.sess.History}
}
func (s *SessionAdapter) LastUpdateTime() time.Time { return s.sess.UpdatedAt }

// StateAdapter adapts internal.Session.Metadata to adk.State
type StateAdapter struct {
	sess  *internal.Session
	store internal.Store
}

func (s *StateAdapter) Get(key string) (any, error) {
	valStr, ok := s.sess.Metadata[key]
	if !ok {
		return nil, session.ErrStateKeyNotExist
	}
	var val any
	if err := json.Unmarshal([]byte(valStr), &val); err == nil {
		return val, nil
	}
	return valStr, nil
}

func (s *StateAdapter) Set(key string, val any) error {
	var valStr string
	if sStr, ok := val.(string); ok {
		valStr = sStr
	} else {
		b, err := json.Marshal(val)
		if err != nil {
			return err
		}
		valStr = string(b)
	}

	if s.sess.Metadata == nil {
		s.sess.Metadata = make(map[string]string)
	}
	s.sess.Metadata[key] = valStr

	return s.store.Update(s.sess)
}

func (s *StateAdapter) All() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for k, vStr := range s.sess.Metadata {
			var val any
			if err := json.Unmarshal([]byte(vStr), &val); err != nil {
				val = vStr
			}
			if !yield(k, val) {
				return
			}
		}
	}
}

// EventsAdapter adapts internal history to adk events
type EventsAdapter struct {
	history []internal.Message
}

func (e *EventsAdapter) All() iter.Seq[*session.Event] {
	return func(yield func(*session.Event) bool) {
		// Truncate history to last 100 messages
		// This prevents context window overflow and excessive processing
		msgs := e.history
		if len(msgs) > 100 {
			msgs = msgs[len(msgs)-100:]
		}

		for _, msg := range msgs {
			// Map Author
			author := "model" // Default for assistant
			if msg.Role == "user" {
				author = "user"
			} else if msg.Role == "assistant" {
				author = "lango-agent" // ADK agent name
			} else if msg.Role == "function" || msg.Role == "tool" {
				author = "tool"
			}

			var parts []*genai.Part
			if msg.Content != "" {
				parts = append(parts, &genai.Part{Text: msg.Content})
			}

			evt := &session.Event{
				ID:        uuid.NewString(), // Generate on fly as we don't store event IDs
				Timestamp: msg.Timestamp,
				Author:    author,
				LLMResponse: model.LLMResponse{
					Content: &genai.Content{
						Role:  msg.Role,
						Parts: parts,
					},
				},
			}
			// Map tool calls if present
			if len(msg.ToolCalls) > 0 {
				for _, tc := range msg.ToolCalls {
					args := make(map[string]any)
					// best effort json unmarshal
					_ = json.Unmarshal([]byte(tc.Input), &args)

					evt.LLMResponse.Content.Parts = append(evt.LLMResponse.Content.Parts, &genai.Part{
						FunctionCall: &genai.FunctionCall{
							Name: tc.Name,
							Args: args,
						},
					})
				}
			}
			if !yield(evt) {
				return
			}
		}
	}
}

func (e *EventsAdapter) Len() int {
	lenMsgs := len(e.history)
	if lenMsgs > 100 {
		return 100
	}
	return lenMsgs
}

func (e *EventsAdapter) At(i int) *session.Event {
	// Reusing logic from All is inefficient but simple for now
	// Ideally we refactor message conversion
	var found *session.Event
	count := 0
	e.All()(func(evt *session.Event) bool {
		if count == i {
			found = evt
			return false
		}
		count++
		return true
	})
	return found
}
