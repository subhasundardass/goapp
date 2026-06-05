package schema

import (
	coremixin "goapp/internal/core/mixin"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Ledger holds the schema definition for the Ledger entity.
type Ledger struct {
	ent.Schema
}

func (Ledger) Mixin() []ent.Mixin {
	return []ent.Mixin{
		coremixin.Base{},
	}
}

func (Ledger) Fields() []ent.Field {
	return []ent.Field{

		field.Int("group_id"),

		field.String("code").
			NotEmpty().
			Unique(),

		field.String("name").
			NotEmpty().
			MaxLen(255),

		field.Float("balance").
			Default(0).
			Comment("Current balance"),
	}
}

func (Ledger) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code").
			Unique(),

		index.Fields("name"),

		index.Fields("group_id"),

		index.Fields("group_id", "name"),
	}
}

func (Ledger) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("group", Acct_Group.Type).
			Field("group_id").
			Required().
			Unique(),

		edge.To("party", PartyMaster.Type).
			Unique(),
	}
}
