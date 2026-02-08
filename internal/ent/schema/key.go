package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Key holds the schema definition for the Key entity.
// Keys are used for encryption/signing operations.
type Key struct {
	ent.Schema
}

// Fields of the Key.
func (Key) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("name").
			Unique().
			NotEmpty().
			Comment("Display name for the key"),
		field.String("remote_key_id").
			NotEmpty().
			Comment("Reference to key in external provider (companion) or 'local'"),
		field.Enum("type").
			Values("encryption", "signing").
			Default("encryption").
			Comment("Key purpose"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("last_used_at").
			Optional().
			Nillable().
			Comment("Last time this key was used for an operation"),
	}
}

// Edges of the Key.
func (Key) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("secrets", Secret.Type).
			Comment("Secrets encrypted with this key"),
	}
}

// Indexes of the Key.
func (Key) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
		index.Fields("type"),
	}
}
