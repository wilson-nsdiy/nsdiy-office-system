package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type Project struct {
	ent.Schema
}

func (Project) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(200).
			NotEmpty(),
		field.String("project_no").
			Unique().
			MaxLen(50).
			NotEmpty(),
		field.String("description").
			MaxLen(2000).
			Optional(),
		field.String("status").
			MaxLen(20).
			Default("TODO"),
		field.String("priority").
			MaxLen(20).
			Default("MEDIUM"),
		field.Time("expected_start_date").
			Optional(),
		field.Time("expected_end_date").
			Optional(),
		field.Time("start_date").
			Optional(),
		field.Time("end_date").
			Optional(),
		field.Int("owner_id").
			Positive(),
	}
}

func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("owned_projects").
			Field("owner_id").
			Unique().
			Required(),
		edge.To("members", ProjectMember.Type),
		edge.To("tasks", Task.Type),
		edge.To("operation_logs", OperationLog.Type),
	}
}
