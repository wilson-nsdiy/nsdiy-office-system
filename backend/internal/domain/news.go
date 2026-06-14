package domain

import "time"

type News struct {
	ID        int
	GroupID   int
	Title     string
	Content   *string
	CreatorID int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewsDetail struct {
	News
	GroupName       string
	CreatorNickname *string
}
