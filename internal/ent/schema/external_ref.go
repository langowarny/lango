package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// ExternalRef holds the schema definition for the ExternalRef entity.
// ExternalRef stores references to external knowledge sources.
type ExternalRef struct {
	ent.Schema
}

// Fields of the ExternalRef.
func (ExternalRef) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("name").
			Unique().
			NotEmpty(),
		field.Enum("ref_type").
			Values("url", "file", "mcp"),
		field.String("location").
			NotEmpty(),
		field.Text("summary").
			Optional(),
		field.JSON("metadata", map[string]interface{}{}).
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the ExternalRef.
func (ExternalRef) Edges() []ent.Edge {
	return nil
}

// Indexes of the ExternalRef.
func (ExternalRef) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ref_type"),
	}
}
