//go:build unit

package testutil

import (
	"fmt"
	"time"

	"oa-nsdiy/backend/internal/domain"
)

// =============================================================================
// Functional option pattern for test data factories
// =============================================================================

type UserOption func(*domain.User)

func WithUserID(id int) UserOption {
	return func(u *domain.User) { u.ID = id }
}

func WithUsername(username string) UserOption {
	return func(u *domain.User) { u.Username = username }
}

func WithEmail(email string) UserOption {
	return func(u *domain.User) { u.Email = email }
}

func WithNickname(nickname string) UserOption {
	return func(u *domain.User) { u.Nickname = &nickname }
}

func WithUserType(userType string) UserOption {
	return func(u *domain.User) { u.UserType = userType }
}

func WithActive(active bool) UserOption {
	return func(u *domain.User) { u.IsActive = active }
}

func WithTokenVersion(v int) UserOption {
	return func(u *domain.User) { u.TokenVersion = v }
}

func WithSalt(salt string) UserOption {
	return func(u *domain.User) { u.Salt = salt }
}

func WithHashedPassword(hash string) UserOption {
	return func(u *domain.User) { u.HashedPassword = hash }
}

func NewTestUser(opts ...UserOption) *domain.User {
	now := time.Now()
	u := &domain.User{
		ID:             1,
		Username:       "testuser",
		Email:          "test@example.com",
		Nickname:       strPtr("Test User"),
		Salt:           "testsalt",
		HashedPassword: "testhash",
		UserType:       "HUMAN",
		IsActive:       true,
		TokenVersion:   1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

// =============================================================================
// Role fixtures
// =============================================================================

type RoleOption func(*domain.Role)

func WithRoleID(id int) RoleOption {
	return func(r *domain.Role) { r.ID = id }
}

func WithRoleName(name string) RoleOption {
	return func(r *domain.Role) { r.Name = name }
}

func WithRoleCode(code string) RoleOption {
	return func(r *domain.Role) { r.Code = code }
}

func WithRoleActive(active bool) RoleOption {
	return func(r *domain.Role) { r.IsActive = active }
}

func NewTestRole(opts ...RoleOption) *domain.Role {
	now := time.Now()
	r := &domain.Role{
		ID:        1,
		Name:      "testrole",
		Code:      "TEST",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// =============================================================================
// Permission fixtures
// =============================================================================

type PermissionOption func(*domain.Permission)

func WithPermID(id int) PermissionOption {
	return func(p *domain.Permission) { p.ID = id }
}

func WithPermName(name string) PermissionOption {
	return func(p *domain.Permission) { p.Name = name }
}

func WithPermCode(code string) PermissionOption {
	return func(p *domain.Permission) { p.ResourceType = code }
}

func WithResourceType(rt string) PermissionOption {
	return func(p *domain.Permission) { p.ResourceType = rt }
}

func NewTestPermission(opts ...PermissionOption) *domain.Permission {
	now := time.Now()
	p := &domain.Permission{
		ID:           1,
		Name:         "testperm",
		ResourceType: "test",
		ResourcePath: "/test",
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// =============================================================================
// Project fixtures
// =============================================================================

type ProjectOption func(*domain.Project)

func WithProjectID(id int) ProjectOption {
	return func(p *domain.Project) { p.ID = id }
}

func WithProjectName(name string) ProjectOption {
	return func(p *domain.Project) { p.Name = name }
}

func WithProjectNo(no string) ProjectOption {
	return func(p *domain.Project) { p.ProjectNo = no }
}

func WithProjectStatus(status string) ProjectOption {
	return func(p *domain.Project) { p.Status = status }
}

func WithProjectPriority(priority string) ProjectOption {
	return func(p *domain.Project) { p.Priority = priority }
}

func WithOwnerID(ownerID int) ProjectOption {
	return func(p *domain.Project) { p.OwnerID = ownerID }
}

func NewTestProject(opts ...ProjectOption) *domain.Project {
	now := time.Now()
	p := &domain.Project{
		ID:        1,
		Name:      "Test Project",
		ProjectNo: "PRJ0001",
		Status:    "TODO",
		Priority:  "MEDIUM",
		OwnerID:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// =============================================================================
// Task fixtures
// =============================================================================

type TaskOption func(*domain.Task)

func WithTaskID(id int) TaskOption {
	return func(t *domain.Task) { t.ID = id }
}

func WithTaskTitle(title string) TaskOption {
	return func(t *domain.Task) { t.Title = title }
}

func WithTaskStatus(status string) TaskOption {
	return func(t *domain.Task) { t.Status = status }
}

func WithTaskPriority(priority string) TaskOption {
	return func(t *domain.Task) { t.Priority = priority }
}

func WithTaskProjectID(projectID int) TaskOption {
	return func(t *domain.Task) { t.ProjectID = projectID }
}

func WithAssigneeID(assigneeID *int) TaskOption {
	return func(t *domain.Task) { t.AssigneeID = assigneeID }
}

func NewTestTask(opts ...TaskOption) *domain.Task {
	now := time.Now()
	t := &domain.Task{
		ID:        1,
		Title:     "Test Task",
		Status:    "TODO",
		Priority:  "MEDIUM",
		ProjectID: 1,
		CreatorID: 1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// =============================================================================
// News fixtures
// =============================================================================

type NewsOption func(*domain.News)

func WithNewsID(id int) NewsOption {
	return func(n *domain.News) { n.ID = id }
}

func WithNewsTitle(title string) NewsOption {
	return func(n *domain.News) { n.Title = title }
}

func WithNewsGroupID(groupID int) NewsOption {
	return func(n *domain.News) { n.GroupID = groupID }
}

func NewTestNews(opts ...NewsOption) *domain.News {
	now := time.Now()
	n := &domain.News{
		ID:        1,
		Title:     "Test News",
		GroupID:   1,
		CreatorID: 1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

// =============================================================================
// Article fixtures
// =============================================================================

type ArticleOption func(*domain.Article)

func WithArticleID(id int) ArticleOption {
	return func(a *domain.Article) { a.ID = id }
}

func WithArticleTitle(title string) ArticleOption {
	return func(a *domain.Article) { a.Title = title }
}

func WithArticleStatus(status string) ArticleOption {
	return func(a *domain.Article) { a.Status = status }
}

func NewTestArticle(opts ...ArticleOption) *domain.Article {
	now := time.Now()
	a := &domain.Article{
		ID:        1,
		Title:     "Test Article",
		ArticleNo: "ART0001",
		Status:    "DRAFT",
		AuthorID:  1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// =============================================================================
// MediaAccount fixtures
// =============================================================================

type MediaAccountOption func(*domain.MediaAccount)

func WithMediaAccountID(id int) MediaAccountOption {
	return func(a *domain.MediaAccount) { a.ID = id }
}

func WithMediaAccountName(name string) MediaAccountOption {
	return func(a *domain.MediaAccount) { a.Name = name }
}

func WithPlatform(platform string) MediaAccountOption {
	return func(a *domain.MediaAccount) { a.Platform = platform }
}

func NewTestMediaAccount(opts ...MediaAccountOption) *domain.MediaAccount {
	now := time.Now()
	a := &domain.MediaAccount{
		ID:        1,
		Name:      "Test Account",
		Platform:  "WEIBO",
		AccountID: "acc123",
		Status:    "ACTIVE",
		CreatedAt: now,
		UpdatedAt: now,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// =============================================================================
// MediaContent fixtures
// =============================================================================

type MediaContentOption func(*domain.MediaContent)

func WithMediaContentID(id int) MediaContentOption {
	return func(c *domain.MediaContent) { c.ID = id }
}

func WithMediaContentTitle(title string) MediaContentOption {
	return func(c *domain.MediaContent) { c.Title = title }
}

func WithMediaContentStatus(status string) MediaContentOption {
	return func(c *domain.MediaContent) { c.Status = status }
}

func WithMediaContentAccountID(accountID int) MediaContentOption {
	return func(c *domain.MediaContent) { c.AccountID = &accountID }
}

func NewTestMediaContent(opts ...MediaContentOption) *domain.MediaContent {
	now := time.Now()
	aid := 1
	c := &domain.MediaContent{
		ID:        1,
		Title:     "Test Content",
		AccountID: &aid,
		Status:    "DRAFT",
		CreatedAt: now,
		UpdatedAt: now,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// =============================================================================
// UploadFile fixtures
// =============================================================================

type UploadFileOption func(*domain.UploadFile)

func WithFileID(id int) UploadFileOption {
	return func(f *domain.UploadFile) { f.ID = id }
}

func WithFileName(name string) UploadFileOption {
	return func(f *domain.UploadFile) { f.Filename = name }
}

func WithFileType(fileType string) UploadFileOption {
	return func(f *domain.UploadFile) { f.FileType = fileType }
}

func NewTestUploadFile(opts ...UploadFileOption) *domain.UploadFile {
	now := time.Now()
	f := &domain.UploadFile{
		ID:               1,
		Filename:         "test.pdf",
		OriginalFilename: "test.pdf",
		FilePath:         "/uploads/test.pdf",
		FileType:         "document",
		FileSize:         1024,
		MimeType:         "application/pdf",
		Extension:        "pdf",
		UploaderID:       1,
		CreatedAt:        now,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// =============================================================================
// ApiToken fixtures
// =============================================================================

type ApiTokenOption func(*domain.ApiToken)

func WithTokenID(id int) ApiTokenOption {
	return func(t *domain.ApiToken) { t.ID = id }
}

func WithTokenName(name string) ApiTokenOption {
	return func(t *domain.ApiToken) { t.Name = name }
}

func WithTokenStatus(status string) ApiTokenOption {
	return func(t *domain.ApiToken) { t.Status = status }
}

func NewTestApiToken(opts ...ApiTokenOption) *domain.ApiToken {
	now := time.Now()
	t := &domain.ApiToken{
		ID:           1,
		UserID:       1,
		Name:         "test-token",
		TokenHash:    fmt.Sprintf("hash-%d", now.UnixNano()),
		Status:       "ACTIVE",
		UsageCount:   0,
		LastUsedAt:   nil,
		ExpiresAt:    nil,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// =============================================================================
// Helpers
// =============================================================================

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
