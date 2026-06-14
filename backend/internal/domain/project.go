package domain

import "time"

type Project struct {
	ID                int
	Name              string
	ProjectNo         string
	Description       *string
	Status            string
	Priority          string
	ExpectedStartDate *time.Time
	ExpectedEndDate   *time.Time
	StartDate         *time.Time
	EndDate           *time.Time
	OwnerID           int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ProjectDetail struct {
	Project
	OwnerNickname *string
}

type ProjectMember struct {
	ID        int
	ProjectID int
	UserID    int
	Role      string
	JoinedAt  time.Time
}

type ProjectMemberDetail struct {
	ProjectMember
	Username string
	Nickname *string
	Email    string
	IsActive bool
}
