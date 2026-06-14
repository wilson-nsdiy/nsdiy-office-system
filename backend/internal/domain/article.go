package domain

import "time"

type Article struct {
	ID               int
	ArticleNo        string
	Title            string
	Content          *string
	Summary          *string
	Status           string
	AuthorID         int
	CoverDescription *string
	CoverUrl         *string
	FirstPublishedAt *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type ArticleDetail struct {
	Article
	AuthorName     string
	AuthorNickname *string
}

type ArticleVersion struct {
	ID               int
	ArticleID        int
	VersionNo        int
	Title            string
	Content          *string
	CoverDescription *string
	Summary          *string
	Status           string
	EditorID         *int
	EditReason       *string
	CreatedAt        time.Time
}

type ArticleVersionDetail struct {
	ArticleVersion
	EditorName     *string
	EditorNickname *string
}
