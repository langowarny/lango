package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// WorkflowRun holds the schema definition for a workflow execution instance.
type WorkflowRun struct {
	ent.Schema
}

// Fields of the WorkflowRun.
func (WorkflowRun) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("workflow_name").
			NotEmpty().
			Comment("Name of the workflow being executed"),
		field.String("description").
			Optional(),
		field.Enum("status").
			Values("pending", "running", "completed", "failed", "cancelled").
			Default("pending"),
		field.Int("total_steps").
			Default(0),
		field.Int("completed_steps").
			Default(0),
		field.String("error_message").
			Optional(),
		field.Time("started_at").
			Default(time.Now).
			Immutable(),
		field.Time("completed_at").
			Optional().
			Nillable(),
	}
}

// Edges of the WorkflowRun.
func (WorkflowRun) Edges() []ent.Edge {
	return nil
}

// Indexes of the WorkflowRun.
func (WorkflowRun) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
		index.Fields("workflow_name"),
		index.Fields("started_at"),
	}
}
