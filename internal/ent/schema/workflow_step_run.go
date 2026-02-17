package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// WorkflowStepRun holds the schema definition for a single step execution within a workflow run.
type WorkflowStepRun struct {
	ent.Schema
}

// Fields of the WorkflowStepRun.
func (WorkflowStepRun) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("run_id", uuid.UUID{}).
			Comment("Reference to the parent WorkflowRun"),
		field.String("step_id").
			NotEmpty().
			Comment("Step identifier within the workflow YAML"),
		field.String("agent").
			Optional().
			Comment("Sub-agent that executed this step"),
		field.Text("prompt").
			Comment("Rendered prompt (after template substitution)"),
		field.Enum("status").
			Values("pending", "running", "completed", "failed", "skipped").
			Default("pending"),
		field.Text("result").
			Optional().
			Comment("Step output/result"),
		field.String("error_message").
			Optional(),
		field.Time("started_at").
			Optional().
			Nillable(),
		field.Time("completed_at").
			Optional().
			Nillable(),
	}
}

// Edges of the WorkflowStepRun.
func (WorkflowStepRun) Edges() []ent.Edge {
	return nil
}

// Indexes of the WorkflowStepRun.
func (WorkflowStepRun) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("run_id"),
		index.Fields("status"),
		index.Fields("run_id", "step_id"),
	}
}
