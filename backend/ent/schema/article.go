package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type Article struct {
	ent.Schema
}

func (Article) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (Article) Fields() []ent.Field {
	return []ent.Field{
		field.String("article_no").
			Unique().
			MaxLen(50).
			NotEmpty(),
		field.String("title").
			MaxLen(200).
			NotEmpty(),
		field.String("content").
			Optional(),
		field.String("summary").
			MaxLen(1000).
			Optional(),
		field.String("status").
			MaxLen(20).
			Default("DRAFT"),
		field.Int("author_id").
			Positive(),
		field.String("cover_description").
			MaxLen(500).
			Optional(),
		field.String("cover_url").
			MaxLen(500).
			Optional(),
		field.Time("first_published_at").
			Optional(),
	}
}

func (Article) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("author", User.Type).
			Ref("articles").
			Field("author_id").
			Unique().
			Required(),
		edge.To("versions", ArticleVersion.Type),
	}
}
