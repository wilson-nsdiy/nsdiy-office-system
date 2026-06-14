package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/config"
	"oa-nsdiy/backend/internal/pkg/response"
	"oa-nsdiy/backend/internal/service"
)

func AuthMiddleware(cfg *config.Config, authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "UNAUTHORIZED", "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "INVALID_AUTH_HEADER", "Authorization header format must be 'Bearer {token}'")
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := authService.VerifyToken(c.Request.Context(), tokenString)
		if err != nil {
			response.Unauthorized(c, "INVALID_TOKEN", "Invalid token")
			c.Abort()
			return
		}

		if claims.TokenType != "access" {
			response.Unauthorized(c, "INVALID_TOKEN_TYPE", "Invalid token type")
			c.Abort()
			return
		}

		// Fetch latest user info from database to validate TokenVersion
		user, err := authService.GetUserByID(c.Request.Context(), claims.UserID)
		if err != nil {
			response.Unauthorized(c, "USER_NOT_FOUND", "User not found")
			c.Abort()
			return
		}

		// Check user is active
		if !user.IsActive {
			response.Unauthorized(c, "USER_INACTIVE", "User account is not active")
			c.Abort()
			return
		}

		// Security: Validate TokenVersion to ensure token hasn't been invalidated by password change
		if claims.TokenVersion != user.TokenVersion {
			response.Unauthorized(c, "TOKEN_REVOKED", "Token has been revoked (password changed)")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("token_version", claims.TokenVersion)
		c.Next()
	}
}

func GetUserID(c *gin.Context) int {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	id, ok := userID.(int)
	if !ok {
		return 0
	}
	return id
}

func GetUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	name, ok := username.(string)
	if !ok {
		return ""
	}
	return name
}
