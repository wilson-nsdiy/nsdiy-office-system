package handler

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAuthHandler,
	NewRoleHandler,
	NewPermissionHandler,
	NewNewsHandler,
	NewArticleHandler,
	NewProjectHandler,
	NewTaskHandler,
	NewMediaHandler,
	NewFileHandler,
	NewApiTokenHandler,
)
