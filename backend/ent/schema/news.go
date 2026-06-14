package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type News struct {
	ent.Schema
}

func (News) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (News) Fields() []ent.Field {
	return []ent.Field{
		field.Int("group_id").
			Positive(),
		field.String("title").
			MaxLen(200).
			NotEmpty(),
		field.String("content").
			Optional(),
		field.Int("creator_id").
			Positive(),
	}
}

func (News) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("group", NewsGroup.Type).
			Ref("news").
			Field("group_id").
			Unique().
			Required(),
		edge.From("creator", User.Type).
			Ref("news").
			Field("creator_id").
			Unique().
			Required(),
	}
}
