package learning

import (
	"encoding/json"
	"fmt"
	"strings"
)

// analysisResult is the expected structure from LLM analysis output.
type analysisResult struct {
	Type       string `json:"type"`                // fact, pattern, correction, preference
	Category   string `json:"category"`            // domain-specific category
	Content    string `json:"content"`             // the extracted knowledge
	Confidence string `json:"confidence"`          // low, medium, high
	Subject    string `json:"subject,omitempty"`   // optional graph subject
	Predicate  string `json:"predicate,omitempty"` // optional graph predicate
	Object     string `json:"object,omitempty"`    // optional graph object
}

// parseAnalysisResponse extracts structured results from an LLM JSON response.
// Handles code fences, single objects, and arrays.
func parseAnalysisResponse(raw string) ([]analysisResult, error) {
	cleaned := stripCodeFence(raw)
	cleaned = strings.TrimSpace(cleaned)

	// Try array first.
	var results []analysisResult
	if err := json.Unmarshal([]byte(cleaned), &results); err == nil {
		return results, nil
	}

	// Try single object.
	var single analysisResult
	if err := json.Unmarshal([]byte(cleaned), &single); err == nil {
		return []analysisResult{single}, nil
	}

	return nil, fmt.Errorf("parse analysis response: invalid JSON")
}

// mapKnowledgeCategory maps LLM analysis type to a valid ent knowledge category enum.
func mapKnowledgeCategory(analysisType string) string {
	switch analysisType {
	case "preference":
		return "preference"
	case "fact":
		return "fact"
	case "rule":
		return "rule"
	case "definition":
		return "definition"
	default:
		return "fact"
	}
}

// mapLearningCategory maps LLM analysis type to a valid ent learning category enum.
func mapLearningCategory(analysisType string) string {
	switch analysisType {
	case "correction":
		return "user_correction"
	case "pattern":
		return "general"
	case "tool_error":
		return "tool_error"
	case "provider_error":
		return "provider_error"
	case "timeout":
		return "timeout"
	case "permission":
		return "permission"
	default:
		return "general"
	}
}

// stripCodeFence removes markdown code fences from LLM output.
func stripCodeFence(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```json") {
		s = strings.TrimPrefix(s, "```json")
	} else if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
	}
	if strings.HasSuffix(s, "```") {
		s = strings.TrimSuffix(s, "```")
	}
	return strings.TrimSpace(s)
}
