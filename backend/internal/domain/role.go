package domain

import "time"

type Role struct {
	ID          int
	Name        string
	Code        string
	Description *string
	IsActive    bool
	IsDefault   bool
	IsBuiltin   bool
	RoleType    *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
