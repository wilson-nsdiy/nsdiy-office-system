package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type OperationLog struct {
	ent.Schema
}

func (OperationLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (OperationLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id").
			Optional(),
		field.String("username").
			MaxLen(100).
			Optional(),
		field.String("operation_type").
			MaxLen(50).
			NotEmpty(),
		field.String("module").
			MaxLen(50).
			NotEmpty(),
		field.String("action").
			MaxLen(50).
			NotEmpty(),
		field.String("resource_type").
			MaxLen(50).
			Optional(),
		field.Int("resource_id").
			Optional(),
		field.String("resource_name").
			MaxLen(200).
			Optional(),
		field.Int("project_id").
			Optional(),
		field.String("detail").
			MaxLen(2000).
			Optional(),
		field.String("ip_address").
			MaxLen(50).
			Optional(),
		field.String("user_agent").
			MaxLen(500).
			Optional(),
		field.String("status").
			MaxLen(20).
			NotEmpty(),
		field.String("error_message").
			MaxLen(2000).
			Optional(),
	}
}

func (OperationLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("operation_logs").
			Field("user_id").
			Unique(),
		edge.From("project", Project.Type).
			Ref("operation_logs").
			Field("project_id").
			Unique(),
	}
}
