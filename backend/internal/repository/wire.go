package repository

import (
	"github.com/google/wire"
	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/internal/db"
)

// ProvideEntClient returns the initialized Ent client from the db package.
// db.Init() must be called before Wire injection.
func ProvideEntClient() *ent.Client {
	return db.Client
}

var ProviderSet = wire.NewSet(
	// Repository bundle
	NewRepositories,

	// Individual repositories
	NewUserRepository,
	NewRoleRepository,
	NewPermissionRepository,
	NewNewsGroupRepository,
	NewNewsRepository,
	NewArticleRepository,
	NewProjectRepository,
	NewTaskRepository,
	NewMediaAccountRepository,
	NewMediaContentRepository,
	NewFileRepository,
	NewApiTokenRepository,
)
