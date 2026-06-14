package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type Role struct {
	ent.Schema
}

func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			MaxLen(100).
			NotEmpty(),
		field.String("code").
			Unique().
			MaxLen(50).
			NotEmpty(),
		field.String("description").
			MaxLen(500).
			Optional(),
		field.Bool("is_active").
			Default(true),
		field.Bool("is_default").
			Default(false),
		field.Bool("is_builtin").
			Default(false),
		field.String("role_type").
			MaxLen(50).
			Optional(),
	}
}

func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
		edge.To("role_perms", RolePerm.Type),
	}
}
