package routes

import (
	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/handler"
)

// SetupRoutes registers all API routes using pre-built handlers and middleware.
// Components are injected via Wire dependency injection.
func SetupRoutes(
	router *gin.Engine,
	authMiddleware gin.HandlerFunc,
	adminRBAC gin.HandlerFunc,
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
) {
	api := router.Group("/api")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Setup routes
	authHandler.SetupRoutes(api, authMiddleware)
	roleHandler.SetupRoutes(api, authMiddleware, adminRBAC)
	permissionHandler.SetupRoutes(api, authMiddleware, adminRBAC)
	newsHandler.SetupRoutes(api, authMiddleware)
	articleHandler.SetupRoutes(api, authMiddleware)
	projectHandler.SetupRoutes(api, authMiddleware)
	taskHandler.SetupRoutes(api, authMiddleware)
	mediaHandler.SetupRoutes(api, authMiddleware)
	fileHandler.SetupRoutes(api, authMiddleware)
	apiTokenHandler.SetupRoutes(api, authMiddleware)
}
