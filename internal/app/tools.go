package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/langowarny/lango/internal/agent"
	"github.com/langowarny/lango/internal/knowledge"
	"github.com/langowarny/lango/internal/learning"
	"github.com/langowarny/lango/internal/skill"
	"github.com/langowarny/lango/internal/supervisor"
	"github.com/langowarny/lango/internal/tools/filesystem"
)

// buildTools creates the set of tools available to the agent.
// MVP includes only exec and filesystem tools.
func buildTools(sv *supervisor.Supervisor, fsCfg filesystem.Config) []*agent.Tool {
	var tools []*agent.Tool

	// Exec tools (delegated to Supervisor for security isolation)
	tools = append(tools, buildExecTools(sv)...)

	// Filesystem tools
	fsTool := filesystem.New(fsCfg)
	tools = append(tools, buildFilesystemTools(fsTool)...)

	return tools
}

func buildExecTools(sv *supervisor.Supervisor) []*agent.Tool {
	return []*agent.Tool{
		{
			Name:        "exec",
			Description: "Execute shell commands",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]interface{}{
						"type":        "string",
						"description": "The shell command to execute",
					},
				},
				"required": []string{"command"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				cmd, ok := params["command"].(string)
				if !ok {
					return nil, fmt.Errorf("missing command parameter")
				}
				return sv.ExecuteTool(ctx, cmd)
			},
		},
		{
			Name:        "exec_bg",
			Description: "Execute a shell command in the background",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]interface{}{
						"type":        "string",
						"description": "The shell command to execute",
					},
				},
				"required": []string{"command"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				cmd, ok := params["command"].(string)
				if !ok {
					return nil, fmt.Errorf("missing command parameter")
				}
				return sv.StartBackground(cmd)
			},
		},
		{
			Name:        "exec_status",
			Description: "Check the status of a background process",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "The background process ID returned by exec_bg",
					},
				},
				"required": []string{"id"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				id, ok := params["id"].(string)
				if !ok {
					return nil, fmt.Errorf("missing id parameter")
				}
				return sv.GetBackgroundStatus(id)
			},
		},
		{
			Name:        "exec_stop",
			Description: "Stop a background process",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "The background process ID returned by exec_bg",
					},
				},
				"required": []string{"id"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				id, ok := params["id"].(string)
				if !ok {
					return nil, fmt.Errorf("missing id parameter")
				}
				return nil, sv.StopBackground(id)
			},
		},
	}
}

func buildFilesystemTools(fsTool *filesystem.Tool) []*agent.Tool {
	return []*agent.Tool{
		{
			Name:        "fs_read",
			Description: "Read a file",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{"type": "string", "description": "The file path to read"},
				},
				"required": []string{"path"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				path, ok := params["path"].(string)
				if !ok {
					return nil, fmt.Errorf("missing path parameter")
				}
				return fsTool.Read(path)
			},
		},
		{
			Name:        "fs_list",
			Description: "List contents of a directory",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{"type": "string", "description": "The directory path to list"},
				},
				"required": []string{"path"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				path, _ := params["path"].(string)
				if path == "" {
					path = "."
				}
				return fsTool.ListDir(path)
			},
		},
		{
			Name:        "fs_write",
			Description: "Write content to a file",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path":    map[string]interface{}{"type": "string", "description": "The file path to write to"},
					"content": map[string]interface{}{"type": "string", "description": "The content to write"},
				},
				"required": []string{"path", "content"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				path, _ := params["path"].(string)
				content, _ := params["content"].(string)
				if path == "" {
					return nil, fmt.Errorf("missing path parameter")
				}
				return nil, fsTool.Write(path, content)
			},
		},
		{
			Name:        "fs_edit",
			Description: "Edit a file by replacing a line range",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path":      map[string]interface{}{"type": "string", "description": "The file path to edit"},
					"startLine": map[string]interface{}{"type": "integer", "description": "The starting line number (1-indexed)"},
					"endLine":   map[string]interface{}{"type": "integer", "description": "The ending line number (inclusive)"},
					"content":   map[string]interface{}{"type": "string", "description": "The new content for the specified range"},
				},
				"required": []string{"path", "startLine", "endLine", "content"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				path, _ := params["path"].(string)
				content, _ := params["content"].(string)
				if path == "" {
					return nil, fmt.Errorf("missing path parameter")
				}

				var startLine, endLine int
				if sl, ok := params["startLine"].(float64); ok {
					startLine = int(sl)
				} else if sl, ok := params["startLine"].(int); ok {
					startLine = sl
				}
				if el, ok := params["endLine"].(float64); ok {
					endLine = int(el)
				} else if el, ok := params["endLine"].(int); ok {
					endLine = el
				}

				return nil, fsTool.Edit(path, startLine, endLine, content)
			},
		},
		{
			Name:        "fs_mkdir",
			Description: "Create a directory",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{"type": "string", "description": "The directory path to create"},
				},
				"required": []string{"path"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				path, _ := params["path"].(string)
				if path == "" {
					return nil, fmt.Errorf("missing path parameter")
				}
				return nil, fsTool.Mkdir(path)
			},
		},
		{
			Name:        "fs_delete",
			Description: "Delete a file or directory",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{"type": "string", "description": "The path to delete"},
				},
				"required": []string{"path"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				path, _ := params["path"].(string)
				if path == "" {
					return nil, fmt.Errorf("missing path parameter")
				}
				return nil, fsTool.Delete(path)
			},
		},
	}
}

// buildMetaTools creates knowledge/learning/skill meta-tools for the agent.
func buildMetaTools(store *knowledge.Store, engine *learning.Engine, registry *skill.Registry, autoApprove bool) []*agent.Tool {
	return []*agent.Tool{
		{
			Name:        "save_knowledge",
			Description: "Save a piece of knowledge (user rule, definition, preference, or fact) for future reference",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"key":      map[string]interface{}{"type": "string", "description": "Unique key for this knowledge entry"},
					"category": map[string]interface{}{"type": "string", "description": "Category: rule, definition, preference, or fact", "enum": []string{"rule", "definition", "preference", "fact"}},
					"content":  map[string]interface{}{"type": "string", "description": "The knowledge content to save"},
					"tags":     map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}, "description": "Optional tags for categorization"},
					"source":   map[string]interface{}{"type": "string", "description": "Where this knowledge came from"},
				},
				"required": []string{"key", "category", "content"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				key, _ := params["key"].(string)
				category, _ := params["category"].(string)
				content, _ := params["content"].(string)
				source, _ := params["source"].(string)

				if key == "" || category == "" || content == "" {
					return nil, fmt.Errorf("key, category, and content are required")
				}

				var tags []string
				if rawTags, ok := params["tags"].([]interface{}); ok {
					for _, t := range rawTags {
						if s, ok := t.(string); ok {
							tags = append(tags, s)
						}
					}
				}

				entry := knowledge.KnowledgeEntry{
					Key:      key,
					Category: category,
					Content:  content,
					Tags:     tags,
					Source:   source,
				}

				if err := store.SaveKnowledge(ctx, "", entry); err != nil {
					return nil, fmt.Errorf("save knowledge: %w", err)
				}

				_ = store.SaveAuditLog(ctx, knowledge.AuditEntry{
					Action: "knowledge_save",
					Actor:  "agent",
					Target: key,
				})

				return map[string]interface{}{
					"status":  "saved",
					"key":     key,
					"message": fmt.Sprintf("Knowledge '%s' saved successfully", key),
				}, nil
			},
		},
		{
			Name:        "search_knowledge",
			Description: "Search stored knowledge entries by query and optional category",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query":    map[string]interface{}{"type": "string", "description": "Search query"},
					"category": map[string]interface{}{"type": "string", "description": "Optional category filter: rule, definition, preference, or fact"},
				},
				"required": []string{"query"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				query, _ := params["query"].(string)
				category, _ := params["category"].(string)

				entries, err := store.SearchKnowledge(ctx, query, category, 10)
				if err != nil {
					return nil, fmt.Errorf("search knowledge: %w", err)
				}

				return map[string]interface{}{
					"results": entries,
					"count":   len(entries),
				}, nil
			},
		},
		{
			Name:        "save_learning",
			Description: "Save a diagnosed error pattern and its fix for future reference",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"trigger":       map[string]interface{}{"type": "string", "description": "What triggered this learning (e.g., tool name or action)"},
					"error_pattern": map[string]interface{}{"type": "string", "description": "The error pattern to match"},
					"diagnosis":     map[string]interface{}{"type": "string", "description": "Diagnosis of the error cause"},
					"fix":           map[string]interface{}{"type": "string", "description": "The fix or workaround"},
					"category":      map[string]interface{}{"type": "string", "description": "Category: tool_error, provider_error, user_correction, timeout, permission, general"},
				},
				"required": []string{"trigger", "fix"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				trigger, _ := params["trigger"].(string)
				errorPattern, _ := params["error_pattern"].(string)
				diagnosis, _ := params["diagnosis"].(string)
				fix, _ := params["fix"].(string)
				category, _ := params["category"].(string)

				if trigger == "" || fix == "" {
					return nil, fmt.Errorf("trigger and fix are required")
				}
				if category == "" {
					category = "general"
				}

				entry := knowledge.LearningEntry{
					Trigger:      trigger,
					ErrorPattern: errorPattern,
					Diagnosis:    diagnosis,
					Fix:          fix,
					Category:     category,
				}

				if err := store.SaveLearning(ctx, "", entry); err != nil {
					return nil, fmt.Errorf("save learning: %w", err)
				}

				_ = store.SaveAuditLog(ctx, knowledge.AuditEntry{
					Action: "learning_save",
					Actor:  "agent",
					Target: trigger,
				})

				return map[string]interface{}{
					"status":  "saved",
					"message": fmt.Sprintf("Learning for '%s' saved successfully", trigger),
				}, nil
			},
		},
		{
			Name:        "search_learnings",
			Description: "Search stored learnings by error pattern or trigger",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query":    map[string]interface{}{"type": "string", "description": "Search query (error message or trigger)"},
					"category": map[string]interface{}{"type": "string", "description": "Optional category filter"},
				},
				"required": []string{"query"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				query, _ := params["query"].(string)
				category, _ := params["category"].(string)

				entries, err := store.SearchLearnings(ctx, query, category, 10)
				if err != nil {
					return nil, fmt.Errorf("search learnings: %w", err)
				}

				return map[string]interface{}{
					"results": entries,
					"count":   len(entries),
				}, nil
			},
		},
		{
			Name:        "create_skill",
			Description: "Create a new reusable skill from a multi-step workflow, script, or template",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name":        map[string]interface{}{"type": "string", "description": "Unique name for the skill"},
					"description": map[string]interface{}{"type": "string", "description": "Description of what the skill does"},
					"type":        map[string]interface{}{"type": "string", "description": "Skill type: composite, script, or template", "enum": []string{"composite", "script", "template"}},
					"definition":  map[string]interface{}{"type": "string", "description": "JSON string of the skill definition"},
					"parameters":  map[string]interface{}{"type": "string", "description": "Optional JSON string of parameter schema"},
				},
				"required": []string{"name", "description", "type", "definition"},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				name, _ := params["name"].(string)
				description, _ := params["description"].(string)
				skillType, _ := params["type"].(string)
				definitionStr, _ := params["definition"].(string)

				if name == "" || description == "" || skillType == "" || definitionStr == "" {
					return nil, fmt.Errorf("name, description, type, and definition are required")
				}

				var definition map[string]interface{}
				if err := json.Unmarshal([]byte(definitionStr), &definition); err != nil {
					return nil, fmt.Errorf("parse definition JSON: %w", err)
				}

				var parameters map[string]interface{}
				if paramStr, ok := params["parameters"].(string); ok && paramStr != "" {
					if err := json.Unmarshal([]byte(paramStr), &parameters); err != nil {
						return nil, fmt.Errorf("parse parameters JSON: %w", err)
					}
				}

				entry := knowledge.SkillEntry{
					Name:             name,
					Description:      description,
					Type:             skillType,
					Definition:       definition,
					Parameters:       parameters,
					CreatedBy:        "agent",
					RequiresApproval: !autoApprove,
				}

				if err := registry.CreateSkill(ctx, entry); err != nil {
					return nil, fmt.Errorf("create skill: %w", err)
				}

				status := "draft"
				if autoApprove {
					if err := registry.ActivateSkill(ctx, name); err != nil {
						return nil, fmt.Errorf("activate skill: %w", err)
					}
					status = "active"
				}

				_ = store.SaveAuditLog(ctx, knowledge.AuditEntry{
					Action: "skill_create",
					Actor:  "agent",
					Target: name,
					Details: map[string]interface{}{
						"type":   skillType,
						"status": status,
					},
				})

				return map[string]interface{}{
					"status":  status,
					"name":    name,
					"message": fmt.Sprintf("Skill '%s' created with status '%s'", name, status),
				}, nil
			},
		},
		{
			Name:        "list_skills",
			Description: "List all active skills with usage statistics",
			Parameters: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
			Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				skills, err := store.ListActiveSkills(ctx)
				if err != nil {
					return nil, fmt.Errorf("list skills: %w", err)
				}

				return map[string]interface{}{
					"skills": skills,
					"count":  len(skills),
				}, nil
			},
		},
	}
}

// wrapWithLearning wraps a tool's handler to call the learning engine after each execution.
func wrapWithLearning(t *agent.Tool, engine *learning.Engine) *agent.Tool {
	original := t.Handler
	return &agent.Tool{
		Name:        t.Name,
		Description: t.Description,
		Parameters:  t.Parameters,
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			result, err := original(ctx, params)
			engine.OnToolResult(ctx, "", t.Name, params, result, err)
			return result, err
		},
	}
}
