package service

import (
	"context"

	"oa-nsdiy/backend/internal/domain"
)

// ArticleRepository defines the interface for article data access required by ArticleService.
type ArticleRepository interface {
	List(ctx context.Context, keyword, status string, page, pageSize int) ([]*domain.ArticleDetail, int64, error)
	GetDetailByID(ctx context.Context, id int) (*domain.ArticleDetail, error)
	GetByID(ctx context.Context, id int) (*domain.Article, error)
	GetByTitle(ctx context.Context, title string) (*domain.Article, error)
	GenerateArticleNo(ctx context.Context) (string, error)
	Create(ctx context.Context, article *domain.Article) error
	Update(ctx context.Context, article *domain.Article) error
	Delete(ctx context.Context, id int) error
	GetMaxVersionNo(ctx context.Context, articleID int) (int, error)
	CreateVersion(ctx context.Context, version *domain.ArticleVersion) error
	GetVersions(ctx context.Context, articleID int, page, pageSize int) ([]*domain.ArticleVersionDetail, int64, error)
	GetVersionByID(ctx context.Context, articleID, versionID int) (*domain.ArticleVersionDetail, error)
}

type ArticleService struct {
	repo ArticleRepository
}

func NewArticleService(repo ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

type ArticleCreateInput struct {
	Title            string
	Content          string
	Summary          string
	Status           string
	CoverDescription string
	CoverUrl         string
}

type ArticleUpdateInput struct {
	Title            string
	Content          string
	Summary          string
	Status           string
	CoverDescription string
	CoverUrl         string
	EditReason       string
}

type ArticleListResult struct {
	Items []*domain.ArticleDetail
	Total int64
}

func (s *ArticleService) GetArticleList(ctx context.Context, keyword, status string, page, pageSize int) (*ArticleListResult, error) {
	items, total, err := s.repo.List(ctx, keyword, status, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &ArticleListResult{Items: items, Total: total}, nil
}

func (s *ArticleService) GetArticle(ctx context.Context, id int) (*domain.ArticleDetail, error) {
	article, err := s.repo.GetDetailByID(ctx, id)
	if err != nil {
		return nil, HandleRepoErr(err, "article_not_found", "Article not found")
	}
	return article, nil
}

func (s *ArticleService) CreateArticle(ctx context.Context, input ArticleCreateInput, authorID int) (*domain.Article, error) {
	// Validate title
	if input.Title == "" {
		return nil, NewServiceError(400, "title_required", "Title is required")
	}

	// Check title uniqueness
	existing, _ := s.repo.GetByTitle(ctx, input.Title)
	if existing != nil {
		return nil, NewServiceError(409, "title_conflict", "Article title already exists")
	}

	// Generate article number
	articleNo, err := s.repo.GenerateArticleNo(ctx)
	if err != nil {
		return nil, err
	}

	article := &domain.Article{
		ArticleNo:        articleNo,
		Title:            input.Title,
		Content:          &input.Content,
		Summary:          &input.Summary,
		Status:           "DRAFT",
		AuthorID:         authorID,
		CoverDescription: &input.CoverDescription,
		CoverUrl:         &input.CoverUrl,
	}

	if input.Status == "PUBLISHED" {
		article.Status = "PUBLISHED"
	}

	if err := s.repo.Create(ctx, article); err != nil {
		return nil, err
	}

	return article, nil
}

func (s *ArticleService) UpdateArticle(ctx context.Context, id int, input ArticleUpdateInput, editorID int) (*domain.Article, error) {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, HandleRepoErr(err, "article_not_found", "Article not found")
	}

	// Check title uniqueness if changed
	if input.Title != "" && input.Title != article.Title {
		existing, _ := s.repo.GetByTitle(ctx, input.Title)
		if existing != nil {
			return nil, NewServiceError(409, "title_conflict", "Article title already exists")
		}
		article.Title = input.Title
	}

	// Check if content fields changed for versioning
	contentChanged := (input.Content != "" && (article.Content == nil || *article.Content != input.Content)) ||
		(input.Summary != "" && (article.Summary == nil || *article.Summary != input.Summary)) ||
		(input.CoverDescription != "" && (article.CoverDescription == nil || *article.CoverDescription != input.CoverDescription))

	// Create version snapshot before update
	if contentChanged {
		maxNo, _ := s.repo.GetMaxVersionNo(ctx, id)
		version := &domain.ArticleVersion{
			ArticleID:        id,
			VersionNo:        maxNo + 1,
			Title:            article.Title,
			Content:          article.Content,
			CoverDescription: article.CoverDescription,
			Summary:          article.Summary,
			Status:           article.Status,
			EditorID:         &editorID,
			EditReason:       &input.EditReason,
		}
		_ = s.repo.CreateVersion(ctx, version)
	}

	// Update article
	if input.Content != "" {
		article.Content = &input.Content
	}
	if input.Summary != "" {
		article.Summary = &input.Summary
	}
	if input.Status != "" {
		article.Status = input.Status
	}
	if input.CoverDescription != "" {
		article.CoverDescription = &input.CoverDescription
	}
	if input.CoverUrl != "" {
		article.CoverUrl = &input.CoverUrl
	}

	if err := s.repo.Update(ctx, article); err != nil {
		return nil, err
	}

	return article, nil
}

func (s *ArticleService) DeleteArticle(ctx context.Context, id int) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return HandleRepoErr(err, "article_not_found", "Article not found")
	}

	return s.repo.Delete(ctx, id)
}

func (s *ArticleService) GetArticleVersions(ctx context.Context, articleID int, page, pageSize int) ([]*domain.ArticleVersionDetail, int64, error) {
	return s.repo.GetVersions(ctx, articleID, page, pageSize)
}

func (s *ArticleService) GetArticleVersion(ctx context.Context, articleID, versionID int) (*domain.ArticleVersionDetail, error) {
	return s.repo.GetVersionByID(ctx, articleID, versionID)
}
