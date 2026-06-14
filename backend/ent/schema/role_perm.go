package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type RolePerm struct {
	ent.Schema
}

func (RolePerm) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (RolePerm) Fields() []ent.Field {
	return []ent.Field{
		field.Int("role_id").
			Positive(),
		field.Int("permission_id").
			Positive(),
	}
}

func (RolePerm) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", Role.Type).
			Ref("role_perms").
			Field("role_id").
			Unique().
			Required(),
		edge.From("permission", Permission.Type).
			Ref("role_perms").
			Field("permission_id").
			Unique().
			Required(),
	}
}
