package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").
			Unique().
			MaxLen(100).
			NotEmpty(),
		field.String("email").
			Unique().
			MaxLen(200).
			NotEmpty(),
		field.String("nickname").
			MaxLen(100).
			Optional(),
		field.String("salt").
			NotEmpty(),
		field.String("hashed_password").
			NotEmpty(),
		field.Int("role_id").
			Optional(),
		field.String("user_type").
			Default("HUMAN").
			MaxLen(20),
		field.Bool("is_active").
			Default(true),
		field.Int("token_version").
			Default(1),
		field.String("verification_code").
			MaxLen(10).
			Optional(),
		field.Time("verification_code_expires_at").
			Optional(),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", Role.Type).
			Ref("users").
			Field("role_id").
			Unique(),
		edge.To("api_tokens", ApiToken.Type),
		edge.To("articles", Article.Type),
		edge.To("news", News.Type),
		edge.To("owned_projects", Project.Type),
		edge.To("project_memberships", ProjectMember.Type),
		edge.To("assigned_tasks", Task.Type),
		edge.To("created_tasks", Task.Type),
		edge.To("edited_article_versions", ArticleVersion.Type),
		edge.To("uploaded_files", UploadFile.Type),
		edge.To("edited_media_versions", MediaContentVersion.Type),
		edge.To("operation_logs", OperationLog.Type),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username", "email").
			Unique(),
	}
}
