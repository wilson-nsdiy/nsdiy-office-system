//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/internal/config"
	"oa-nsdiy/backend/internal/handler"
	"oa-nsdiy/backend/internal/middleware"
	"oa-nsdiy/backend/internal/repository"
	"oa-nsdiy/backend/internal/server"
	"oa-nsdiy/backend/internal/service"
)

// Application holds the components built by Wire dependency injection.
type Application struct {
	Router  *gin.Engine
	Cleanup func()
}

// initializeApplication is the Wire injection function. It takes an
// already-initialized ent client and config, and builds the full
// dependency graph. The parameters themselves serve as Wire bindings.
func initializeApplication(client *ent.Client, cfg *config.Config) (*Application, error) {
	wire.Build(
		repository.ProviderSet,
		service.ProviderSet,
		handler.ProviderSet,
		middleware.ProviderSet,
		server.ProviderSet,

		// Bind concrete repository types to service layer interfaces.
		// Each repository pointer implements one or more service interfaces.
		wire.Bind(new(service.UserRepository), new(*repository.UserRepository)),
		wire.Bind(new(service.RoleRepository), new(*repository.RoleRepository)),
		wire.Bind(new(service.PermissionValidator), new(*repository.PermissionRepository)),
		wire.Bind(new(service.PermissionRepository), new(*repository.PermissionRepository)),
		wire.Bind(new(service.NewsGroupRepository), new(*repository.NewsGroupRepository)),
		wire.Bind(new(service.NewsGroupValidator), new(*repository.NewsGroupRepository)),
		wire.Bind(new(service.NewsRepository), new(*repository.NewsRepository)),
		wire.Bind(new(service.ArticleRepository), new(*repository.ArticleRepository)),
		wire.Bind(new(service.ProjectRepository), new(*repository.ProjectRepository)),
		wire.Bind(new(service.ProjectValidator), new(*repository.ProjectRepository)),
		wire.Bind(new(service.TaskRepository), new(*repository.TaskRepository)),
		wire.Bind(new(service.MediaAccountRepository), new(*repository.MediaAccountRepository)),
		wire.Bind(new(service.MediaContentRepository), new(*repository.MediaContentRepository)),
		wire.Bind(new(service.FileRepository), new(*repository.FileRepository)),
		wire.Bind(new(service.ApiTokenRepository), new(*repository.ApiTokenRepository)),

		// Cleanup provider
		provideCleanup,

		wire.Struct(new(Application), "Router", "Cleanup"),
	)
	return nil, nil
}

// provideCleanup creates a cleanup function that closes the ent client.
func provideCleanup(client *ent.Client) func() {
	return func() {
		_ = client.Close()
	}
}
