package repository

import (
	"context"
	"fmt"
	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/ent/apitoken"
	"oa-nsdiy/backend/internal/domain"
	"time"
)

type ApiTokenRepository struct {
	client *ent.Client
}

func NewApiTokenRepository(client *ent.Client) *ApiTokenRepository {
	return &ApiTokenRepository{client: client}
}

// Type alias for backward compatibility
type ApiToken = domain.ApiToken

func (r *ApiTokenRepository) GetByID(ctx context.Context, id int) (*ApiToken, error) {
	e, err := r.client.ApiToken.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toApiTokenEntity(e), nil
}

func (r *ApiTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*ApiToken, error) {
	e, err := r.client.ApiToken.Query().Where(apitoken.TokenHashEQ(tokenHash)).First(ctx)
	if err != nil {
		return nil, err
	}
	return toApiTokenEntity(e), nil
}

func (r *ApiTokenRepository) Create(ctx context.Context, token *ApiToken) error {
	e, err := r.client.ApiToken.Create().
		SetUserID(token.UserID).
		SetName(token.Name).
		SetTokenHash(token.TokenHash).
		SetTokenPrefix(token.TokenPrefix).
		SetStatus(token.Status).
		SetNillableExpiresAt(token.ExpiresAt).
		Save(ctx)
	if err != nil {
		return err
	}
	token.ID = e.ID
	return nil
}

func (r *ApiTokenRepository) Update(ctx context.Context, token *ApiToken) error {
	_, err := r.client.ApiToken.UpdateOneID(token.ID).
		SetName(token.Name).
		SetStatus(token.Status).
		Save(ctx)
	return err
}

func (r *ApiTokenRepository) UpdateUsage(ctx context.Context, id int) error {
	_, err := r.client.ApiToken.UpdateOneID(id).
		AddUsageCount(1).
		SetLastUsedAt(time.Now()).
		Save(ctx)
	return err
}

func (r *ApiTokenRepository) Delete(ctx context.Context, id int) error {
	return r.client.ApiToken.DeleteOneID(id).Exec(ctx)
}

func (r *ApiTokenRepository) ListByUserID(ctx context.Context, userID int, keyword, status string, page, pageSize int) ([]*ApiToken, int64, error) {
	q := r.client.ApiToken.Query().Where(apitoken.UserIDEQ(userID))

	if keyword != "" {
		q = q.Where(apitoken.NameContainsFold(keyword))
	}

	if status != "" {
		q = q.Where(apitoken.StatusEQ(status))
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	es, err := q.
		Order(ent.Desc(apitoken.FieldID)).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return toApiTokenEntities(es), int64(total), nil
}

func (r *ApiTokenRepository) ValidateToken(ctx context.Context, tokenHash string) (*ApiToken, error) {
	token, err := r.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	if token.Status != "active" {
		return nil, fmt.Errorf("api token is not active")
	}

	if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("api token has expired")
	}

	return token, nil
}

func toApiTokenEntity(e *ent.ApiToken) *ApiToken {
	if e == nil {
		return nil
	}
	return &ApiToken{
		ID:          e.ID,
		UserID:      e.UserID,
		Name:        e.Name,
		TokenHash:   e.TokenHash,
		TokenPrefix: e.TokenPrefix,
		Status:      e.Status,
		ExpiresAt:   timePtr(e.ExpiresAt),
		LastUsedAt:  timePtr(e.LastUsedAt),
		UsageCount:  e.UsageCount,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toApiTokenEntities(es []*ent.ApiToken) []*ApiToken {
	result := make([]*ApiToken, len(es))
	for i, e := range es {
		result[i] = toApiTokenEntity(e)
	}
	return result
}
