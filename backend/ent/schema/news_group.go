package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type NewsGroup struct {
	ent.Schema
}

func (NewsGroup) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (NewsGroup) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			MaxLen(100).
			NotEmpty(),
		field.String("description").
			MaxLen(500).
			Optional(),
		field.Int("sort_order").
			Default(0),
	}
}

func (NewsGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("news", News.Type),
	}
}
