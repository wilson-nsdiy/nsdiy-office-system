package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type ProjectMember struct {
	ent.Schema
}

func (ProjectMember) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (ProjectMember) Fields() []ent.Field {
	return []ent.Field{
		field.Int("project_id").
			Positive(),
		field.Int("user_id").
			Positive(),
		field.String("role").
			MaxLen(20).
			Default("MEMBER"),
		field.Time("joined_at").
			Default(time.Now),
	}
}

func (ProjectMember) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("members").
			Field("project_id").
			Unique().
			Required(),
		edge.From("user", User.Type).
			Ref("project_memberships").
			Field("user_id").
			Unique().
			Required(),
	}
}

func (ProjectMember) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id", "user_id").
			Unique(),
	}
}
