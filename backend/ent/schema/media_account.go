package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type MediaAccount struct {
	ent.Schema
}

func (MediaAccount) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (MediaAccount) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(100).
			NotEmpty(),
		field.String("platform").
			MaxLen(50).
			NotEmpty(),
		field.String("account_id").
			MaxLen(100).
			NotEmpty(),
		field.String("avatar").
			MaxLen(500).
			Optional(),
		field.String("status").
			MaxLen(20).
			Default("active"),
		field.String("access_token").
			MaxLen(500).
			Optional(),
		field.String("refresh_token").
			MaxLen(500).
			Optional(),
		field.Time("token_expires_at").
			Optional(),
	}
}

func (MediaAccount) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("contents", MediaContent.Type),
		edge.To("fans_snapshots", MediaAccountFansSnapshot.Type),
	}
}
