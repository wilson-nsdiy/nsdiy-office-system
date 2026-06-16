package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"oa-nsdiy/backend/internal/config"
	"oa-nsdiy/backend/internal/handler"
	"oa-nsdiy/backend/internal/middleware"
)

// ProvideRouter creates the Gin engine with all routes wired up.
// This is the entry point for the server layer in the Wire dependency graph.
func ProvideRouter(
	cfg *config.Config,
	authMiddleware gin.HandlerFunc,
	adminRBAC middleware.AdminRBACMiddleware,
	authHandler *handler.AuthHandler,
	roleHandler *handler.RoleHandler,
	permissionHandler *handler.PermissionHandler,
	newsHandler *handler.NewsHandler,
	articleHandler *handler.ArticleHandler,
	projectHandler *handler.ProjectHandler,
	taskHandler *handler.TaskHandler,
	mediaHandler *handler.MediaHandler,
	fileHandler *handler.FileHandler,
	apiTokenHandler *handler.ApiTokenHandler,
	setupHandler *handler.SetupHandler,
) *gin.Engine {
	srv := NewServer(cfg)
	srv.SetupRoutes(
		authMiddleware,
		adminRBAC,
		authHandler,
		roleHandler,
		permissionHandler,
		newsHandler,
		articleHandler,
		projectHandler,
		taskHandler,
		mediaHandler,
		fileHandler,
		apiTokenHandler,
		setupHandler,
	)
	return srv.Router
}

var ProviderSet = wire.NewSet(
	NewServer,
	ProvideRouter,
)
