package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type MediaContentVersion struct {
	ent.Schema
}

func (MediaContentVersion) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (MediaContentVersion) Fields() []ent.Field {
	return []ent.Field{
		field.Int("content_id").
			Positive(),
		field.Int("version_no").
			Positive(),
		field.String("title").
			MaxLen(200).
			NotEmpty(),
		field.String("content").
			Optional(),
		field.String("cover_image").
			MaxLen(500).
			Optional(),
		field.String("status").
			MaxLen(20).
			NotEmpty(),
		field.Int("editor_id").
			Optional(),
		field.String("edit_reason").
			MaxLen(500).
			Optional(),
	}
}

func (MediaContentVersion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("media_content", MediaContent.Type).
			Ref("versions").
			Field("content_id").
			Unique().
			Required(),
		edge.From("editor", User.Type).
			Ref("edited_media_versions").
			Field("editor_id").
			Unique(),
	}
}

func (MediaContentVersion) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("content_id", "version_no").
			Unique(),
	}
}
