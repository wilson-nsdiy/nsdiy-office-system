package domain

import "time"

// User entity - keeps pointer types for service layer compatibility.
type User struct {
	ID                        int
	Username                  string
	Email                     string
	Nickname                  *string
	Salt                      string
	HashedPassword            string
	RoleID                    *int
	UserType                  string
	IsActive                  bool
	TokenVersion              int
	VerificationCode          *string
	VerificationCodeExpiresAt *time.Time
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}
