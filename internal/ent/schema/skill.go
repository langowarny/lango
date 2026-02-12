package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Skill holds the schema definition for the Skill entity.
// Skill stores dynamic reusable skills (composite tool chains, scripts, templates).
type Skill struct {
	ent.Schema
}

// Fields of the Skill.
func (Skill) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("name").
			Unique().
			NotEmpty(),
		field.Text("description").
			NotEmpty(),
		field.Enum("skill_type").
			Values("composite", "script", "template"),
		field.JSON("definition", map[string]interface{}{}).
			Comment("The actual skill definition (steps, script, or template)"),
		field.JSON("parameters", map[string]interface{}{}).
			Optional().
			Comment("Parameter schema for the skill"),
		field.Enum("status").
			Values("active", "draft", "disabled").
			Default("draft"),
		field.String("created_by").
			Optional(),
		field.Int("use_count").
			Default(0),
		field.Int("success_count").
			Default(0),
		field.Time("last_used_at").
			Optional().
			Nillable(),
		field.Bool("requires_approval").
			Default(true),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Skill.
func (Skill) Edges() []ent.Edge {
	return nil
}

// Indexes of the Skill.
func (Skill) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
		index.Fields("skill_type"),
	}
}
