package repository

import (
	"context"

	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/ent/task"
	"oa-nsdiy/backend/internal/domain"
)

type TaskRepository struct {
	client *ent.Client
}

func NewTaskRepository(client *ent.Client) *TaskRepository {
	return &TaskRepository{client: client}
}

// Type aliases for backward compatibility
type Task = domain.Task
type TaskDetail = domain.TaskDetail

func (r *TaskRepository) GetByID(ctx context.Context, id int) (*Task, error) {
	e, err := r.client.Task.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toTask(e), nil
}

func (r *TaskRepository) GetDetailByID(ctx context.Context, id int) (*TaskDetail, error) {
	e, err := r.client.Task.Query().
		Where(task.IDEQ(id)).
		WithAssignee().
		WithCreator().
		First(ctx)
	if err != nil {
		return nil, err
	}
	detail := &TaskDetail{Task: *toTask(e)}
	if e.Edges.Assignee != nil {
		detail.AssigneeName = stringPtr(e.Edges.Assignee.Username)
		detail.AssigneeNickname = stringPtr(e.Edges.Assignee.Nickname)
	}
	if e.Edges.Creator != nil {
		detail.CreatorName = e.Edges.Creator.Username
		detail.CreatorNickname = stringPtr(e.Edges.Creator.Nickname)
	}
	return detail, nil
}

func (r *TaskRepository) Create(ctx context.Context, t *Task) error {
	e, err := r.client.Task.Create().
		SetProjectID(t.ProjectID).
		SetNillableParentID(t.ParentID).
		SetTitle(t.Title).
		SetNillableDescription(t.Description).
		SetStatus(t.Status).
		SetPriority(t.Priority).
		SetNillableAssigneeID(t.AssigneeID).
		SetCreatorID(t.CreatorID).
		SetNillablePlannedStartDate(t.PlannedStartDate).
		SetNillablePlannedEndDate(t.PlannedEndDate).
		SetNillableEstimatedHours(t.EstimatedHours).
		Save(ctx)
	if err != nil {
		return err
	}
	t.ID = e.ID
	return nil
}

func (r *TaskRepository) Update(ctx context.Context, t *Task) error {
	_, err := r.client.Task.UpdateOneID(t.ID).
		SetNillableParentID(t.ParentID).
		SetTitle(t.Title).
		SetNillableDescription(t.Description).
		SetStatus(t.Status).
		SetPriority(t.Priority).
		SetNillableAssigneeID(t.AssigneeID).
		SetNillablePlannedStartDate(t.PlannedStartDate).
		SetNillablePlannedEndDate(t.PlannedEndDate).
		SetNillableActualStartTime(t.ActualStartTime).
		SetNillableActualEndTime(t.ActualEndTime).
		SetNillableEstimatedHours(t.EstimatedHours).
		Save(ctx)
	return err
}

func (r *TaskRepository) Delete(ctx context.Context, id int) error {
	return r.client.Task.DeleteOneID(id).Exec(ctx)
}

func (r *TaskRepository) List(ctx context.Context, projectID *int, status, priority string, assigneeID *int, page, pageSize int) ([]*TaskDetail, int64, error) {

	q := r.client.Task.Query()
	if projectID != nil {
		q.Where(task.ProjectIDEQ(*projectID))
	}
	if status != "" {
		q.Where(task.StatusEQ(status))
	}
	if priority != "" {
		q.Where(task.PriorityEQ(priority))
	}
	if assigneeID != nil {
		q.Where(task.AssigneeIDEQ(*assigneeID))
	}

	count, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	total := int64(count)

	offset := (page - 1) * pageSize
	entities, err := q.WithAssignee().
		WithCreator().
		Order(ent.Desc(task.FieldID)).
		Limit(pageSize).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	var tasks []*TaskDetail
	for _, e := range entities {
		detail := &TaskDetail{Task: *toTask(e)}
		if e.Edges.Assignee != nil {
			detail.AssigneeName = stringPtr(e.Edges.Assignee.Username)
			detail.AssigneeNickname = stringPtr(e.Edges.Assignee.Nickname)
		}
		if e.Edges.Creator != nil {
			detail.CreatorName = e.Edges.Creator.Username
			detail.CreatorNickname = stringPtr(e.Edges.Creator.Nickname)
		}
		tasks = append(tasks, detail)
	}

	return tasks, total, nil
}

// --- Converters ---

func toTask(e *ent.Task) *Task {
	return &Task{
		ID:               e.ID,
		ProjectID:        e.ProjectID,
		ParentID:         e.ParentID,
		Title:            e.Title,
		Description:      stringPtr(e.Description),
		Status:           e.Status,
		Priority:         e.Priority,
		AssigneeID:       e.AssigneeID,
		CreatorID:        e.CreatorID,
		PlannedStartDate: e.PlannedStartDate,
		PlannedEndDate:   e.PlannedEndDate,
		ActualStartTime:  e.ActualStartTime,
		ActualEndTime:    e.ActualEndTime,
		EstimatedHours:   e.EstimatedHours,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}
