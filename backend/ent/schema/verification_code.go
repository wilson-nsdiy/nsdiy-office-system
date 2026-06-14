package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"oa-nsdiy/backend/ent/schema/mixins"
)

type VerificationCode struct {
	ent.Schema
}

func (VerificationCode) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (VerificationCode) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			MaxLen(10).
			NotEmpty(),
		field.String("target").
			MaxLen(200).
			NotEmpty(),
		field.String("channel").
			MaxLen(50).
			NotEmpty(),
		field.String("scene").
			MaxLen(50).
			NotEmpty(),
		field.Bool("is_used").
			Default(false),
		field.Int("attempts").
			Default(0),
		field.Int("max_attempts").
			Default(5),
		field.Time("expires_at").
			Default(time.Now().Add(30 * time.Minute)),
		field.Time("used_at").
			Optional(),
		field.String("ip_address").
			MaxLen(50).
			Optional(),
		field.String("user_agent").
			MaxLen(500).
			Optional(),
	}
}
