package service

import (
	"github.com/google/wire"
	"oa-nsdiy/backend/internal/config"
)

// ProvideAuthService creates AuthService by extracting JWT configuration
// from the app config.
func ProvideAuthService(userRepo UserRepository, cfg *config.Config) *AuthService {
	return NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)
}

var ProviderSet = wire.NewSet(
	ProvideAuthService,
	NewRoleService,
	NewPermissionService,
	NewNewsGroupService,
	NewNewsService,
	NewArticleService,
	NewProjectService,
	NewTaskService,
	NewMediaAccountService,
	NewMediaContentService,
	NewFileService,
	NewApiTokenService,
)
