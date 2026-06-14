package repository

import (
	"context"
	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/ent/mediaaccount"
	"oa-nsdiy/backend/internal/domain"
)

type MediaAccountRepository struct {
	client *ent.Client
}

func NewMediaAccountRepository(client *ent.Client) *MediaAccountRepository {
	return &MediaAccountRepository{client: client}
}

// Type alias for backward compatibility
type MediaAccount = domain.MediaAccount

func (r *MediaAccountRepository) GetByID(ctx context.Context, id int) (*MediaAccount, error) {
	e, err := r.client.MediaAccount.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toMediaAccountEntity(e), nil
}

func (r *MediaAccountRepository) Create(ctx context.Context, account *MediaAccount) error {
	e, err := r.client.MediaAccount.Create().
		SetName(account.Name).
		SetPlatform(account.Platform).
		SetAccountID(account.AccountID).
		SetNillableAvatar(account.Avatar).
		SetStatus(account.Status).
		SetNillableAccessToken(account.AccessToken).
		SetNillableRefreshToken(account.RefreshToken).
		SetNillableTokenExpiresAt(account.TokenExpiresAt).
		Save(ctx)
	if err != nil {
		return err
	}
	account.ID = e.ID
	return nil
}

func (r *MediaAccountRepository) Update(ctx context.Context, account *MediaAccount) error {
	_, err := r.client.MediaAccount.UpdateOneID(account.ID).
		SetName(account.Name).
		SetStatus(account.Status).
		SetNillableAccessToken(account.AccessToken).
		SetNillableRefreshToken(account.RefreshToken).
		SetNillableTokenExpiresAt(account.TokenExpiresAt).
		Save(ctx)
	return err
}

func (r *MediaAccountRepository) Delete(ctx context.Context, id int) error {
	return r.client.MediaAccount.DeleteOneID(id).Exec(ctx)
}

func (r *MediaAccountRepository) List(ctx context.Context, keyword, platform, status string, page, pageSize int) ([]*MediaAccount, int64, error) {
	q := r.client.MediaAccount.Query()

	if keyword != "" {
		q = q.Where(mediaaccount.Or(
			mediaaccount.NameContainsFold(keyword),
			mediaaccount.AccountIDContainsFold(keyword),
		))
	}

	if platform != "" {
		q = q.Where(mediaaccount.PlatformEQ(platform))
	}

	if status != "" {
		q = q.Where(mediaaccount.StatusEQ(status))
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	es, err := q.
		Order(ent.Desc(mediaaccount.FieldID)).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return toMediaAccountEntities(es), int64(total), nil
}

func toMediaAccountEntity(e *ent.MediaAccount) *MediaAccount {
	if e == nil {
		return nil
	}
	return &MediaAccount{
		ID:             e.ID,
		Name:           e.Name,
		Platform:       e.Platform,
		AccountID:      e.AccountID,
		Avatar:         stringPtr(e.Avatar),
		Status:         e.Status,
		AccessToken:    stringPtr(e.AccessToken),
		RefreshToken:   stringPtr(e.RefreshToken),
		TokenExpiresAt: timePtr(e.TokenExpiresAt),
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

func toMediaAccountEntities(es []*ent.MediaAccount) []*MediaAccount {
	result := make([]*MediaAccount, len(es))
	for i, e := range es {
		result[i] = toMediaAccountEntity(e)
	}
	return result
}
