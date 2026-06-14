package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type Task struct {
	ent.Schema
}

func (Task) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.Int("project_id").
			Positive(),
		field.Int("parent_id").
			Optional().
			Nillable(),
		field.String("title").
			MaxLen(200).
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
		field.Int("assignee_id").
			Optional().
			Nillable(),
		field.Int("creator_id").
			Positive(),
		field.Time("planned_start_date").
			Optional().
			Nillable(),
		field.Time("planned_end_date").
			Optional().
			Nillable(),
		field.Time("actual_start_time").
			Optional().
			Nillable(),
		field.Time("actual_end_time").
			Optional().
			Nillable(),
		field.Float("estimated_hours").
			Optional().
			Nillable(),
	}
}

func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("tasks").
			Field("project_id").
			Unique().
			Required(),
		edge.From("assignee", User.Type).
			Ref("assigned_tasks").
			Field("assignee_id").
			Unique(),
		edge.From("creator", User.Type).
			Ref("created_tasks").
			Field("creator_id").
			Unique().
			Required(),
	}
}
