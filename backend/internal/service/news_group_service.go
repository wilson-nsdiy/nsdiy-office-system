package service

import (
	"context"

	"oa-nsdiy/backend/internal/domain"
)

// NewsGroupRepository defines the interface for news group data access required by NewsGroupService.
type NewsGroupRepository interface {
	ListAll(ctx context.Context) ([]*domain.NewsGroup, error)
	GetByID(ctx context.Context, id int) (*domain.NewsGroup, error)
	GetByName(ctx context.Context, name string) (*domain.NewsGroup, error)
	Create(ctx context.Context, group *domain.NewsGroup) error
	Update(ctx context.Context, group *domain.NewsGroup) error
	Delete(ctx context.Context, id int) error
	GetNewsCount(ctx context.Context, groupID int) (int, error)
}

type NewsGroupService struct {
	repo NewsGroupRepository
}

func NewNewsGroupService(repo NewsGroupRepository) *NewsGroupService {
	return &NewsGroupService{repo: repo}
}

type NewsGroupCreateInput struct {
	Name        string
	Description string
	SortOrder   int
}

type NewsGroupUpdateInput struct {
	Name        string
	Description string
	SortOrder   int
}

func (s *NewsGroupService) GetAllGroups(ctx context.Context) ([]*domain.NewsGroup, error) {
	return s.repo.ListAll(ctx)
}

func (s *NewsGroupService) GetGroup(ctx context.Context, id int) (*domain.NewsGroup, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *NewsGroupService) CreateGroup(ctx context.Context, input NewsGroupCreateInput) (*domain.NewsGroup, error) {
	existing, _ := s.repo.GetByName(ctx, input.Name)
	if existing != nil {
		return nil, NewServiceError(409, "name_conflict", "Group name already exists")
	}

	group := &domain.NewsGroup{
		Name:        input.Name,
		Description: &input.Description,
		SortOrder:   input.SortOrder,
	}

	if err := s.repo.Create(ctx, group); err != nil {
		return nil, err
	}

	return group, nil
}

func (s *NewsGroupService) UpdateGroup(ctx context.Context, id int, input NewsGroupUpdateInput) (*domain.NewsGroup, error) {
	group, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, HandleRepoErr(err, "group_not_found", "News group not found")
	}

	if input.Name != "" && input.Name != group.Name {
		existing, _ := s.repo.GetByName(ctx, input.Name)
		if existing != nil {
			return nil, NewServiceError(409, "name_conflict", "Group name already exists")
		}
		group.Name = input.Name
	}

	if input.Description != "" {
		group.Description = &input.Description
	}

	group.SortOrder = input.SortOrder

	if err := s.repo.Update(ctx, group); err != nil {
		return nil, err
	}

	return group, nil
}

func (s *NewsGroupService) DeleteGroup(ctx context.Context, id int) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return HandleRepoErr(err, "group_not_found", "News group not found")
	}

	count, err := s.repo.GetNewsCount(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return NewServiceError(409, "group_has_news", "Group has news items")
	}

	return s.repo.Delete(ctx, id)
}
