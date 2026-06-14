package repository

import (
	"context"

	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/ent/mediacontent"
	"oa-nsdiy/backend/ent/mediacontentversion"
	"oa-nsdiy/backend/internal/domain"
)

type MediaContentRepository struct {
	client *ent.Client
}

func NewMediaContentRepository(client *ent.Client) *MediaContentRepository {
	return &MediaContentRepository{client: client}
}

// Type aliases for backward compatibility
type MediaContent = domain.MediaContent
type MediaContentDetail = domain.MediaContentDetail
type MediaContentVersion = domain.MediaContentVersion
type MediaContentVersionDetail = domain.MediaContentVersionDetail

func (r *MediaContentRepository) GetByID(ctx context.Context, id int) (*MediaContent, error) {
	e, err := r.client.MediaContent.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toMediaContent(e), nil
}

func (r *MediaContentRepository) GetDetailByID(ctx context.Context, id int) (*MediaContentDetail, error) {
	e, err := r.client.MediaContent.Query().
		Where(mediacontent.IDEQ(id)).
		WithAccount().
		First(ctx)
	if err != nil {
		return nil, err
	}
	detail := &MediaContentDetail{MediaContent: *toMediaContent(e)}
	if e.Edges.Account != nil {
		detail.AccountName = stringPtr(e.Edges.Account.Name)
	}
	return detail, nil
}

func (r *MediaContentRepository) Create(ctx context.Context, c *MediaContent) error {
	e, err := r.client.MediaContent.Create().
		SetTitle(c.Title).
		SetNillableContent(c.Content).
		SetNillableCoverImage(c.CoverImage).
		SetPlatform(c.Platform).
		SetNillableAccountID(c.AccountID).
		SetStatus(c.Status).
		SetViews(c.Views).
		SetLikes(c.Likes).
		SetComments(c.Comments).
		SetShares(c.Shares).
		SetNillablePublishTime(c.PublishTime).
		Save(ctx)
	if err != nil {
		return err
	}
	c.ID = e.ID
	return nil
}

func (r *MediaContentRepository) Update(ctx context.Context, c *MediaContent) error {
	_, err := r.client.MediaContent.UpdateOneID(c.ID).
		SetTitle(c.Title).
		SetNillableContent(c.Content).
		SetNillableCoverImage(c.CoverImage).
		SetNillableAccountID(c.AccountID).
		SetStatus(c.Status).
		SetViews(c.Views).
		SetLikes(c.Likes).
		SetComments(c.Comments).
		SetShares(c.Shares).
		SetNillablePublishTime(c.PublishTime).
		Save(ctx)
	return err
}

func (r *MediaContentRepository) Delete(ctx context.Context, id int) error {
	return r.client.MediaContent.DeleteOneID(id).Exec(ctx)
}

func (r *MediaContentRepository) List(ctx context.Context, keyword, platform, status string, accountID *int, page, pageSize int) ([]*MediaContentDetail, int64, error) {

	q := r.client.MediaContent.Query()
	if keyword != "" {
		q.Where(mediacontent.TitleContains(keyword))
	}
	if platform != "" {
		q.Where(mediacontent.PlatformEQ(platform))
	}
	if status != "" {
		q.Where(mediacontent.StatusEQ(status))
	}
	if accountID != nil {
		q.Where(mediacontent.AccountIDEQ(*accountID))
	}

	count, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	total := int64(count)

	offset := (page - 1) * pageSize
	entities, err := q.WithAccount().
		Order(ent.Desc(mediacontent.FieldID)).
		Limit(pageSize).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	var contents []*MediaContentDetail
	for _, e := range entities {
		detail := &MediaContentDetail{MediaContent: *toMediaContent(e)}
		if e.Edges.Account != nil {
			detail.AccountName = stringPtr(e.Edges.Account.Name)
		}
		contents = append(contents, detail)
	}

	return contents, total, nil
}

func (r *MediaContentRepository) GetMaxVersionNo(ctx context.Context, contentID int) (int, error) {
	var maxV int
	err := r.client.MediaContentVersion.Query().
		Where(mediacontentversion.ContentIDEQ(contentID)).
		Aggregate(ent.Max(mediacontentversion.FieldVersionNo)).
		Scan(ctx, &maxV)
	if err != nil {
		return 0, err
	}
	return maxV, nil
}

func (r *MediaContentRepository) CreateVersion(ctx context.Context, v *MediaContentVersion) error {
	e, err := r.client.MediaContentVersion.Create().
		SetContentID(v.ContentID).
		SetVersionNo(v.VersionNo).
		SetTitle(v.Title).
		SetNillableContent(v.Content).
		SetNillableCoverImage(v.CoverImage).
		SetStatus(v.Status).
		SetNillableEditorID(v.EditorID).
		SetNillableEditReason(v.EditReason).
		Save(ctx)
	if err != nil {
		return err
	}
	v.ID = e.ID
	return nil
}

func (r *MediaContentRepository) ListVersions(ctx context.Context, contentID int, page, pageSize int) ([]*MediaContentVersionDetail, int64, error) {

	q := r.client.MediaContentVersion.Query().
		Where(mediacontentversion.ContentIDEQ(contentID))

	count, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	total := int64(count)

	offset := (page - 1) * pageSize
	entities, err := q.WithEditor().
		Order(ent.Desc(mediacontentversion.FieldVersionNo)).
		Limit(pageSize).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	var versions []*MediaContentVersionDetail
	for _, e := range entities {
		detail := &MediaContentVersionDetail{
			MediaContentVersion: *toMediaContentVersion(e),
		}
		if e.Edges.Editor != nil {
			detail.EditorName = stringPtr(e.Edges.Editor.Username)
			detail.EditorNickname = stringPtr(e.Edges.Editor.Nickname)
		}
		versions = append(versions, detail)
	}

	return versions, total, nil
}

func (r *MediaContentRepository) GetVersionByID(ctx context.Context, versionID int) (*MediaContentVersionDetail, error) {
	e, err := r.client.MediaContentVersion.Query().
		Where(mediacontentversion.IDEQ(versionID)).
		WithEditor().
		First(ctx)
	if err != nil {
		return nil, err
	}
	detail := &MediaContentVersionDetail{
		MediaContentVersion: *toMediaContentVersion(e),
	}
	if e.Edges.Editor != nil {
		detail.EditorName = stringPtr(e.Edges.Editor.Username)
		detail.EditorNickname = stringPtr(e.Edges.Editor.Nickname)
	}
	return detail, nil
}

func (r *MediaContentRepository) DeleteVersions(ctx context.Context, contentID int) error {
	_, err := r.client.MediaContentVersion.Delete().
		Where(mediacontentversion.ContentIDEQ(contentID)).
		Exec(ctx)
	return err
}

// --- Converters ---

func toMediaContent(e *ent.MediaContent) *MediaContent {
	return &MediaContent{
		ID:          e.ID,
		Title:       e.Title,
		Content:     stringPtr(e.Content),
		CoverImage:  stringPtr(e.CoverImage),
		Platform:    e.Platform,
		AccountID:   intPtr(e.AccountID),
		Status:      e.Status,
		Views:       e.Views,
		Likes:       e.Likes,
		Comments:    e.Comments,
		Shares:      e.Shares,
		PublishTime: timePtr(e.PublishTime),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toMediaContentVersion(e *ent.MediaContentVersion) *MediaContentVersion {
	return &MediaContentVersion{
		ID:         e.ID,
		ContentID:  e.ContentID,
		VersionNo:  e.VersionNo,
		Title:      e.Title,
		Content:    stringPtr(e.Content),
		CoverImage: stringPtr(e.CoverImage),
		Status:     e.Status,
		EditorID:   intPtr(e.EditorID),
		EditReason: stringPtr(e.EditReason),
		CreatedAt:  e.CreatedAt,
	}
}
