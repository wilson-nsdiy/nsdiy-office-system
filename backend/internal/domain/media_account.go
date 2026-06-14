package domain

import "time"

type MediaAccount struct {
	ID             int
	Name           string
	Platform       string
	AccountID      string
	Avatar         *string
	Status         string
	AccessToken    *string
	RefreshToken   *string
	TokenExpiresAt *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
