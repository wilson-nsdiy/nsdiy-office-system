package domain

import "time"

type Permission struct {
	ID           int
	Pid          *int
	Name         string
	ResourceType string
	ResourcePath string
	HTTPMethod   *string
	Description  *string
	IsActive     bool
	IsBuiltin    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
