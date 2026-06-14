package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type ArticleVersion struct {
	ent.Schema
}

func (ArticleVersion) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (ArticleVersion) Fields() []ent.Field {
	return []ent.Field{
		field.Int("article_id").
			Positive(),
		field.Int("version_no").
			Positive(),
		field.String("title").
			MaxLen(200).
			NotEmpty(),
		field.String("content").
			Optional(),
		field.String("cover_description").
			MaxLen(500).
			Optional(),
		field.String("summary").
			MaxLen(1000).
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

func (ArticleVersion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("article", Article.Type).
			Ref("versions").
			Field("article_id").
			Unique().
			Required(),
		edge.From("editor", User.Type).
			Ref("edited_article_versions").
			Field("editor_id").
			Unique(),
	}
}

func (ArticleVersion) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("article_id", "version_no").
			Unique(),
	}
}
