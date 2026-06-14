package domain

import "time"

type NewsGroup struct {
	ID          int
	Name        string
	Description *string
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
