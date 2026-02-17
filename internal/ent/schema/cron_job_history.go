package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// CronJobHistory holds the schema definition for a cron job execution record.
type CronJobHistory struct {
	ent.Schema
}

// Fields of the CronJobHistory.
func (CronJobHistory) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("job_id", uuid.UUID{}).
			Comment("Reference to the CronJob that was executed"),
		field.String("job_name").
			NotEmpty().
			Comment("Snapshot of job name at execution time"),
		field.Enum("status").
			Values("running", "completed", "failed").
			Default("running"),
		field.Text("prompt").
			Comment("Prompt that was executed"),
		field.Text("result").
			Optional().
			Comment("Agent response"),
		field.String("error_message").
			Optional().
			Comment("Error details if execution failed"),
		field.Int("tokens_used").
			Default(0),
		field.Time("started_at").
			Default(time.Now).
			Immutable(),
		field.Time("completed_at").
			Optional().
			Nillable(),
	}
}

// Edges of the CronJobHistory.
func (CronJobHistory) Edges() []ent.Edge {
	return nil
}

// Indexes of the CronJobHistory.
func (CronJobHistory) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("job_id"),
		index.Fields("status"),
		index.Fields("started_at"),
	}
}
