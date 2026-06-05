package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	coremixin "goapp/internal/core/mixin"
)

// Test holds the schema definition for the Test entity.
type Test struct {
	ent.Schema
}

func (Test) Mixin() []ent.Mixin {
	return []ent.Mixin{
		coremixin.Base{},
	}
}

func (Test) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(255),

		field.String("description").
			Optional().
			Nillable(),
	}
}

func (Test) Edges() []ent.Edge {
	return nil
}
