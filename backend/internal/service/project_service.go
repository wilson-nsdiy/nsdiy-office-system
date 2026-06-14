package service

import (
	"context"

	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/internal/domain"
)

// ProjectRepository defines the interface for project data access required by ProjectService.
type ProjectRepository interface {
	ListByUserID(ctx context.Context, userID int, keyword string, page, pageSize int) ([]*domain.ProjectDetail, int64, error)
	GetByID(ctx context.Context, id int) (*domain.Project, error)
	GetByProjectNo(ctx context.Context, projectNo string) (*domain.Project, error)
	GenerateUniqueProjectNo(ctx context.Context) (string, error)
	Create(ctx context.Context, project *domain.Project) error
	Update(ctx context.Context, project *domain.Project) error
	Delete(ctx context.Context, id int) error
	AddMember(ctx context.Context, member *domain.ProjectMember) error
	GetMembers(ctx context.Context, projectID int) ([]*domain.ProjectMemberDetail, error)
	GetMember(ctx context.Context, projectID, userID int) (*domain.ProjectMember, error)
	UpdateMemberRole(ctx context.Context, projectID, userID int, role string) error
	RemoveMember(ctx context.Context, projectID, userID int) error
	IsMember(ctx context.Context, projectID, userID int) (bool, error)
}

type ProjectService struct {
	repo ProjectRepository
}

func NewProjectService(repo ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

type ProjectCreateInput struct {
	Name              string
	Description       string
	Priority          string
	ExpectedStartDate string
	ExpectedEndDate   string
}

type ProjectUpdateInput struct {
	Name              string
	Description       string
	Status            string
	Priority          string
	ExpectedStartDate string
	ExpectedEndDate   string
	StartDate         string
	EndDate           string
}

type ProjectListResult struct {
	Items []*domain.ProjectDetail
	Total int64
}

func (s *ProjectService) GetProjectList(ctx context.Context, userID int, keyword string, page, pageSize int) (*ProjectListResult, error) {
	items, total, err := s.repo.ListByUserID(ctx, userID, keyword, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &ProjectListResult{Items: items, Total: total}, nil
}

func (s *ProjectService) GetProject(ctx context.Context, projectNo string) (*domain.Project, error) {
	project, err := s.repo.GetByProjectNo(ctx, projectNo)
	if err != nil {
		return nil, HandleRepoErr(err, "project_not_found", "Project not found")
	}
	return project, nil
}

func (s *ProjectService) CreateProject(ctx context.Context, input ProjectCreateInput, ownerID int) (*domain.Project, error) {
	if input.Name == "" {
		return nil, NewServiceError(400, "name_required", "Name is required")
	}

	projectNo, err := s.repo.GenerateUniqueProjectNo(ctx)
	if err != nil {
		return nil, err
	}

	project := &domain.Project{
		Name:        input.Name,
		ProjectNo:   projectNo,
		Description: &input.Description,
		Status:      "TODO",
		Priority:    "MEDIUM",
		OwnerID:     ownerID,
	}

	if input.Priority != "" {
		project.Priority = input.Priority
	}

	if err := s.repo.Create(ctx, project); err != nil {
		return nil, err
	}

	// Auto-add creator as OWNER member
	member := &domain.ProjectMember{
		ProjectID: project.ID,
		UserID:    ownerID,
		Role:      "OWNER",
	}
	_ = s.repo.AddMember(ctx, member)

	return project, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, projectNo string, input ProjectUpdateInput) (*domain.Project, error) {
	project, err := s.repo.GetByProjectNo(ctx, projectNo)
	if err != nil {
		return nil, HandleRepoErr(err, "project_not_found", "Project not found")
	}

	if input.Name != "" {
		project.Name = input.Name
	}
	if input.Description != "" {
		project.Description = &input.Description
	}
	if input.Status != "" {
		project.Status = input.Status
	}
	if input.Priority != "" {
		project.Priority = input.Priority
	}

	if err := s.repo.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, projectNo string) error {
	project, err := s.repo.GetByProjectNo(ctx, projectNo)
	if err != nil {
		return HandleRepoErr(err, "project_not_found", "Project not found")
	}

	return s.repo.Delete(ctx, project.ID)
}

func (s *ProjectService) GetMembers(ctx context.Context, projectNo string) ([]*domain.ProjectMemberDetail, error) {
	project, err := s.repo.GetByProjectNo(ctx, projectNo)
	if err != nil {
		return nil, HandleRepoErr(err, "project_not_found", "Project not found")
	}

	return s.repo.GetMembers(ctx, project.ID)
}

func (s *ProjectService) AddMember(ctx context.Context, projectNo string, userID int, role string) error {
	project, err := s.repo.GetByProjectNo(ctx, projectNo)
	if err != nil {
		return HandleRepoErr(err, "project_not_found", "Project not found")
	}

	// Check if already a member
	existing, _ := s.repo.GetMember(ctx, project.ID, userID)
	if existing != nil {
		return NewServiceError(409, "already_member", "User is already a member")
	}

	if role == "" {
		role = "MEMBER"
	}

	member := &domain.ProjectMember{
		ProjectID: project.ID,
		UserID:    userID,
		Role:      role,
	}

	return s.repo.AddMember(ctx, member)
}

func (s *ProjectService) UpdateMemberRole(ctx context.Context, projectNo string, userID int, role string) error {
	project, err := s.repo.GetByProjectNo(ctx, projectNo)
	if err != nil {
		return HandleRepoErr(err, "project_not_found", "Project not found")
	}

	return s.repo.UpdateMemberRole(ctx, project.ID, userID, role)
}

func (s *ProjectService) RemoveMember(ctx context.Context, projectNo string, userID int) error {
	project, err := s.repo.GetByProjectNo(ctx, projectNo)
	if err != nil {
		return HandleRepoErr(err, "project_not_found", "Project not found")
	}

	// Block removing project owner
	if project.OwnerID == userID {
		return NewServiceError(403, "cannot_remove_owner", "Cannot remove project owner")
	}

	return s.repo.RemoveMember(ctx, project.ID, userID)
}

func (s *ProjectService) IsMember(ctx context.Context, projectNo string, userID int) (bool, error) {
	project, err := s.repo.GetByProjectNo(ctx, projectNo)
	if err != nil {
		return false, HandleRepoErr(err, "project_not_found", "Project not found")
	}

	return s.repo.IsMember(ctx, project.ID, userID)
}

func (s *ProjectService) GetUserRole(ctx context.Context, projectNo string, userID int) (string, error) {
	project, err := s.repo.GetByProjectNo(ctx, projectNo)
	if err != nil {
		return "", HandleRepoErr(err, "project_not_found", "Project not found")
	}

	// Check if user is project owner
	if project.OwnerID == userID {
		return "OWNER", nil
	}

	// Get member role
	member, err := s.repo.GetMember(ctx, project.ID, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return "", nil
		}
		return "", err
	}

	return member.Role, nil
}
