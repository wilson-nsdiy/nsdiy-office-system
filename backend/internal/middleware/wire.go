package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"oa-nsdiy/backend/internal/config"
	"oa-nsdiy/backend/internal/service"
)

// AdminRBACMiddleware is a named type for the admin RBAC handler, avoiding
// ambiguity with other gin.HandlerFunc values in the Wire graph.
type AdminRBACMiddleware gin.HandlerFunc

// ProvideAuthMiddleware creates the JWT auth middleware handler.
func ProvideAuthMiddleware(cfg *config.Config, authService *service.AuthService) gin.HandlerFunc {
	return AuthMiddleware(cfg, authService)
}

// ProvideAdminRBACMiddleware creates the admin RBAC middleware handler.
func ProvideAdminRBACMiddleware(roleService *service.RoleService) AdminRBACMiddleware {
	return AdminRBACMiddleware(RBAC(roleService, "admin:access"))
}

var ProviderSet = wire.NewSet(
	ProvideAuthMiddleware,
	ProvideAdminRBACMiddleware,
)
