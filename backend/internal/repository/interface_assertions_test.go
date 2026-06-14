package repository_test

import (
	"oa-nsdiy/backend/internal/repository"
	"oa-nsdiy/backend/internal/service"
)

// Compile-time interface implementation checks.
// Using external test package (package repository_test) avoids circular imports.
var (
	_ service.UserRepository        = (*repository.UserRepository)(nil)
	_ service.RoleRepository        = (*repository.RoleRepository)(nil)
	_ service.PermissionValidator   = (*repository.PermissionRepository)(nil)
	_ service.PermissionRepository  = (*repository.PermissionRepository)(nil)
	_ service.ArticleRepository     = (*repository.ArticleRepository)(nil)
	_ service.NewsRepository        = (*repository.NewsRepository)(nil)
	_ service.NewsGroupValidator    = (*repository.NewsGroupRepository)(nil)
	_ service.NewsGroupRepository   = (*repository.NewsGroupRepository)(nil)
	_ service.ProjectRepository     = (*repository.ProjectRepository)(nil)
	_ service.ProjectValidator      = (*repository.ProjectRepository)(nil)
	_ service.TaskRepository        = (*repository.TaskRepository)(nil)
	_ service.MediaAccountRepository = (*repository.MediaAccountRepository)(nil)
	_ service.MediaContentRepository = (*repository.MediaContentRepository)(nil)
	_ service.FileRepository        = (*repository.FileRepository)(nil)
	_ service.ApiTokenRepository    = (*repository.ApiTokenRepository)(nil)
)
