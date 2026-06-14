package repository

import (
	"context"

	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/ent/news"
	"oa-nsdiy/backend/internal/domain"

	"entgo.io/ent/dialect/sql"
)

type NewsRepository struct {
	client *ent.Client
}

func NewNewsRepository(client *ent.Client) *NewsRepository {
	return &NewsRepository{client: client}
}

// Type aliases for backward compatibility
type News = domain.News
type NewsDetail = domain.NewsDetail

func toNews(e *ent.News) *News {
	return &News{
		ID:        e.ID,
		GroupID:   e.GroupID,
		Title:     e.Title,
		Content:   stringPtr(e.Content),
		CreatorID: e.CreatorID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func toNewsDetail(e *ent.News) *NewsDetail {
	d := &NewsDetail{
		News: *toNews(e),
	}
	if e.Edges.Group != nil {
		d.GroupName = e.Edges.Group.Name
	}
	if e.Edges.Creator != nil {
		d.CreatorNickname = stringPtr(e.Edges.Creator.Nickname)
	}
	return d
}

func toNewsDetails(es []*ent.News) []*NewsDetail {
	result := make([]*NewsDetail, len(es))
	for i, e := range es {
		result[i] = toNewsDetail(e)
	}
	return result
}

func (r *NewsRepository) GetByID(ctx context.Context, id int) (*News, error) {
	e, err := r.client.News.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toNews(e), nil
}

func (r *NewsRepository) GetDetailByID(ctx context.Context, id int) (*NewsDetail, error) {
	e, err := r.client.News.Query().
		Where(news.IDEQ(id)).
		WithGroup().
		WithCreator().
		First(ctx)
	if err != nil {
		return nil, err
	}
	return toNewsDetail(e), nil
}

func (r *NewsRepository) Create(ctx context.Context, n *News) error {
	e, err := r.client.News.Create().
		SetGroupID(n.GroupID).
		SetTitle(n.Title).
		SetNillableContent(n.Content).
		SetCreatorID(n.CreatorID).
		Save(ctx)
	if err != nil {
		return err
	}
	n.ID = e.ID
	return nil
}

func (r *NewsRepository) Update(ctx context.Context, n *News) error {
	_, err := r.client.News.UpdateOneID(n.ID).
		SetGroupID(n.GroupID).
		SetTitle(n.Title).
		SetNillableContent(n.Content).
		Save(ctx)
	return err
}

func (r *NewsRepository) Delete(ctx context.Context, id int) error {
	return r.client.News.DeleteOneID(id).Exec(ctx)
}

func (r *NewsRepository) List(ctx context.Context, groupID *int, keyword string, page, pageSize int) ([]*NewsDetail, int64, error) {
	q := r.client.News.Query()

	if groupID != nil {
		q.Where(news.GroupIDEQ(*groupID))
	}
	if keyword != "" {
		q.Where(news.TitleContains(keyword))
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	es, err := q.
		Order(news.ByID(sql.OrderDesc())).
		Limit(pageSize).
		Offset(offset).
		WithGroup().
		WithCreator().
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return toNewsDetails(es), int64(total), nil
}
