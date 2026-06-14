package mixins

import (
	"context"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// SoftDeleteMixin implements soft-delete via deleted_at timestamp.
type SoftDeleteMixin struct {
	mixin.Schema
}

func (SoftDeleteMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("deleted_at").
			Nillable().
			Optional(),
	}
}

type softDeleteKey struct{}

// SkipSoftDelete returns a context that disables soft-delete filtering.
func SkipSoftDelete(parent context.Context) context.Context {
	return context.WithValue(parent, softDeleteKey{}, true)
}

// SkipSoftDeleteFrom checks whether soft-delete should be skipped.
func SkipSoftDeleteFrom(ctx context.Context) bool {
	v, _ := ctx.Value(softDeleteKey{}).(bool)
	return v
}

// Interceptors adds automatic deleted_at IS NULL filtering to all queries.
func (d SoftDeleteMixin) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.InterceptFunc(func(next ent.Querier) ent.Querier {
			return ent.QuerierFunc(func(ctx context.Context, query ent.Query) (ent.Value, error) {
				if SkipSoftDeleteFrom(ctx) {
					return next.Query(ctx, query)
				}
				// Try to add the soft-delete predicate to the query
				type whereP interface {
					WhereP(...func(*sql.Selector))
				}
				if wp, ok := query.(whereP); ok {
					wp.WhereP(sql.FieldIsNull(d.Fields()[0].Descriptor().Name))
				}
				return next.Query(ctx, query)
			})
		}),
	}
}

// Hooks intercepts DELETE operations and converts them to UPDATE SET deleted_at = NOW().
func (d SoftDeleteMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if m.Op() != ent.OpDelete && m.Op() != ent.OpDeleteOne {
					return next.Mutate(ctx, m)
				}
				if SkipSoftDeleteFrom(ctx) {
					return next.Mutate(ctx, m)
				}

				mx, ok := m.(interface {
					SetOp(ent.Op)
					SetDeletedAt(time.Time)
					WhereP(...func(*sql.Selector))
				})
				if !ok {
					return next.Mutate(ctx, m)
				}

				mx.WhereP(sql.FieldIsNull(d.Fields()[0].Descriptor().Name))
				mx.SetOp(ent.OpUpdate)
				mx.SetDeletedAt(time.Now())
				return next.Mutate(ctx, m)
			})
		},
	}
}
