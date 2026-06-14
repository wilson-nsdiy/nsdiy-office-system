package service

import (
	"context"

	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/internal/domain"
)

// NewsRepository defines the interface for news data access required by NewsService.
type NewsRepository interface {
	List(ctx context.Context, groupID *int, keyword string, page, pageSize int) ([]*domain.NewsDetail, int64, error)
	GetDetailByID(ctx context.Context, id int) (*domain.NewsDetail, error)
	GetByID(ctx context.Context, id int) (*domain.News, error)
	Create(ctx context.Context, news *domain.News) error
	Update(ctx context.Context, news *domain.News) error
	Delete(ctx context.Context, id int) error
}

// NewsGroupValidator defines the minimal news group interface used by NewsService for validation.
type NewsGroupValidator interface {
	GetByID(ctx context.Context, id int) (*domain.NewsGroup, error)
}

type NewsService struct {
	repo      NewsRepository
	groupRepo NewsGroupValidator
}

func NewNewsService(repo NewsRepository, groupRepo NewsGroupValidator) *NewsService {
	return &NewsService{repo: repo, groupRepo: groupRepo}
}

type NewsCreateInput struct {
	GroupID int
	Title   string
	Content string
}

type NewsUpdateInput struct {
	GroupID *int
	Title   string
	Content string
}

type NewsListResult struct {
	Items []*domain.NewsDetail
	Total int64
}

func (s *NewsService) GetNewsList(ctx context.Context, groupID *int, keyword string, page, pageSize int) (*NewsListResult, error) {
	items, total, err := s.repo.List(ctx, groupID, keyword, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &NewsListResult{Items: items, Total: total}, nil
}

func (s *NewsService) GetNews(ctx context.Context, id int) (*domain.NewsDetail, error) {
	news, err := s.repo.GetDetailByID(ctx, id)
	if err != nil {
		return nil, HandleRepoErr(err, "news_not_found", "News not found")
	}
	return news, nil
}

func (s *NewsService) CreateNews(ctx context.Context, input NewsCreateInput, creatorID int) (*domain.News, error) {
	// Validate group exists
	_, err := s.groupRepo.GetByID(ctx, input.GroupID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, BadRequestErr("group_not_found", "News group not found")
		}
		return nil, err
	}

	news := &domain.News{
		GroupID:   input.GroupID,
		Title:     input.Title,
		Content:   &input.Content,
		CreatorID: creatorID,
	}

	if err := s.repo.Create(ctx, news); err != nil {
		return nil, err
	}

	return news, nil
}

func (s *NewsService) UpdateNews(ctx context.Context, id int, input NewsUpdateInput) (*domain.NewsDetail, error) {
	news, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, HandleRepoErr(err, "news_not_found", "News not found")
	}

	if input.GroupID != nil {
		// Validate group exists
		_, err := s.groupRepo.GetByID(ctx, *input.GroupID)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, BadRequestErr("group_not_found", "News group not found")
			}
			return nil, err
		}
		news.GroupID = *input.GroupID
	}

	if input.Title != "" {
		news.Title = input.Title
	}

	if input.Content != "" {
		news.Content = &input.Content
	}

	if err := s.repo.Update(ctx, news); err != nil {
		return nil, err
	}

	return s.repo.GetDetailByID(ctx, id)
}

func (s *NewsService) DeleteNews(ctx context.Context, id int) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return HandleRepoErr(err, "news_not_found", "News not found")
	}

	return s.repo.Delete(ctx, id)
}
