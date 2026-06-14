package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/config"
	"oa-nsdiy/backend/internal/handler"
	"oa-nsdiy/backend/internal/middleware"
	"oa-nsdiy/backend/internal/server/routes"
	"oa-nsdiy/backend/internal/web"
)

type Server struct {
	Router *gin.Engine
	cfg    *config.Config
}

func NewServer(cfg *config.Config) *Server {
	if cfg.Server.IsDebug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestLogger())

	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
	}
	router.Use(cors.New(corsConfig))

	// Serve embedded frontend (only active when built with -tags embed)
	router.Use(web.ServeEmbeddedFrontend())

	return &Server{
		Router: router,
		cfg:    cfg,
	}
}

// SetupRoutes sets up all API routes with pre-built handlers from Wire dependency injection.
func (s *Server) SetupRoutes(
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
) {
	routes.SetupRoutes(
		s.Router,
		authMiddleware,
		gin.HandlerFunc(adminRBAC),
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
	)
}

// NewHTTPServer creates an *http.Server from the configured Gin engine.
// The caller is responsible for calling ListenAndServe() and Shutdown().
func (s *Server) NewHTTPServer(addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: s.Router,
	}
}