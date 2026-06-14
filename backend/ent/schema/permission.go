package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type Permission struct {
	ent.Schema
}

func (Permission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.Int("pid").
			Optional(),
		field.String("name").
			Unique().
			MaxLen(100).
			NotEmpty(),
		field.String("resource_type").
			MaxLen(50).
			NotEmpty(),
		field.String("resource_path").
			MaxLen(200).
			NotEmpty(),
		field.String("http_method").
			MaxLen(10).
			Optional(),
		field.String("description").
			MaxLen(500).
			Optional(),
		field.Bool("is_active").
			Default(true),
		field.Bool("is_builtin").
			Default(false),
	}
}

func (Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("role_perms", RolePerm.Type),
	}
}
