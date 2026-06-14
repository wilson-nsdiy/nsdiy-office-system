package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type ApiToken struct {
	ent.Schema
}

func (ApiToken) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (ApiToken) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id").
			Positive(),
		field.String("name").
			MaxLen(100).
			NotEmpty(),
		field.String("token_hash").
			Unique().
			NotEmpty(),
		field.String("token_prefix").
			MaxLen(20).
			NotEmpty(),
		field.String("status").
			MaxLen(20).
			Default("active"),
		field.Time("expires_at").
			Optional(),
		field.Time("last_used_at").
			Optional(),
		field.Int("usage_count").
			Default(0),
	}
}

func (ApiToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("api_tokens").
			Field("user_id").
			Unique().
			Required(),
	}
}
