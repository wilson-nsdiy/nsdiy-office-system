package repository

import (
	"context"
	"testing"
)

func TestUserRepository_CRUD(t *testing.T) {
	client := setupTestDB(t)
	repo := NewUserRepository(client)

	// Test Create
	nickname := "Test User"
	user := &User{
		Username:       "testuser",
		Email:          "test@example.com",
		Nickname:       &nickname,
		Salt:           "randomsalt",
		HashedPassword: "hashedpwd",
		UserType:       "HUMAN",
		IsActive:       true,
		TokenVersion:   1,
	}
	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if user.ID == 0 {
		t.Error("expected non-zero ID after create")
	}

	// Test GetByID
	got, err := repo.GetByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.Username != "testuser" {
		t.Errorf("expected Username=testuser, got %s", got.Username)
	}
	if got.Email != "test@example.com" {
		t.Errorf("expected Email=test@example.com, got %s", got.Email)
	}
	if got.IsActive != true {
		t.Error("expected IsActive=true")
	}

	// Test GetByUsername
	got2, err := repo.GetByUsername(context.Background(), "testuser")
	if err != nil {
		t.Fatalf("GetByUsername failed: %v", err)
	}
	if got2.ID != user.ID {
		t.Errorf("expected ID=%d, got %d", user.ID, got2.ID)
	}

	// Test GetByEmail
	got3, err := repo.GetByEmail(context.Background(), "test@example.com")
	if err != nil {
		t.Fatalf("GetByEmail failed: %v", err)
	}
	if got3.ID != user.ID {
		t.Errorf("expected ID=%d, got %d", user.ID, got3.ID)
	}

	// Test GetActiveByID
	active, err := repo.GetActiveByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("GetActiveByID failed: %v", err)
	}
	if !active.IsActive {
		t.Error("expected user to be active")
	}

	// Test Update
	user.IsActive = false
	err = repo.Update(context.Background(), user)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	updated, err := repo.GetByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}
	if updated.IsActive {
		t.Error("expected IsActive=false after update")
	}

	// Test UpdatePassword
	err = repo.UpdatePassword(context.Background(), user.ID, "newsalt", "newhash")
	if err != nil {
		t.Fatalf("UpdatePassword failed: %v", err)
	}
	pwdUser, err := repo.GetByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("GetByID after password update failed: %v", err)
	}
	if pwdUser.Salt != "newsalt" {
		t.Errorf("expected Salt=newsalt, got %s", pwdUser.Salt)
	}
	if pwdUser.TokenVersion != 2 {
		t.Errorf("expected TokenVersion=2, got %d", pwdUser.TokenVersion)
	}

	// Test List
	users, total, err := repo.List(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total=1, got %d", total)
	}
	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}

	// Test Delete
	err = repo.Delete(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err = repo.GetByID(context.Background(), user.ID)
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestUserRepository_NotFound(t *testing.T) {
	client := setupTestDB(t)
	repo := NewUserRepository(client)

	_, err := repo.GetByID(context.Background(), 999)
	if err == nil {
		t.Error("expected error for non-existent user")
	}

	_, err = repo.GetByUsername(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error for non-existent username")
	}
}

func TestUserRepository_Duplicate(t *testing.T) {
	client := setupTestDB(t)
	repo := NewUserRepository(client)

	user1 := &User{
		Username:       "dupuser",
		Email:          "dup@example.com",
		Salt:           "salt",
		HashedPassword: "pwd",
		UserType:       "HUMAN",
		IsActive:       true,
		TokenVersion:   1,
	}
	err := repo.Create(context.Background(), user1)
	if err != nil {
		t.Fatalf("first Create failed: %v", err)
	}

	user2 := &User{
		Username:       "dupuser",
		Email:          "dup2@example.com",
		Salt:           "salt",
		HashedPassword: "pwd",
		UserType:       "HUMAN",
		IsActive:       true,
		TokenVersion:   1,
	}
	err = repo.Create(context.Background(), user2)
	if err == nil {
		t.Error("expected error for duplicate username")
	}
}
