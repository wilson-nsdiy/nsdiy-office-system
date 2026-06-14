package domain

import "time"

type MediaContent struct {
	ID          int
	Title       string
	Content     *string
	CoverImage  *string
	Platform    string
	AccountID   *int
	Status      string
	Views       int
	Likes       int
	Comments    int
	Shares      int
	PublishTime *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type MediaContentDetail struct {
	MediaContent
	AccountName *string
}

type MediaContentVersion struct {
	ID         int
	ContentID  int
	VersionNo  int
	Title      string
	Content    *string
	CoverImage *string
	Status     string
	EditorID   *int
	EditReason *string
	CreatedAt  time.Time
}

type MediaContentVersionDetail struct {
	MediaContentVersion
	EditorName     *string
	EditorNickname *string
}
