package repository

import (
	"context"
	"testing"
)

func TestTaskRepository_CRUD(t *testing.T) {
	client := setupTestDB(t)

	// Create prerequisite user
	userRepo := NewUserRepository(client)
	user := &User{
		Username:       "dev",
		Email:          "dev@example.com",
		Salt:           "salt",
		HashedPassword: "pwd",
		UserType:       "HUMAN",
		IsActive:       true,
		TokenVersion:   1,
	}
	if err := userRepo.Create(context.Background(), user); err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	// Create prerequisite project
	projRepo := NewProjectRepository(client)
	project := &Project{
		Name:      "Test Project",
		ProjectNo: "PRJ0001",
		Status:    "TODO",
		Priority:  "MEDIUM",
		OwnerID:   user.ID,
	}
	if err := projRepo.Create(context.Background(), project); err != nil {
		t.Fatalf("create project failed: %v", err)
	}

	repo := NewTaskRepository(client)
	desc := "Task description"
	hours := 8.0
	task := &Task{
		ProjectID:      project.ID,
		Title:          "Test Task",
		Description:    &desc,
		Status:         "TODO",
		Priority:       "HIGH",
		AssigneeID:     &user.ID,
		CreatorID:      user.ID,
		EstimatedHours: &hours,
	}

	err := repo.Create(context.Background(), task)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if task.ID == 0 {
		t.Error("expected non-zero ID")
	}

	got, err := repo.GetByID(context.Background(), task.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.Title != "Test Task" {
		t.Errorf("expected Title=Test Task, got %s", got.Title)
	}
	if got.EstimatedHours == nil || *got.EstimatedHours != 8.0 {
		t.Errorf("expected EstimatedHours=8.0, got %v", got.EstimatedHours)
	}

	detail, err := repo.GetDetailByID(context.Background(), task.ID)
	if err != nil {
		t.Fatalf("GetDetailByID failed: %v", err)
	}
	if detail.CreatorName != "dev" {
		t.Errorf("expected CreatorName=dev, got %s", detail.CreatorName)
	}

	// Test List with filters
	tasks, total, err := repo.List(context.Background(), &project.ID, "TODO", "HIGH", &user.ID, 1, 10)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total=1, got %d", total)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}

	// Update
	task.Status = "DONE"
	err = repo.Update(context.Background(), task)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, err := repo.GetByID(context.Background(), task.ID)
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}
	if updated.Status != "DONE" {
		t.Errorf("expected Status=DONE, got %s", updated.Status)
	}

	// Delete
	err = repo.Delete(context.Background(), task.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}
