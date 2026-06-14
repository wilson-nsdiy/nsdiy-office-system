package service

import (
	"context"

	"oa-nsdiy/backend/internal/domain"
)

// RoleRepository defines the interface for role data access required by RoleService.
type RoleRepository interface {
	ListActive(ctx context.Context) ([]*domain.Role, error)
	GetByID(ctx context.Context, id int) (*domain.Role, error)
	GetByUserID(ctx context.Context, userID int) (*domain.Role, error)
	GetByName(ctx context.Context, name string) (*domain.Role, error)
	GetByCode(ctx context.Context, code string) (*domain.Role, error)
	Create(ctx context.Context, role *domain.Role) error
	Update(ctx context.Context, role *domain.Role) error
	Delete(ctx context.Context, id int) error
	IsUsedByUsers(ctx context.Context, roleID int) (bool, error)
	GetPermissions(ctx context.Context, roleID int) ([]*domain.Permission, error)
	UpdatePermissions(ctx context.Context, roleID int, permissionIDs []int) error
}

// PermissionValidator defines the minimal permission interface used by RoleService for validation.
type PermissionValidator interface {
	GetByID(ctx context.Context, id int) (*domain.Permission, error)
	GetByRoleID(ctx context.Context, roleID int) ([]*domain.Permission, error)
}

type RoleService struct {
	repo     RoleRepository
	permRepo PermissionValidator
}

func NewRoleService(repo RoleRepository, permRepo PermissionValidator) *RoleService {
	return &RoleService{repo: repo, permRepo: permRepo}
}

type RoleCreateInput struct {
	Name        string
	Code        string
	Description string
	RoleType    string
	IsActive    bool
}

type RoleUpdateInput struct {
	Name        string
	Description string
	RoleType    string
	IsActive    bool
}

func (s *RoleService) GetRoles(ctx context.Context) ([]*domain.Role, error) {
	return s.repo.ListActive(ctx)
}

func (s *RoleService) GetRole(ctx context.Context, id int) (*domain.Role, error) {
	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, HandleRepoErr(err, "role_not_found", "Role not found")
	}
	return role, nil
}

func (s *RoleService) CreateRole(ctx context.Context, input RoleCreateInput) (*domain.Role, error) {
	// Check name uniqueness
	existing, _ := s.repo.GetByName(ctx, input.Name)
	if existing != nil {
		return nil, NewServiceError(409, "name_conflict", "Role name already exists")
	}

	// Check code uniqueness
	existing, _ = s.repo.GetByCode(ctx, input.Code)
	if existing != nil {
		return nil, NewServiceError(409, "code_conflict", "Role code already exists")
	}

	role := &domain.Role{
		Name:        input.Name,
		Code:        input.Code,
		Description: &input.Description,
		IsActive:    input.IsActive,
		RoleType:    &input.RoleType,
	}

	if err := s.repo.Create(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) UpdateRole(ctx context.Context, id int, input RoleUpdateInput) (*domain.Role, error) {
	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, HandleRepoErr(err, "role_not_found", "Role not found")
	}

	if role.IsBuiltin {
		return nil, NewServiceError(403, "builtin_immutable", "Cannot modify built-in role")
	}

	// Check name uniqueness if changed
	if input.Name != "" && input.Name != role.Name {
		existing, _ := s.repo.GetByName(ctx, input.Name)
		if existing != nil {
			return nil, NewServiceError(409, "name_conflict", "Role name already exists")
		}
		role.Name = input.Name
	}

	if input.Description != "" {
		role.Description = &input.Description
	}

	if input.RoleType != "" {
		role.RoleType = &input.RoleType
	}

	role.IsActive = input.IsActive

	if err := s.repo.Update(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) DeleteRole(ctx context.Context, id int) error {
	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return HandleRepoErr(err, "role_not_found", "Role not found")
	}

	if role.IsBuiltin {
		return NewServiceError(403, "builtin_immutable", "Cannot delete built-in role")
	}

	// Check if role is used by any user
	used, err := s.repo.IsUsedByUsers(ctx, id)
	if err != nil {
		return err
	}
	if used {
		return NewServiceError(409, "role_assigned", "Role is assigned to users")
	}

	return s.repo.Delete(ctx, id)
}

func (s *RoleService) GetRolePermissions(ctx context.Context, roleID int) ([]*domain.Permission, error) {
	return s.repo.GetPermissions(ctx, roleID)
}

func (s *RoleService) UpdateRolePermissions(ctx context.Context, roleID int, permissionIDs []int) error {
	// Validate all permission IDs exist
	for _, permID := range permissionIDs {
		_, err := s.permRepo.GetByID(ctx, permID)
		if err != nil {
			return HandleRepoErr(err, "invalid_permission_id", "Invalid permission ID")
		}
	}

	return s.repo.UpdatePermissions(ctx, roleID, permissionIDs)
}

// HasPermission checks if a user has ALL of the specified permissions.
func (s *RoleService) HasPermission(ctx context.Context, userID int, permissionCodes ...string) (bool, error) {
	// Get user's role
	role, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	if role == nil {
		return false, nil
	}
	// Get role's permissions
	perms, err := s.permRepo.GetByRoleID(ctx, role.ID)
	if err != nil {
		return false, err
	}
	// Build a set of permission codes the user has
	permSet := make(map[string]bool, len(perms))
	for _, p := range perms {
		permSet[p.Name] = true
	}
	// Check all required permissions exist
	for _, code := range permissionCodes {
		if !permSet[code] {
			return false, nil
		}
	}
	return true, nil
}
