package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	coremixin "goapp/internal/core/mixin"
)

// Accounting holds the schema definition for the Accounting entity.
type Accounting struct {
	ent.Schema
}

func (Accounting) Mixin() []ent.Mixin {
	return []ent.Mixin{
		coremixin.Base{},
	}
}

func (Accounting) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(255),

		field.String("description").
			Optional().
			Nillable(),
	}
}

func (Accounting) Edges() []ent.Edge {
	return nil
}
