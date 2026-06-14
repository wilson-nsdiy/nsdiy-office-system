package repository

import (
	"context"
	"fmt"

	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/ent/article"
	"oa-nsdiy/backend/ent/articleversion"
	"oa-nsdiy/backend/internal/domain"

	"entgo.io/ent/dialect/sql"
)

type ArticleRepository struct {
	client *ent.Client
}

func NewArticleRepository(client *ent.Client) *ArticleRepository {
	return &ArticleRepository{client: client}
}

// Type aliases for backward compatibility
type Article = domain.Article
type ArticleDetail = domain.ArticleDetail
type ArticleVersion = domain.ArticleVersion
type ArticleVersionDetail = domain.ArticleVersionDetail

func toArticle(e *ent.Article) *Article {
	return &Article{
		ID:               e.ID,
		ArticleNo:        e.ArticleNo,
		Title:            e.Title,
		Content:          stringPtr(e.Content),
		Summary:          stringPtr(e.Summary),
		Status:           e.Status,
		AuthorID:         e.AuthorID,
		CoverDescription: stringPtr(e.CoverDescription),
		CoverUrl:         stringPtr(e.CoverURL),
		FirstPublishedAt: timePtr(e.FirstPublishedAt),
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}

func toArticleDetail(e *ent.Article) *ArticleDetail {
	d := &ArticleDetail{
		Article: *toArticle(e),
	}
	if e.Edges.Author != nil {
		d.AuthorName = e.Edges.Author.Username
		d.AuthorNickname = stringPtr(e.Edges.Author.Nickname)
	}
	return d
}

func toArticleDetails(es []*ent.Article) []*ArticleDetail {
	result := make([]*ArticleDetail, len(es))
	for i, e := range es {
		result[i] = toArticleDetail(e)
	}
	return result
}

func (r *ArticleRepository) GetByID(ctx context.Context, id int) (*Article, error) {
	e, err := r.client.Article.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toArticle(e), nil
}

func (r *ArticleRepository) GetDetailByID(ctx context.Context, id int) (*ArticleDetail, error) {
	e, err := r.client.Article.Query().
		Where(article.IDEQ(id)).
		WithAuthor().
		First(ctx)
	if err != nil {
		return nil, err
	}
	return toArticleDetail(e), nil
}

func (r *ArticleRepository) GetByTitle(ctx context.Context, title string) (*Article, error) {
	e, err := r.client.Article.Query().
		Where(article.TitleEQ(title)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return toArticle(e), nil
}

func (r *ArticleRepository) GenerateArticleNo(ctx context.Context) (string, error) {
	count, err := r.client.Article.Query().Count(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("ART%04d", count+1), nil
}

func (r *ArticleRepository) Create(ctx context.Context, a *Article) error {
	e, err := r.client.Article.Create().
		SetArticleNo(a.ArticleNo).
		SetTitle(a.Title).
		SetNillableContent(a.Content).
		SetNillableSummary(a.Summary).
		SetStatus(a.Status).
		SetAuthorID(a.AuthorID).
		SetNillableCoverDescription(a.CoverDescription).
		SetNillableCoverURL(a.CoverUrl).
		SetNillableFirstPublishedAt(a.FirstPublishedAt).
		Save(ctx)
	if err != nil {
		return err
	}
	a.ID = e.ID
	return nil
}

func (r *ArticleRepository) Update(ctx context.Context, a *Article) error {
	_, err := r.client.Article.UpdateOneID(a.ID).
		SetTitle(a.Title).
		SetNillableContent(a.Content).
		SetNillableSummary(a.Summary).
		SetStatus(a.Status).
		SetNillableCoverDescription(a.CoverDescription).
		SetNillableCoverURL(a.CoverUrl).
		SetNillableFirstPublishedAt(a.FirstPublishedAt).
		Save(ctx)
	return err
}

func (r *ArticleRepository) Delete(ctx context.Context, id int) error {
	return r.client.Article.DeleteOneID(id).Exec(ctx)
}

func (r *ArticleRepository) List(ctx context.Context, keyword, status string, page, pageSize int) ([]*ArticleDetail, int64, error) {
	q := r.client.Article.Query()

	if keyword != "" {
		q.Where(article.Or(
			article.TitleContains(keyword),
			article.ArticleNoContains(keyword),
		))
	}
	if status != "" {
		q.Where(article.StatusEQ(status))
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	es, err := q.
		Order(article.ByID(sql.OrderDesc())).
		Limit(pageSize).
		Offset(offset).
		WithAuthor().
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return toArticleDetails(es), int64(total), nil
}

func (r *ArticleRepository) GetMaxVersionNo(ctx context.Context, articleID int) (int, error) {
	var maxV int
	err := r.client.ArticleVersion.Query().
		Where(articleversion.ArticleIDEQ(articleID)).
		Aggregate(ent.Max(articleversion.FieldVersionNo)).
		Scan(ctx, &maxV)
	if err != nil {
		if ent.IsNotFound(err) {
			return 0, nil
		}
		return 0, err
	}
	return maxV, nil
}

func toArticleVersion(e *ent.ArticleVersion) *ArticleVersion {
	return &ArticleVersion{
		ID:               e.ID,
		ArticleID:        e.ArticleID,
		VersionNo:        e.VersionNo,
		Title:            e.Title,
		Content:          stringPtr(e.Content),
		CoverDescription: stringPtr(e.CoverDescription),
		Summary:          stringPtr(e.Summary),
		Status:           e.Status,
		EditorID:         intPtr(e.EditorID),
		EditReason:       stringPtr(e.EditReason),
		CreatedAt:        e.CreatedAt,
	}
}

func toArticleVersionDetail(e *ent.ArticleVersion) *ArticleVersionDetail {
	d := &ArticleVersionDetail{
		ArticleVersion: *toArticleVersion(e),
	}
	if e.Edges.Editor != nil {
		d.EditorName = stringPtr(e.Edges.Editor.Username)
		d.EditorNickname = stringPtr(e.Edges.Editor.Nickname)
	}
	return d
}

func toArticleVersionDetails(es []*ent.ArticleVersion) []*ArticleVersionDetail {
	result := make([]*ArticleVersionDetail, len(es))
	for i, e := range es {
		result[i] = toArticleVersionDetail(e)
	}
	return result
}

func (r *ArticleRepository) CreateVersion(ctx context.Context, version *ArticleVersion) error {
	e, err := r.client.ArticleVersion.Create().
		SetArticleID(version.ArticleID).
		SetVersionNo(version.VersionNo).
		SetTitle(version.Title).
		SetNillableContent(version.Content).
		SetNillableCoverDescription(version.CoverDescription).
		SetNillableSummary(version.Summary).
		SetStatus(version.Status).
		SetNillableEditorID(version.EditorID).
		SetNillableEditReason(version.EditReason).
		Save(ctx)
	if err != nil {
		return err
	}
	version.ID = e.ID
	return nil
}

func (r *ArticleRepository) GetVersions(ctx context.Context, articleID int, page, pageSize int) ([]*ArticleVersionDetail, int64, error) {
	q := r.client.ArticleVersion.Query().
		Where(articleversion.ArticleIDEQ(articleID))

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	es, err := q.
		Order(ent.Desc(articleversion.FieldVersionNo)).
		Limit(pageSize).
		Offset(offset).
		WithEditor().
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return toArticleVersionDetails(es), int64(total), nil
}

func (r *ArticleRepository) GetVersionByID(ctx context.Context, articleID, versionID int) (*ArticleVersionDetail, error) {
	e, err := r.client.ArticleVersion.Query().
		Where(
			articleversion.ArticleIDEQ(articleID),
			articleversion.IDEQ(versionID),
		).
		WithEditor().
		First(ctx)
	if err != nil {
		return nil, err
	}
	return toArticleVersionDetail(e), nil
}
