package repository

import (
	"context"
	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/ent/news"
	"oa-nsdiy/backend/ent/newsgroup"
	"oa-nsdiy/backend/internal/domain"
)

type NewsGroupRepository struct {
	client *ent.Client
}

func NewNewsGroupRepository(client *ent.Client) *NewsGroupRepository {
	return &NewsGroupRepository{client: client}
}

// Type alias for backward compatibility
type NewsGroup = domain.NewsGroup

func (r *NewsGroupRepository) GetByID(ctx context.Context, id int) (*NewsGroup, error) {
	e, err := r.client.NewsGroup.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toNewsGroupEntity(e), nil
}

func (r *NewsGroupRepository) GetByName(ctx context.Context, name string) (*NewsGroup, error) {
	e, err := r.client.NewsGroup.Query().Where(newsgroup.NameEQ(name)).First(ctx)
	if err != nil {
		return nil, err
	}
	return toNewsGroupEntity(e), nil
}

func (r *NewsGroupRepository) Create(ctx context.Context, group *NewsGroup) error {
	e, err := r.client.NewsGroup.Create().
		SetName(group.Name).
		SetNillableDescription(group.Description).
		SetSortOrder(group.SortOrder).
		Save(ctx)
	if err != nil {
		return err
	}
	group.ID = e.ID
	return nil
}

func (r *NewsGroupRepository) Update(ctx context.Context, group *NewsGroup) error {
	_, err := r.client.NewsGroup.UpdateOneID(group.ID).
		SetName(group.Name).
		SetNillableDescription(group.Description).
		SetSortOrder(group.SortOrder).
		Save(ctx)
	return err
}

func (r *NewsGroupRepository) Delete(ctx context.Context, id int) error {
	return r.client.NewsGroup.DeleteOneID(id).Exec(ctx)
}

func (r *NewsGroupRepository) ListAll(ctx context.Context) ([]*NewsGroup, error) {
	es, err := r.client.NewsGroup.Query().
		Order(ent.Asc(newsgroup.FieldSortOrder), ent.Asc(newsgroup.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toNewsGroupEntities(es), nil
}

func (r *NewsGroupRepository) GetNewsCount(ctx context.Context, groupID int) (int, error) {
	count, err := r.client.News.Query().Where(news.GroupIDEQ(groupID)).Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func toNewsGroupEntity(e *ent.NewsGroup) *NewsGroup {
	if e == nil {
		return nil
	}
	return &NewsGroup{
		ID:          e.ID,
		Name:        e.Name,
		Description: stringPtr(e.Description),
		SortOrder:   e.SortOrder,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toNewsGroupEntities(es []*ent.NewsGroup) []*NewsGroup {
	result := make([]*NewsGroup, len(es))
	for i, e := range es {
		result[i] = toNewsGroupEntity(e)
	}
	return result
}
