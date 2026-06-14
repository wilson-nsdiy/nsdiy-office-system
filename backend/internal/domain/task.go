package domain

import "time"

type Task struct {
	ID               int
	ProjectID        int
	ParentID         *int
	Title            string
	Description      *string
	Status           string
	Priority         string
	AssigneeID       *int
	CreatorID        int
	PlannedStartDate *time.Time
	PlannedEndDate   *time.Time
	ActualStartTime  *time.Time
	ActualEndTime    *time.Time
	EstimatedHours   *float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type TaskDetail struct {
	Task
	AssigneeName     *string
	AssigneeNickname *string
	CreatorName      string
	CreatorNickname  *string
}
