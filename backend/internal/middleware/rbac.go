package middleware

import (
	"oa-nsdiy/backend/internal/pkg/response"
	"oa-nsdiy/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// RBAC creates middleware that checks if the authenticated user has the required permissions.
func RBAC(roleService *service.RoleService, requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			response.Unauthorized(c, "not_authenticated", "User not authenticated")
			c.Abort()
			return
		}

		hasPermission, err := roleService.HasPermission(c.Request.Context(), userID, requiredPermissions...)
		if err != nil || !hasPermission {
			response.Forbidden(c, "access_denied", "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}
