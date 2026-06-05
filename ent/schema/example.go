package schema

import "entgo.io/ent"

// Example holds the schema definition for the Example entity.
type Example struct {
	ent.Schema
}

// Fields of the Example.
func (Example) Fields() []ent.Field {
	return nil
}

// Edges of the Example.
func (Example) Edges() []ent.Edge {
	return nil
}
