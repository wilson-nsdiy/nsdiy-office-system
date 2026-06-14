package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type MediaContent struct {
	ent.Schema
}

func (MediaContent) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (MediaContent) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").
			MaxLen(200).
			NotEmpty(),
		field.String("content").
			Optional(),
		field.String("cover_image").
			MaxLen(500).
			Optional(),
		field.String("platform").
			MaxLen(50).
			NotEmpty(),
		field.Int("account_id").
			Optional(),
		field.String("status").
			MaxLen(20).
			Default("draft"),
		field.Int("views").
			Default(0),
		field.Int("likes").
			Default(0),
		field.Int("comments").
			Default(0),
		field.Int("shares").
			Default(0),
		field.Time("publish_time").
			Optional(),
	}
}

func (MediaContent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", MediaAccount.Type).
			Ref("contents").
			Field("account_id").
			Unique(),
		edge.To("versions", MediaContentVersion.Type),
	}
}
