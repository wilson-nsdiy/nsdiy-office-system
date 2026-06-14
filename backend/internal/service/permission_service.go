package service

import (
	"context"

	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/internal/domain"
)

// PermissionRepository defines the interface for permission data access required by PermissionService.
type PermissionRepository interface {
	Search(ctx context.Context, resourceType, keyword string) ([]*domain.Permission, error)
	ListActive(ctx context.Context) ([]*domain.Permission, error)
	ListAll(ctx context.Context) ([]*domain.Permission, error)
	GetByID(ctx context.Context, id int) (*domain.Permission, error)
	GetByName(ctx context.Context, name string) (*domain.Permission, error)
	Create(ctx context.Context, perm *domain.Permission) error
	Update(ctx context.Context, perm *domain.Permission) error
	Delete(ctx context.Context, id int) error
	IsUsedByRoles(ctx context.Context, permID int) (bool, error)
}

type PermissionService struct {
	repo PermissionRepository
}

func NewPermissionService(repo PermissionRepository) *PermissionService {
	return &PermissionService{repo: repo}
}

type PermissionCreateInput struct {
	Pid          *int
	Name         string
	ResourceType string
	ResourcePath string
	HTTPMethod   string
	Description  string
	IsActive     bool
}

type PermissionUpdateInput struct {
	Pid          *int
	Name         string
	ResourceType string
	ResourcePath string
	HTTPMethod   string
	Description  string
	IsActive     bool
}

func (s *PermissionService) GetPermissions(ctx context.Context, resourceType, keyword string) ([]*domain.Permission, error) {
	if resourceType != "" || keyword != "" {
		return s.repo.Search(ctx, resourceType, keyword)
	}
	return s.repo.ListActive(ctx)
}

func (s *PermissionService) GetAllPermissions(ctx context.Context) ([]*domain.Permission, error) {
	return s.repo.ListAll(ctx)
}

func (s *PermissionService) GetPermission(ctx context.Context, id int) (*domain.Permission, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PermissionService) CreatePermission(ctx context.Context, input PermissionCreateInput) (*domain.Permission, error) {
	// Check name uniqueness
	existing, _ := s.repo.GetByName(ctx, input.Name)
	if existing != nil {
		return nil, NewServiceError(409, "name_conflict", "Permission name already exists")
	}

	// Validate parent exists if pid is provided
	if input.Pid != nil {
		_, err := s.repo.GetByID(ctx, *input.Pid)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, BadRequestErr("parent_not_found", "Parent permission not found")
			}
			return nil, err
		}
	}

	perm := &domain.Permission{
		Pid:          input.Pid,
		Name:         input.Name,
		ResourceType: input.ResourceType,
		ResourcePath: input.ResourcePath,
		HTTPMethod:   &input.HTTPMethod,
		Description:  &input.Description,
		IsActive:     input.IsActive,
	}

	if err := s.repo.Create(ctx, perm); err != nil {
		return nil, err
	}

	return perm, nil
}

func (s *PermissionService) UpdatePermission(ctx context.Context, id int, input PermissionUpdateInput) (*domain.Permission, error) {
	perm, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, HandleRepoErr(err, "permission_not_found", "Permission not found")
	}

	if perm.IsBuiltin {
		return nil, NewServiceError(403, "builtin_immutable", "Cannot modify built-in permission")
	}

	// Check name uniqueness if changed
	if input.Name != "" && input.Name != perm.Name {
		existing, _ := s.repo.GetByName(ctx, input.Name)
		if existing != nil {
			return nil, NewServiceError(409, "name_conflict", "Permission name already exists")
		}
		perm.Name = input.Name
	}

	if input.Pid != nil {
		// Prevent self-referencing
		if *input.Pid == id {
			return nil, NewServiceError(400, "self_parent", "Cannot be parent of itself")
		}
		perm.Pid = input.Pid
	}

	if input.ResourceType != "" {
		perm.ResourceType = input.ResourceType
	}

	if input.ResourcePath != "" {
		perm.ResourcePath = input.ResourcePath
	}

	if input.HTTPMethod != "" {
		perm.HTTPMethod = &input.HTTPMethod
	}

	if input.Description != "" {
		perm.Description = &input.Description
	}

	perm.IsActive = input.IsActive

	if err := s.repo.Update(ctx, perm); err != nil {
		return nil, err
	}

	return perm, nil
}

func (s *PermissionService) DeletePermission(ctx context.Context, id int) error {
	perm, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return HandleRepoErr(err, "permission_not_found", "Permission not found")
	}

	if perm.IsBuiltin {
		return NewServiceError(403, "builtin_immutable", "Cannot delete built-in permission")
	}

	// Check if permission is used by any role
	used, err := s.repo.IsUsedByRoles(ctx, id)
	if err != nil {
		return err
	}
	if used {
		return NewServiceError(409, "permission_assigned", "Permission is assigned to roles")
	}

	return s.repo.Delete(ctx, id)
}
