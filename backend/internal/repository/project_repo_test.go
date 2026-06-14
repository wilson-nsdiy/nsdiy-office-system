package repository

import (
	"context"
	"testing"
)

func TestProjectRepository_CRUD(t *testing.T) {
	client := setupTestDB(t)
	// Create a user first (owner)
	userRepo := NewUserRepository(client)
	user := &User{
		Username:       "owner",
		Email:          "owner@example.com",
		Salt:           "salt",
		HashedPassword: "pwd",
		UserType:       "HUMAN",
		IsActive:       true,
		TokenVersion:   1,
	}
	if err := userRepo.Create(context.Background(), user); err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	repo := NewProjectRepository(client)
	desc := "A test project"
	project := &Project{
		Name:        "Test Project",
		ProjectNo:   "PRJ0001",
		Description: &desc,
		Status:      "TODO",
		Priority:    "HIGH",
		OwnerID:     user.ID,
	}

	err := repo.Create(context.Background(), project)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if project.ID == 0 {
		t.Error("expected non-zero ID")
	}

	got, err := repo.GetByID(context.Background(), project.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.Name != "Test Project" {
		t.Errorf("expected Name=Test Project, got %s", got.Name)
	}

	got2, err := repo.GetByProjectNo(context.Background(), "PRJ0001")
	if err != nil {
		t.Fatalf("GetByProjectNo failed: %v", err)
	}
	if got2.ID != project.ID {
		t.Errorf("expected ID=%d, got %d", project.ID, got2.ID)
	}

	detail, err := repo.GetDetailByID(context.Background(), project.ID)
	if err != nil {
		t.Fatalf("GetDetailByID failed: %v", err)
	}
	if detail.OwnerNickname == nil {
		t.Log("owner nickname is nil (no nickname set)")
	}

	// Test ListByUserID
	projects, total, err := repo.ListByUserID(context.Background(), user.ID, "", 1, 10)
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total=1, got %d", total)
	}
	if len(projects) != 1 {
		t.Errorf("expected 1 project, got %d", len(projects))
	}

	// Test member operations
	member := &ProjectMember{
		ProjectID: project.ID,
		UserID:    user.ID,
		Role:      "MEMBER",
	}
	err = repo.AddMember(context.Background(), member)
	if err != nil {
		t.Fatalf("AddMember failed: %v", err)
	}

	isMember, err := repo.IsMember(context.Background(), project.ID, user.ID)
	if err != nil {
		t.Fatalf("IsMember failed: %v", err)
	}
	if !isMember {
		t.Error("expected user to be a member")
	}

	err = repo.UpdateMemberRole(context.Background(), project.ID, user.ID, "OWNER")
	if err != nil {
		t.Fatalf("UpdateMemberRole failed: %v", err)
	}

	members, err := repo.GetMembers(context.Background(), project.ID)
	if err != nil {
		t.Fatalf("GetMembers failed: %v", err)
	}
	if len(members) != 1 {
		t.Errorf("expected 1 member, got %d", len(members))
	}

	err = repo.RemoveMember(context.Background(), project.ID, user.ID)
	if err != nil {
		t.Fatalf("RemoveMember failed: %v", err)
	}

	err = repo.Delete(context.Background(), project.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestProjectRepository_GenerateNo(t *testing.T) {
	client := setupTestDB(t)
	repo := NewProjectRepository(client)

	no, err := repo.GenerateUniqueProjectNo(context.Background())
	if err != nil {
		t.Fatalf("GenerateUniqueProjectNo failed: %v", err)
	}
	if no != "PRJ0001" {
		t.Errorf("expected PRJ0001, got %s", no)
	}
}
