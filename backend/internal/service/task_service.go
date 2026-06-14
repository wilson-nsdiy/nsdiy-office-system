package service

import (
	"context"
	"time"

	"oa-nsdiy/backend/internal/domain"
)

// TaskRepository defines the interface for task data access required by TaskService.
type TaskRepository interface {
	List(ctx context.Context, projectID *int, status, priority string, assigneeID *int, page, pageSize int) ([]*domain.TaskDetail, int64, error)
	GetDetailByID(ctx context.Context, id int) (*domain.TaskDetail, error)
	GetByID(ctx context.Context, id int) (*domain.Task, error)
	Create(ctx context.Context, task *domain.Task) error
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id int) error
}

// ProjectValidator defines the minimal project interface used by TaskService for validation.
type ProjectValidator interface {
	GetByID(ctx context.Context, id int) (*domain.Project, error)
}

type TaskService struct {
	repo        TaskRepository
	projectRepo ProjectValidator
}

func NewTaskService(repo TaskRepository, projectRepo ProjectValidator) *TaskService {
	return &TaskService{repo: repo, projectRepo: projectRepo}
}

type TaskCreateInput struct {
	Title            string
	Description      string
	Status           string
	Priority         string
	AssigneeID       *int
	ParentID         *int
	PlannedStartDate string
	PlannedEndDate   string
	EstimatedHours   float64
}

type TaskUpdateInput struct {
	Title            string
	Description      string
	Status           string
	Priority         string
	AssigneeID       *int
	ParentID         *int
	PlannedStartDate string
	PlannedEndDate   string
	EstimatedHours   float64
}

type TaskListResult struct {
	Items []*domain.TaskDetail
	Total int64
}

func (s *TaskService) GetTaskList(ctx context.Context, projectID *int, status, priority string, assigneeID *int, page, pageSize int) (*TaskListResult, error) {
	items, total, err := s.repo.List(ctx, projectID, status, priority, assigneeID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &TaskListResult{Items: items, Total: total}, nil
}

func (s *TaskService) GetTask(ctx context.Context, id int) (*domain.TaskDetail, error) {
	return s.repo.GetDetailByID(ctx, id)
}

func (s *TaskService) CreateTask(ctx context.Context, input TaskCreateInput, projectID, creatorID int) (*domain.Task, error) {
	if input.Title == "" {
		return nil, NewServiceError(400, "title_required", "Title is required")
	}

	// Validate project exists
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, NewServiceError(400, "project_not_found", "Project not found")
	}

	task := &domain.Task{
		ProjectID:   projectID,
		ParentID:    input.ParentID,
		Title:       input.Title,
		Description: &input.Description,
		Status:      "TODO",
		Priority:    "MEDIUM",
		AssigneeID:  input.AssigneeID,
		CreatorID:   creatorID,
	}

	if input.Status != "" {
		task.Status = input.Status
	}
	if input.Priority != "" {
		task.Priority = input.Priority
	}
	if input.EstimatedHours > 0 {
		task.EstimatedHours = &input.EstimatedHours
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id int, input TaskUpdateInput) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, HandleRepoErr(err, "task_not_found", "Task not found")
	}

	if input.Title != "" {
		task.Title = input.Title
	}
	if input.Description != "" {
		task.Description = &input.Description
	}
	if input.Status != "" {
		// Auto-set timestamps on status change
		if input.Status == "IN_PROGRESS" && task.ActualStartTime == nil {
			now := time.Now()
			task.ActualStartTime = &now
		}
		if input.Status == "DONE" && task.ActualEndTime == nil {
			now := time.Now()
			task.ActualEndTime = &now
		}
		task.Status = input.Status
	}
	if input.Priority != "" {
		task.Priority = input.Priority
	}
	if input.AssigneeID != nil {
		task.AssigneeID = input.AssigneeID
	}
	if input.ParentID != nil {
		task.ParentID = input.ParentID
	}
	if input.EstimatedHours > 0 {
		task.EstimatedHours = &input.EstimatedHours
	}

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return HandleRepoErr(err, "task_not_found", "Task not found")
	}

	return s.repo.Delete(ctx, id)
}
