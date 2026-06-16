//go:build !embed

package web

import (
	"github.com/gin-gonic/gin"
)

// ServeEmbeddedFrontend is a no-op middleware in non-embed builds.
// In development mode, frontend is served separately by Next.js dev server.
func ServeEmbeddedFrontend() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// HasEmbeddedFrontend returns false when frontend is not embedded (no embed build tag)
func HasEmbeddedFrontend() bool {
	return false
}