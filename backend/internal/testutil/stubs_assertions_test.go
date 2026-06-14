package testutil

import (
	"oa-nsdiy/backend/internal/service"
)

// Compile-time interface assertions to ensure stubs implement service interfaces.

var _ service.UserRepository = (*StubUserRepository)(nil)
var _ service.UserRepository = (*StubUserRepositoryWithData)(nil)

var _ service.RoleRepository = (*StubRoleRepository)(nil)
var _ service.RoleRepository = (*StubRoleRepositoryWithData)(nil)

var _ service.PermissionValidator = (*StubPermissionValidator)(nil)
var _ service.PermissionValidator = (*StubPermissionValidatorWithData)(nil)

var _ service.PermissionRepository = (*StubPermissionRepository)(nil)

var _ service.ArticleRepository = (*StubArticleRepository)(nil)

var _ service.NewsRepository = (*StubNewsRepository)(nil)

var _ service.NewsGroupValidator = (*StubNewsGroupValidator)(nil)
var _ service.NewsGroupValidator = (*StubNewsGroupValidatorWithData)(nil)

var _ service.NewsGroupRepository = (*StubNewsGroupRepository)(nil)

var _ service.ProjectRepository = (*StubProjectRepository)(nil)

var _ service.TaskRepository = (*StubTaskRepository)(nil)

var _ service.ProjectValidator = (*StubProjectValidator)(nil)
var _ service.ProjectValidator = (*StubProjectValidatorWithData)(nil)

var _ service.MediaAccountRepository = (*StubMediaAccountRepository)(nil)

var _ service.MediaContentRepository = (*StubMediaContentRepository)(nil)

var _ service.FileRepository = (*StubFileRepository)(nil)

var _ service.ApiTokenRepository = (*StubApiTokenRepository)(nil)
