package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Secret holds the schema definition for the Secret entity.
// Secrets store encrypted values that can be retrieved by AI agents.
type Secret struct {
	ent.Schema
}

// Fields of the Secret.
func (Secret) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("name").
			Unique().
			NotEmpty().
			Comment("Lookup key for the secret"),
		field.Bytes("encrypted_value").
			Comment("Encrypted secret data"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Int("access_count").
			Default(0).
			Comment("Number of times this secret has been accessed"),
	}
}

// Edges of the Secret.
func (Secret) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("key", Key.Type).
			Ref("secrets").
			Unique().
			Required().
			Comment("Encryption key used for this secret"),
	}
}

// Indexes of the Secret.
func (Secret) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
		index.Fields("created_at"),
	}
}
