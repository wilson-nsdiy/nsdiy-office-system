package repository

import (
	"context"
	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/ent/role"
	"oa-nsdiy/backend/ent/roleperm"
	"oa-nsdiy/backend/ent/user"
	"oa-nsdiy/backend/internal/domain"
)

type RoleRepository struct {
	client *ent.Client
}

func NewRoleRepository(client *ent.Client) *RoleRepository {
	return &RoleRepository{client: client}
}

// Type alias for backward compatibility
type Role = domain.Role

func (r *RoleRepository) GetByID(ctx context.Context, id int) (*Role, error) {
	e, err := r.client.Role.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toRoleEntity(e), nil
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (*Role, error) {
	e, err := r.client.Role.Query().Where(role.NameEQ(name)).First(ctx)
	if err != nil {
		return nil, err
	}
	return toRoleEntity(e), nil
}

func (r *RoleRepository) GetByCode(ctx context.Context, code string) (*Role, error) {
	e, err := r.client.Role.Query().Where(role.CodeEQ(code)).First(ctx)
	if err != nil {
		return nil, err
	}
	return toRoleEntity(e), nil
}

func (r *RoleRepository) Create(ctx context.Context, role *Role) error {
	e, err := r.client.Role.Create().
		SetName(role.Name).
		SetCode(role.Code).
		SetNillableDescription(role.Description).
		SetIsActive(role.IsActive).
		SetIsDefault(role.IsDefault).
		SetIsBuiltin(role.IsBuiltin).
		SetNillableRoleType(role.RoleType).
		Save(ctx)
	if err != nil {
		return err
	}
	role.ID = e.ID
	return nil
}

func (r *RoleRepository) Update(ctx context.Context, role *Role) error {
	_, err := r.client.Role.UpdateOneID(role.ID).
		SetName(role.Name).
		SetNillableDescription(role.Description).
		SetNillableRoleType(role.RoleType).
		SetIsActive(role.IsActive).
		Save(ctx)
	return err
}

func (r *RoleRepository) Delete(ctx context.Context, id int) error {
	return r.client.Role.DeleteOneID(id).Exec(ctx)
}

func (r *RoleRepository) GetByUserID(ctx context.Context, userID int) (*Role, error) {
	// Get user to obtain RoleID
	u, err := r.client.User.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u.RoleID == 0 {
		return nil, nil
	}
	e, err := r.client.Role.Get(ctx, u.RoleID)
	if err != nil {
		return nil, err
	}
	return toRoleEntity(e), nil
}

func (r *RoleRepository) ListActive(ctx context.Context) ([]*Role, error) {
	es, err := r.client.Role.Query().Where(role.IsActive(true)).All(ctx)
	if err != nil {
		return nil, err
	}
	return toRoleEntities(es), nil
}

func (r *RoleRepository) IsUsedByUsers(ctx context.Context, roleID int) (bool, error) {
	count, err := r.client.User.Query().Where(user.RoleIDEQ(roleID)).Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RoleRepository) GetPermissions(ctx context.Context, roleID int) ([]*Permission, error) {
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

func (r *RoleRepository) UpdatePermissions(ctx context.Context, roleID int, permissionIDs []int) error {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing permissions
	_, err = tx.RolePerm.Delete().Where(roleperm.RoleIDEQ(roleID)).Exec(ctx)
	if err != nil {
		return err
	}

	// Insert new permissions
	for _, pid := range permissionIDs {
		_, err = tx.RolePerm.Create().SetRoleID(roleID).SetPermissionID(pid).Save(ctx)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func toRoleEntity(e *ent.Role) *Role {
	if e == nil {
		return nil
	}
	return &Role{
		ID:          e.ID,
		Name:        e.Name,
		Code:        e.Code,
		Description: stringPtr(e.Description),
		IsActive:    e.IsActive,
		IsDefault:   e.IsDefault,
		IsBuiltin:   e.IsBuiltin,
		RoleType:    stringPtr(e.RoleType),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toRoleEntities(es []*ent.Role) []*Role {
	result := make([]*Role, len(es))
	for i, e := range es {
		result[i] = toRoleEntity(e)
	}
	return result
}
