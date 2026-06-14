package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type UploadFile struct {
	ent.Schema
}

func (UploadFile) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (UploadFile) Fields() []ent.Field {
	return []ent.Field{
		field.String("filename").
			Unique().
			MaxLen(200).
			NotEmpty(),
		field.String("original_filename").
			MaxLen(200).
			NotEmpty(),
		field.String("file_path").
			MaxLen(500).
			NotEmpty(),
		field.Int64("file_size").
			NonNegative(),
		field.String("mime_type").
			MaxLen(100).
			NotEmpty(),
		field.String("file_type").
			MaxLen(50).
			NotEmpty(),
		field.String("extension").
			MaxLen(20).
			NotEmpty(),
		field.Int("uploader_id").
			Positive(),
		field.String("purpose").
			MaxLen(100).
			Optional(),
		field.String("md5").
			MaxLen(32).
			Optional(),
		field.Int("reference_count").
			Default(1),
	}
}

func (UploadFile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("uploader", User.Type).
			Ref("uploaded_files").
			Field("uploader_id").
			Unique().
			Required(),
	}
}
