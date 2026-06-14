package repository

import (
	"context"
	"testing"
)

func TestRoleRepository_CRUD(t *testing.T) {
	client := setupTestDB(t)
	repo := NewRoleRepository(client)

	roleType := "SYSTEM"
	role := &Role{
		Name:        "Admin",
		Code:        "admin",
		Description: strPtr("Administrator role"),
		IsActive:    true,
		IsDefault:   false,
		IsBuiltin:   true,
		RoleType:    &roleType,
	}

	err := repo.Create(context.Background(), role)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if role.ID == 0 {
		t.Error("expected non-zero ID")
	}

	got, err := repo.GetByID(context.Background(), role.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.Name != "Admin" {
		t.Errorf("expected Name=Admin, got %s", got.Name)
	}

	got2, err := repo.GetByName(context.Background(), "Admin")
	if err != nil {
		t.Fatalf("GetByName failed: %v", err)
	}
	if got2.ID != role.ID {
		t.Errorf("expected ID=%d, got %d", role.ID, got2.ID)
	}

	roles, err := repo.ListActive(context.Background())
	if err != nil {
		t.Fatalf("ListActive failed: %v", err)
	}
	if len(roles) != 1 {
		t.Errorf("expected 1 active role, got %d", len(roles))
	}

	role.IsActive = false
	err = repo.Update(context.Background(), role)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	used, err := repo.IsUsedByUsers(context.Background(), role.ID)
	if err != nil {
		t.Fatalf("IsUsedByUsers failed: %v", err)
	}
	if used {
		t.Error("expected role not to be used by any users")
	}

	err = repo.Delete(context.Background(), role.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.GetByID(context.Background(), role.ID)
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestRoleRepository_UniqueName(t *testing.T) {
	client := setupTestDB(t)
	repo := NewRoleRepository(client)

	r1 := &Role{
		Name:        "Tester",
		Code:        "tester",
		Description: strPtr("Test role 1"),
		IsActive:    true,
	}
	err := repo.Create(context.Background(), r1)
	if err != nil {
		t.Fatalf("first Create failed: %v", err)
	}

	r2 := &Role{
		Name:        "Tester",
		Code:        "tester2",
		Description: strPtr("Test role 2"),
		IsActive:    true,
	}
	err = repo.Create(context.Background(), r2)
	if err == nil {
		t.Error("expected duplicate name error")
	}
}
