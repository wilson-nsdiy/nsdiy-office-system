package service

import (
	"context"
	"time"

	"oa-nsdiy/backend/internal/domain"
)

// ApiTokenRepository defines the interface for API token data access required by ApiTokenService.
type ApiTokenRepository interface {
	ListByUserID(ctx context.Context, userID int, keyword, status string, page, pageSize int) ([]*domain.ApiToken, int64, error)
	GetByID(ctx context.Context, id int) (*domain.ApiToken, error)
	Create(ctx context.Context, token *domain.ApiToken) error
	Update(ctx context.Context, token *domain.ApiToken) error
	Delete(ctx context.Context, id int) error
	ValidateToken(ctx context.Context, tokenHash string) (*domain.ApiToken, error)
	UpdateUsage(ctx context.Context, id int) error
}

type ApiTokenService struct {
	repo        ApiTokenRepository
	authService *AuthService
}

func NewApiTokenService(repo ApiTokenRepository, authService *AuthService) *ApiTokenService {
	return &ApiTokenService{repo: repo, authService: authService}
}

type ApiTokenCreateInput struct {
	Name      string
	ExpiresAt string
}

type ApiTokenUpdateInput struct {
	Name   string
	Status string
}

type ApiTokenListResult struct {
	Items []*domain.ApiToken
	Total int64
}

type ApiTokenCreateResult struct {
	Token    *domain.ApiToken
	RawToken string
}

func (s *ApiTokenService) ListTokens(ctx context.Context, userID int, keyword, status string, page, pageSize int) (*ApiTokenListResult, error) {
	items, total, err := s.repo.ListByUserID(ctx, userID, keyword, status, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &ApiTokenListResult{Items: items, Total: total}, nil
}

func (s *ApiTokenService) GetToken(ctx context.Context, id, userID int) (*domain.ApiToken, error) {
	token, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, HandleRepoErr(err, "token_not_found", "API token not found")
	}

	// Check ownership
	if token.UserID != userID {
		return nil, NewServiceError(403, "access_denied", "Access denied")
	}

	return token, nil
}

func (s *ApiTokenService) CreateToken(ctx context.Context, userID int, input ApiTokenCreateInput) (*ApiTokenCreateResult, error) {
	// Generate token
	rawToken, prefix, hash := s.authService.GenerateApiToken(ctx)

	token := &domain.ApiToken{
		UserID:      userID,
		Name:        input.Name,
		TokenHash:   hash,
		TokenPrefix: prefix,
		Status:      "active",
	}

	// Parse expires_at if provided
	if input.ExpiresAt != "" {
		expiresAt, err := time.Parse(time.RFC3339, input.ExpiresAt)
		if err != nil {
			return nil, BadRequestErr("invalid_expires_at", "Invalid expires_at format, expected RFC3339")
		}
		if expiresAt.Before(time.Now()) {
			return nil, BadRequestErr("invalid_expires_at", "Expiration time must be in the future")
		}
		token.ExpiresAt = &expiresAt
	}

	if err := s.repo.Create(ctx, token); err != nil {
		return nil, err
	}

	return &ApiTokenCreateResult{
		Token:    token,
		RawToken: rawToken,
	}, nil
}

func (s *ApiTokenService) UpdateToken(ctx context.Context, id, userID int, input ApiTokenUpdateInput) (*domain.ApiToken, error) {
	token, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, HandleRepoErr(err, "token_not_found", "API token not found")
	}

	// Check ownership
	if token.UserID != userID {
		return nil, NewServiceError(403, "access_denied", "Access denied")
	}

	if input.Name != "" {
		token.Name = input.Name
	}
	if input.Status != "" {
		token.Status = input.Status
	}

	if err := s.repo.Update(ctx, token); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *ApiTokenService) DeleteToken(ctx context.Context, id, userID int) error {
	token, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return HandleRepoErr(err, "token_not_found", "API token not found")
	}

	// Check ownership
	if token.UserID != userID {
		return NewServiceError(403, "access_denied", "Access denied")
	}

	return s.repo.Delete(ctx, id)
}

func (s *ApiTokenService) ValidateToken(ctx context.Context, tokenHash string) (*domain.ApiToken, error) {
	return s.repo.ValidateToken(ctx, tokenHash)
}

func (s *ApiTokenService) UpdateTokenUsage(ctx context.Context, id int) error {
	return s.repo.UpdateUsage(ctx, id)
}
