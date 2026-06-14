package repository

import (
	"context"
	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/ent/permission"
	"oa-nsdiy/backend/ent/role"
	"oa-nsdiy/backend/ent/roleperm"
	"oa-nsdiy/backend/internal/domain"
)

type PermissionRepository struct {
	client *ent.Client
}

func NewPermissionRepository(client *ent.Client) *PermissionRepository {
	return &PermissionRepository{client: client}
}

// Type alias for backward compatibility
type Permission = domain.Permission

func (r *PermissionRepository) GetByID(ctx context.Context, id int) (*Permission, error) {
	e, err := r.client.Permission.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toPermissionEntity(e), nil
}

func (r *PermissionRepository) GetByRoleID(ctx context.Context, roleID int) ([]*Permission, error) {
	es, err := r.client.Role.Query().
		Where(role.IDEQ(roleID)).
		QueryRolePerms().
		QueryPermission().
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toPermissionEntities(es), nil
}

func (r *PermissionRepository) GetByName(ctx context.Context, name string) (*Permission, error) {
	e, err := r.client.Permission.Query().Where(permission.NameEQ(name)).First(ctx)
	if err != nil {
		return nil, err
	}
	return toPermissionEntity(e), nil
}

func (r *PermissionRepository) Create(ctx context.Context, perm *Permission) error {
	e, err := r.client.Permission.Create().
		SetNillablePid(perm.Pid).
		SetName(perm.Name).
		SetResourceType(perm.ResourceType).
		SetResourcePath(perm.ResourcePath).
		SetNillableHTTPMethod(perm.HTTPMethod).
		SetNillableDescription(perm.Description).
		SetIsActive(perm.IsActive).
		SetIsBuiltin(perm.IsBuiltin).
		Save(ctx)
	if err != nil {
		return err
	}
	perm.ID = e.ID
	return nil
}

func (r *PermissionRepository) Update(ctx context.Context, perm *Permission) error {
	_, err := r.client.Permission.UpdateOneID(perm.ID).
		SetNillablePid(perm.Pid).
		SetName(perm.Name).
		SetResourceType(perm.ResourceType).
		SetResourcePath(perm.ResourcePath).
		SetNillableHTTPMethod(perm.HTTPMethod).
		SetNillableDescription(perm.Description).
		SetIsActive(perm.IsActive).
		Save(ctx)
	return err
}

func (r *PermissionRepository) Delete(ctx context.Context, id int) error {
	return r.client.Permission.DeleteOneID(id).Exec(ctx)
}

func (r *PermissionRepository) ListActive(ctx context.Context) ([]*Permission, error) {
	es, err := r.client.Permission.Query().Where(permission.IsActive(true)).All(ctx)
	if err != nil {
		return nil, err
	}
	return toPermissionEntities(es), nil
}

func (r *PermissionRepository) ListAll(ctx context.Context) ([]*Permission, error) {
	es, err := r.client.Permission.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return toPermissionEntities(es), nil
}

func (r *PermissionRepository) Search(ctx context.Context, resourceType, keyword string) ([]*Permission, error) {
	q := r.client.Permission.Query().Where(permission.IsActive(true))

	if resourceType != "" {
		q = q.Where(permission.ResourceTypeEQ(resourceType))
	}

	if keyword != "" {
		q = q.Where(permission.Or(
			permission.NameContainsFold(keyword),
			permission.ResourcePathContainsFold(keyword),
		))
	}

	es, err := q.All(ctx)
	if err != nil {
		return nil, err
	}
	return toPermissionEntities(es), nil
}

func (r *PermissionRepository) IsUsedByRoles(ctx context.Context, permID int) (bool, error) {
	count, err := r.client.RolePerm.Query().Where(roleperm.PermissionIDEQ(permID)).Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func toPermissionEntity(e *ent.Permission) *Permission {
	if e == nil {
		return nil
	}
	return &Permission{
		ID:           e.ID,
		Pid:          intPtr(e.Pid),
		Name:         e.Name,
		ResourceType: e.ResourceType,
		ResourcePath: e.ResourcePath,
		HTTPMethod:   stringPtr(e.HTTPMethod),
		Description:  stringPtr(e.Description),
		IsActive:     e.IsActive,
		IsBuiltin:    e.IsBuiltin,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func toPermissionEntities(es []*ent.Permission) []*Permission {
	result := make([]*Permission, len(es))
	for i, e := range es {
		result[i] = toPermissionEntity(e)
	}
	return result
}
