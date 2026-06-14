package domain

import "time"

type ApiToken struct {
	ID          int
	UserID      int
	Name        string
	TokenHash   string
	TokenPrefix string
	Status      string
	ExpiresAt   *time.Time
	LastUsedAt  *time.Time
	UsageCount  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
