package testutil

import (
	"context"
	"fmt"
	"time"

	"oa-nsdiy/backend/internal/domain"
)

// =============================================================================
// UserRepository stubs
// =============================================================================

type StubUserRepository struct{}

func (s *StubUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubUserRepository) UpdatePassword(ctx context.Context, id int, salt, hashedPassword string) error {
	return fmt.Errorf("not implemented")
}

func (s *StubUserRepository) SetVerificationCode(ctx context.Context, id int, code string, expiresAt time.Time) error {
	return fmt.Errorf("not implemented")
}

func (s *StubUserRepository) ClearVerificationCode(ctx context.Context, id int) error {
	return fmt.Errorf("not implemented")
}

// StubUserRepositoryWithData provides configurable user data for testing.
type StubUserRepositoryWithData struct {
	Users         map[int]*domain.User
	Err           error
	PasswordErr   error
	SetCodeErr    error
	ClearCodeErr  error
	UpdatePassErr error
}

func (s *StubUserRepositoryWithData) GetByID(ctx context.Context, id int) (*domain.User, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	if u, ok := s.Users[id]; ok {
		return u, nil
	}
	return nil, fmt.Errorf("not found")
}

func (s *StubUserRepositoryWithData) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	for _, u := range s.Users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (s *StubUserRepositoryWithData) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	for _, u := range s.Users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (s *StubUserRepositoryWithData) UpdatePassword(ctx context.Context, id int, salt, hashedPassword string) error {
	if s.UpdatePassErr != nil {
		return s.UpdatePassErr
	}
	if u, ok := s.Users[id]; ok {
		u.Salt = salt
		u.HashedPassword = hashedPassword
		u.TokenVersion++
		return nil
	}
	return fmt.Errorf("not found")
}

func (s *StubUserRepositoryWithData) SetVerificationCode(ctx context.Context, id int, code string, expiresAt time.Time) error {
	if s.SetCodeErr != nil {
		return s.SetCodeErr
	}
	if u, ok := s.Users[id]; ok {
		u.VerificationCode = &code
		u.VerificationCodeExpiresAt = &expiresAt
		return nil
	}
	return fmt.Errorf("not found")
}

func (s *StubUserRepositoryWithData) ClearVerificationCode(ctx context.Context, id int) error {
	if s.ClearCodeErr != nil {
		return s.ClearCodeErr
	}
	if u, ok := s.Users[id]; ok {
		u.VerificationCode = nil
		u.VerificationCodeExpiresAt = nil
		return nil
	}
	return fmt.Errorf("not found")
}

// =============================================================================
// RoleRepository stubs
// =============================================================================

type StubRoleRepository struct{}

func (s *StubRoleRepository) ListActive(ctx context.Context) ([]*domain.Role, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *StubRoleRepository) GetByID(ctx context.Context, id int) (*domain.Role, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubRoleRepository) GetByUserID(ctx context.Context, userID int) (*domain.Role, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubRoleRepository) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubRoleRepository) GetByCode(ctx context.Context, code string) (*domain.Role, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubRoleRepository) Create(ctx context.Context, role *domain.Role) error {
	return fmt.Errorf("not implemented")
}

func (s *StubRoleRepository) Update(ctx context.Context, role *domain.Role) error {
	return fmt.Errorf("not implemented")
}

func (s *StubRoleRepository) Delete(ctx context.Context, id int) error {
	return fmt.Errorf("not implemented")
}

func (s *StubRoleRepository) IsUsedByUsers(ctx context.Context, roleID int) (bool, error) {
	return false, fmt.Errorf("not implemented")
}

func (s *StubRoleRepository) GetPermissions(ctx context.Context, roleID int) ([]*domain.Permission, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *StubRoleRepository) UpdatePermissions(ctx context.Context, roleID int, permissionIDs []int) error {
	return fmt.Errorf("not implemented")
}

// StubRoleRepositoryWithData provides configurable role data for testing.
type StubRoleRepositoryWithData struct {
	Roles map[int]*domain.Role
	ByUID map[int]*domain.Role
	Perms map[int][]*domain.Permission
}

func (s *StubRoleRepositoryWithData) ListActive(ctx context.Context) ([]*domain.Role, error) {
	var result []*domain.Role
	for _, r := range s.Roles {
		if r.IsActive {
			result = append(result, r)
		}
	}
	return result, nil
}

func (s *StubRoleRepositoryWithData) GetByID(ctx context.Context, id int) (*domain.Role, error) {
	if r, ok := s.Roles[id]; ok {
		return r, nil
	}
	return nil, fmt.Errorf("not found")
}

func (s *StubRoleRepositoryWithData) GetByUserID(ctx context.Context, userID int) (*domain.Role, error) {
	if r, ok := s.ByUID[userID]; ok {
		return r, nil
	}
	return nil, nil
}

func (s *StubRoleRepositoryWithData) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	for _, r := range s.Roles {
		if r.Name == name {
			return r, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (s *StubRoleRepositoryWithData) GetByCode(ctx context.Context, code string) (*domain.Role, error) {
	for _, r := range s.Roles {
		if r.Code == code {
			return r, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (s *StubRoleRepositoryWithData) Create(ctx context.Context, role *domain.Role) error {
	role.ID = len(s.Roles) + 1
	s.Roles[role.ID] = role
	return nil
}

func (s *StubRoleRepositoryWithData) Update(ctx context.Context, role *domain.Role) error {
	if _, ok := s.Roles[role.ID]; !ok {
		return fmt.Errorf("not found")
	}
	s.Roles[role.ID] = role
	return nil
}

func (s *StubRoleRepositoryWithData) Delete(ctx context.Context, id int) error {
	delete(s.Roles, id)
	return nil
}

func (s *StubRoleRepositoryWithData) IsUsedByUsers(ctx context.Context, roleID int) (bool, error) {
	return false, nil
}

func (s *StubRoleRepositoryWithData) GetPermissions(ctx context.Context, roleID int) ([]*domain.Permission, error) {
	if s.Perms == nil {
		return nil, nil
	}
	return s.Perms[roleID], nil
}

func (s *StubRoleRepositoryWithData) UpdatePermissions(ctx context.Context, roleID int, permissionIDs []int) error {
	return nil
}

// =============================================================================
// PermissionValidator stubs
// =============================================================================

type StubPermissionValidator struct{}

func (s *StubPermissionValidator) GetByID(ctx context.Context, id int) (*domain.Permission, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubPermissionValidator) GetByRoleID(ctx context.Context, roleID int) ([]*domain.Permission, error) {
	return nil, fmt.Errorf("not implemented")
}

type StubPermissionValidatorWithData struct {
	Perms    map[int]*domain.Permission
	ByRoleID map[int][]*domain.Permission
}

func (s *StubPermissionValidatorWithData) GetByID(ctx context.Context, id int) (*domain.Permission, error) {
	if p, ok := s.Perms[id]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("not found")
}

func (s *StubPermissionValidatorWithData) GetByRoleID(ctx context.Context, roleID int) ([]*domain.Permission, error) {
	if s.ByRoleID == nil {
		return nil, nil
	}
	return s.ByRoleID[roleID], nil
}

// =============================================================================
// PermissionRepository stubs
// =============================================================================

type StubPermissionRepository struct{}

func (s *StubPermissionRepository) Search(ctx context.Context, resourceType, keyword string) ([]*domain.Permission, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *StubPermissionRepository) ListActive(ctx context.Context) ([]*domain.Permission, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *StubPermissionRepository) ListAll(ctx context.Context) ([]*domain.Permission, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *StubPermissionRepository) GetByID(ctx context.Context, id int) (*domain.Permission, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubPermissionRepository) GetByName(ctx context.Context, name string) (*domain.Permission, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubPermissionRepository) Create(ctx context.Context, perm *domain.Permission) error {
	return fmt.Errorf("not implemented")
}

func (s *StubPermissionRepository) Update(ctx context.Context, perm *domain.Permission) error {
	return fmt.Errorf("not implemented")
}

func (s *StubPermissionRepository) Delete(ctx context.Context, id int) error {
	return fmt.Errorf("not implemented")
}

func (s *StubPermissionRepository) IsUsedByRoles(ctx context.Context, permID int) (bool, error) {
	return false, fmt.Errorf("not implemented")
}

// =============================================================================
// ArticleRepository stubs
// =============================================================================

type StubArticleRepository struct{}

func (s *StubArticleRepository) List(ctx context.Context, keyword, status string, page, pageSize int) ([]*domain.ArticleDetail, int64, error) {
	return nil, 0, fmt.Errorf("not implemented")
}

func (s *StubArticleRepository) GetDetailByID(ctx context.Context, id int) (*domain.ArticleDetail, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubArticleRepository) GetByID(ctx context.Context, id int) (*domain.Article, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubArticleRepository) GetByTitle(ctx context.Context, title string) (*domain.Article, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubArticleRepository) GenerateArticleNo(ctx context.Context) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (s *StubArticleRepository) Create(ctx context.Context, article *domain.Article) error {
	return fmt.Errorf("not implemented")
}

func (s *StubArticleRepository) Update(ctx context.Context, article *domain.Article) error {
	return fmt.Errorf("not implemented")
}

func (s *StubArticleRepository) Delete(ctx context.Context, id int) error {
	return fmt.Errorf("not implemented")
}

func (s *StubArticleRepository) GetMaxVersionNo(ctx context.Context, articleID int) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (s *StubArticleRepository) CreateVersion(ctx context.Context, version *domain.ArticleVersion) error {
	return fmt.Errorf("not implemented")
}

func (s *StubArticleRepository) GetVersions(ctx context.Context, articleID int, page, pageSize int) ([]*domain.ArticleVersionDetail, int64, error) {
	return nil, 0, fmt.Errorf("not implemented")
}

func (s *StubArticleRepository) GetVersionByID(ctx context.Context, articleID, versionID int) (*domain.ArticleVersionDetail, error) {
	return nil, fmt.Errorf("not found")
}

// =============================================================================
// NewsRepository stubs
// =============================================================================

type StubNewsRepository struct{}

func (s *StubNewsRepository) List(ctx context.Context, groupID *int, keyword string, page, pageSize int) ([]*domain.NewsDetail, int64, error) {
	return nil, 0, fmt.Errorf("not implemented")
}

func (s *StubNewsRepository) GetDetailByID(ctx context.Context, id int) (*domain.NewsDetail, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubNewsRepository) GetByID(ctx context.Context, id int) (*domain.News, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubNewsRepository) Create(ctx context.Context, news *domain.News) error {
	return fmt.Errorf("not implemented")
}

func (s *StubNewsRepository) Update(ctx context.Context, news *domain.News) error {
	return fmt.Errorf("not implemented")
}

func (s *StubNewsRepository) Delete(ctx context.Context, id int) error {
	return fmt.Errorf("not implemented")
}

// =============================================================================
// NewsGroupValidator stubs
// =============================================================================

type StubNewsGroupValidator struct{}

func (s *StubNewsGroupValidator) GetByID(ctx context.Context, id int) (*domain.NewsGroup, error) {
	return nil, fmt.Errorf("not found")
}

type StubNewsGroupValidatorWithData struct {
	Groups map[int]*domain.NewsGroup
}

func (s *StubNewsGroupValidatorWithData) GetByID(ctx context.Context, id int) (*domain.NewsGroup, error) {
	if g, ok := s.Groups[id]; ok {
		return g, nil
	}
	return nil, fmt.Errorf("not found")
}

// =============================================================================
// NewsGroupRepository stubs
// =============================================================================

type StubNewsGroupRepository struct{}

func (s *StubNewsGroupRepository) ListAll(ctx context.Context) ([]*domain.NewsGroup, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *StubNewsGroupRepository) GetByID(ctx context.Context, id int) (*domain.NewsGroup, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubNewsGroupRepository) GetByName(ctx context.Context, name string) (*domain.NewsGroup, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubNewsGroupRepository) Create(ctx context.Context, group *domain.NewsGroup) error {
	return fmt.Errorf("not implemented")
}

func (s *StubNewsGroupRepository) Update(ctx context.Context, group *domain.NewsGroup) error {
	return fmt.Errorf("not implemented")
}

func (s *StubNewsGroupRepository) Delete(ctx context.Context, id int) error {
	return fmt.Errorf("not implemented")
}

func (s *StubNewsGroupRepository) GetNewsCount(ctx context.Context, groupID int) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

// =============================================================================
// ProjectRepository stubs
// =============================================================================

type StubProjectRepository struct{}

func (s *StubProjectRepository) ListByUserID(ctx context.Context, userID int, keyword string, page, pageSize int) ([]*domain.ProjectDetail, int64, error) {
	return nil, 0, fmt.Errorf("not implemented")
}

func (s *StubProjectRepository) GetByID(ctx context.Context, id int) (*domain.Project, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubProjectRepository) GetByProjectNo(ctx context.Context, projectNo string) (*domain.Project, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubProjectRepository) GenerateUniqueProjectNo(ctx context.Context) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (s *StubProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	return fmt.Errorf("not implemented")
}

func (s *StubProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	return fmt.Errorf("not implemented")
}

func (s *StubProjectRepository) Delete(ctx context.Context, id int) error {
	return fmt.Errorf("not implemented")
}

func (s *StubProjectRepository) AddMember(ctx context.Context, member *domain.ProjectMember) error {
	return fmt.Errorf("not implemented")
}

func (s *StubProjectRepository) GetMembers(ctx context.Context, projectID int) ([]*domain.ProjectMemberDetail, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *StubProjectRepository) GetMember(ctx context.Context, projectID, userID int) (*domain.ProjectMember, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubProjectRepository) UpdateMemberRole(ctx context.Context, projectID, userID int, role string) error {
	return fmt.Errorf("not implemented")
}

func (s *StubProjectRepository) RemoveMember(ctx context.Context, projectID, userID int) error {
	return fmt.Errorf("not implemented")
}

func (s *StubProjectRepository) IsMember(ctx context.Context, projectID, userID int) (bool, error) {
	return false, fmt.Errorf("not implemented")
}

// =============================================================================
// TaskRepository stubs
// =============================================================================

type StubTaskRepository struct{}

func (s *StubTaskRepository) List(ctx context.Context, projectID *int, status, priority string, assigneeID *int, page, pageSize int) ([]*domain.TaskDetail, int64, error) {
	return nil, 0, fmt.Errorf("not implemented")
}

func (s *StubTaskRepository) GetDetailByID(ctx context.Context, id int) (*domain.TaskDetail, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubTaskRepository) GetByID(ctx context.Context, id int) (*domain.Task, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	return fmt.Errorf("not implemented")
}

func (s *StubTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	return fmt.Errorf("not implemented")
}

func (s *StubTaskRepository) Delete(ctx context.Context, id int) error {
	return fmt.Errorf("not implemented")
}

// =============================================================================
// ProjectValidator stubs
// =============================================================================

type StubProjectValidator struct{}

func (s *StubProjectValidator) GetByID(ctx context.Context, id int) (*domain.Project, error) {
	return nil, fmt.Errorf("not found")
}

type StubProjectValidatorWithData struct {
	Projects map[int]*domain.Project
}

func (s *StubProjectValidatorWithData) GetByID(ctx context.Context, id int) (*domain.Project, error) {
	if p, ok := s.Projects[id]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("not found")
}

// =============================================================================
// MediaAccountRepository stubs
// =============================================================================

type StubMediaAccountRepository struct{}

func (s *StubMediaAccountRepository) List(ctx context.Context, keyword, platform, status string, page, pageSize int) ([]*domain.MediaAccount, int64, error) {
	return nil, 0, fmt.Errorf("not implemented")
}

func (s *StubMediaAccountRepository) GetByID(ctx context.Context, id int) (*domain.MediaAccount, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubMediaAccountRepository) Create(ctx context.Context, account *domain.MediaAccount) error {
	return fmt.Errorf("not implemented")
}

func (s *StubMediaAccountRepository) Update(ctx context.Context, account *domain.MediaAccount) error {
	return fmt.Errorf("not implemented")
}

func (s *StubMediaAccountRepository) Delete(ctx context.Context, id int) error {
	return fmt.Errorf("not implemented")
}

// =============================================================================
// MediaContentRepository stubs
// =============================================================================

type StubMediaContentRepository struct{}

func (s *StubMediaContentRepository) List(ctx context.Context, keyword, platform, status string, accountID *int, page, pageSize int) ([]*domain.MediaContentDetail, int64, error) {
	return nil, 0, fmt.Errorf("not implemented")
}

func (s *StubMediaContentRepository) GetDetailByID(ctx context.Context, id int) (*domain.MediaContentDetail, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubMediaContentRepository) GetByID(ctx context.Context, id int) (*domain.MediaContent, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubMediaContentRepository) Create(ctx context.Context, content *domain.MediaContent) error {
	return fmt.Errorf("not implemented")
}

func (s *StubMediaContentRepository) Update(ctx context.Context, content *domain.MediaContent) error {
	return fmt.Errorf("not implemented")
}

func (s *StubMediaContentRepository) Delete(ctx context.Context, id int) error {
	return fmt.Errorf("not implemented")
}

func (s *StubMediaContentRepository) GetMaxVersionNo(ctx context.Context, contentID int) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (s *StubMediaContentRepository) CreateVersion(ctx context.Context, version *domain.MediaContentVersion) error {
	return fmt.Errorf("not implemented")
}

func (s *StubMediaContentRepository) DeleteVersions(ctx context.Context, contentID int) error {
	return fmt.Errorf("not implemented")
}

func (s *StubMediaContentRepository) ListVersions(ctx context.Context, contentID int, page, pageSize int) ([]*domain.MediaContentVersionDetail, int64, error) {
	return nil, 0, fmt.Errorf("not implemented")
}

func (s *StubMediaContentRepository) GetVersionByID(ctx context.Context, versionID int) (*domain.MediaContentVersionDetail, error) {
	return nil, fmt.Errorf("not found")
}

// =============================================================================
// FileRepository stubs
// =============================================================================

type StubFileRepository struct{}

func (s *StubFileRepository) List(ctx context.Context, keyword, fileType string, uploaderID *int, page, pageSize int) ([]*domain.UploadFileDetail, int64, error) {
	return nil, 0, fmt.Errorf("not implemented")
}

func (s *StubFileRepository) GetDetailByID(ctx context.Context, id int) (*domain.UploadFileDetail, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubFileRepository) GetByID(ctx context.Context, id int) (*domain.UploadFile, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubFileRepository) Create(ctx context.Context, file *domain.UploadFile) error {
	return fmt.Errorf("not implemented")
}

func (s *StubFileRepository) Update(ctx context.Context, file *domain.UploadFile) error {
	return fmt.Errorf("not implemented")
}

func (s *StubFileRepository) Delete(ctx context.Context, id int) error {
	return fmt.Errorf("not implemented")
}

// =============================================================================
// ApiTokenRepository stubs
// =============================================================================

type StubApiTokenRepository struct{}

func (s *StubApiTokenRepository) ListByUserID(ctx context.Context, userID int, keyword, status string, page, pageSize int) ([]*domain.ApiToken, int64, error) {
	return nil, 0, fmt.Errorf("not implemented")
}

func (s *StubApiTokenRepository) GetByID(ctx context.Context, id int) (*domain.ApiToken, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubApiTokenRepository) Create(ctx context.Context, token *domain.ApiToken) error {
	return fmt.Errorf("not implemented")
}

func (s *StubApiTokenRepository) Update(ctx context.Context, token *domain.ApiToken) error {
	return fmt.Errorf("not implemented")
}

func (s *StubApiTokenRepository) Delete(ctx context.Context, id int) error {
	return fmt.Errorf("not implemented")
}

func (s *StubApiTokenRepository) ValidateToken(ctx context.Context, tokenHash string) (*domain.ApiToken, error) {
	return nil, fmt.Errorf("not found")
}

func (s *StubApiTokenRepository) UpdateUsage(ctx context.Context, id int) error {
	return fmt.Errorf("not implemented")
}
