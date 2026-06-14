package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type MediaAccountFansSnapshot struct {
	ent.Schema
}

func (MediaAccountFansSnapshot) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (MediaAccountFansSnapshot) Fields() []ent.Field {
	return []ent.Field{
		field.Int("account_id").
			Positive(),
		field.Int("fans_count").
			NonNegative(),
		field.Time("snapshot_date").
			Default(time.Now),
	}
}

func (MediaAccountFansSnapshot) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", MediaAccount.Type).
			Ref("fans_snapshots").
			Field("account_id").
			Unique().
			Required(),
	}
}
