package service

import (
	"context"

	"oa-nsdiy/backend/internal/domain"
)

// MediaAccountRepository defines the interface for media account data access required by MediaAccountService.
type MediaAccountRepository interface {
	List(ctx context.Context, keyword, platform, status string, page, pageSize int) ([]*domain.MediaAccount, int64, error)
	GetByID(ctx context.Context, id int) (*domain.MediaAccount, error)
	Create(ctx context.Context, account *domain.MediaAccount) error
	Update(ctx context.Context, account *domain.MediaAccount) error
	Delete(ctx context.Context, id int) error
}

type MediaAccountService struct {
	repo MediaAccountRepository
}

func NewMediaAccountService(repo MediaAccountRepository) *MediaAccountService {
	return &MediaAccountService{repo: repo}
}

type MediaAccountCreateInput struct {
	Name      string
	Platform  string
	AccountId string
	Avatar    string
}

type MediaAccountUpdateInput struct {
	Name           string
	Status         string
	AccessToken    string
	RefreshToken   string
	TokenExpiresAt string
}

type MediaAccountListResult struct {
	Items []*domain.MediaAccount
	Total int64
}

func (s *MediaAccountService) ListAccounts(ctx context.Context, keyword, platform, status string, page, pageSize int) (*MediaAccountListResult, error) {
	items, total, err := s.repo.List(ctx, keyword, platform, status, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &MediaAccountListResult{Items: items, Total: total}, nil
}

func (s *MediaAccountService) GetAccount(ctx context.Context, id int) (*domain.MediaAccount, error) {
	account, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, HandleRepoErr(err, "account_not_found", "Media account not found")
	}
	return account, nil
}

func (s *MediaAccountService) CreateAccount(ctx context.Context, input MediaAccountCreateInput) (*domain.MediaAccount, error) {
	account := &domain.MediaAccount{
		Name:      input.Name,
		Platform:  input.Platform,
		AccountID: input.AccountId,
		Avatar:    &input.Avatar,
		Status:    "active",
	}

	if err := s.repo.Create(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *MediaAccountService) UpdateAccount(ctx context.Context, id int, input MediaAccountUpdateInput) (*domain.MediaAccount, error) {
	account, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, HandleRepoErr(err, "account_not_found", "Media account not found")
	}

	if input.Name != "" {
		account.Name = input.Name
	}
	if input.Status != "" {
		account.Status = input.Status
	}

	if err := s.repo.Update(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *MediaAccountService) DeleteAccount(ctx context.Context, id int) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return HandleRepoErr(err, "account_not_found", "Media account not found")
	}

	return s.repo.Delete(ctx, id)
}

// MediaContentRepository defines the interface for media content data access required by MediaContentService.
type MediaContentRepository interface {
	List(ctx context.Context, keyword, platform, status string, accountID *int, page, pageSize int) ([]*domain.MediaContentDetail, int64, error)
	GetDetailByID(ctx context.Context, id int) (*domain.MediaContentDetail, error)
	GetByID(ctx context.Context, id int) (*domain.MediaContent, error)
	Create(ctx context.Context, content *domain.MediaContent) error
	Update(ctx context.Context, content *domain.MediaContent) error
	Delete(ctx context.Context, id int) error
	GetMaxVersionNo(ctx context.Context, contentID int) (int, error)
	CreateVersion(ctx context.Context, version *domain.MediaContentVersion) error
	DeleteVersions(ctx context.Context, contentID int) error
	ListVersions(ctx context.Context, contentID int, page, pageSize int) ([]*domain.MediaContentVersionDetail, int64, error)
	GetVersionByID(ctx context.Context, versionID int) (*domain.MediaContentVersionDetail, error)
}

type MediaContentService struct {
	repo MediaContentRepository
}

func NewMediaContentService(repo MediaContentRepository) *MediaContentService {
	return &MediaContentService{repo: repo}
}

type MediaContentCreateInput struct {
	Title       string
	Content     string
	CoverImage  string
	Platform    string
	AccountId   *int
	Status      string
	PublishTime string
}

type MediaContentUpdateInput struct {
	Title       string
	Content     string
	CoverImage  string
	Status      string
	Views       *int
	Likes       *int
	Comments    *int
	Shares      *int
	PublishTime string
}

type MediaContentListResult struct {
	Items []*domain.MediaContentDetail
	Total int64
}

func (s *MediaContentService) ListContents(ctx context.Context, keyword, platform, status string, accountID *int, page, pageSize int) (*MediaContentListResult, error) {
	items, total, err := s.repo.List(ctx, keyword, platform, status, accountID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &MediaContentListResult{Items: items, Total: total}, nil
}

func (s *MediaContentService) GetContent(ctx context.Context, id int) (*domain.MediaContentDetail, error) {
	return s.repo.GetDetailByID(ctx, id)
}

func (s *MediaContentService) CreateContent(ctx context.Context, input MediaContentCreateInput) (*domain.MediaContent, error) {
	content := &domain.MediaContent{
		Title:      input.Title,
		Content:    &input.Content,
		CoverImage: &input.CoverImage,
		Platform:   input.Platform,
		AccountID:  input.AccountId,
		Status:     "draft",
	}

	if input.Status != "" {
		content.Status = input.Status
	}

	if err := s.repo.Create(ctx, content); err != nil {
		return nil, err
	}

	// Create initial version
	maxNo, _ := s.repo.GetMaxVersionNo(ctx, content.ID)
	version := &domain.MediaContentVersion{
		ContentID:  content.ID,
		VersionNo:  maxNo + 1,
		Title:      content.Title,
		Content:    content.Content,
		CoverImage: content.CoverImage,
		Status:     content.Status,
	}
	_ = s.repo.CreateVersion(ctx, version)

	return content, nil
}

func (s *MediaContentService) UpdateContent(ctx context.Context, id int, input MediaContentUpdateInput) (*domain.MediaContent, error) {
	content, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create version before update
	maxNo, _ := s.repo.GetMaxVersionNo(ctx, id)
	version := &domain.MediaContentVersion{
		ContentID:  id,
		VersionNo:  maxNo + 1,
		Title:      content.Title,
		Content:    content.Content,
		CoverImage: content.CoverImage,
		Status:     content.Status,
	}
	_ = s.repo.CreateVersion(ctx, version)

	// Update content
	if input.Title != "" {
		content.Title = input.Title
	}
	if input.Content != "" {
		content.Content = &input.Content
	}
	if input.CoverImage != "" {
		content.CoverImage = &input.CoverImage
	}
	if input.Status != "" {
		content.Status = input.Status
	}
	if input.Views != nil {
		content.Views = *input.Views
	}
	if input.Likes != nil {
		content.Likes = *input.Likes
	}
	if input.Comments != nil {
		content.Comments = *input.Comments
	}
	if input.Shares != nil {
		content.Shares = *input.Shares
	}

	if err := s.repo.Update(ctx, content); err != nil {
		return nil, err
	}

	return content, nil
}

func (s *MediaContentService) DeleteContent(ctx context.Context, id int) error {
	// Delete all versions first
	_ = s.repo.DeleteVersions(ctx, id)

	return s.repo.Delete(ctx, id)
}

func (s *MediaContentService) ListVersions(ctx context.Context, contentID int, page, pageSize int) ([]*domain.MediaContentVersionDetail, int64, error) {
	return s.repo.ListVersions(ctx, contentID, page, pageSize)
}

func (s *MediaContentService) GetVersion(ctx context.Context, versionID int) (*domain.MediaContentVersionDetail, error) {
	return s.repo.GetVersionByID(ctx, versionID)
}
