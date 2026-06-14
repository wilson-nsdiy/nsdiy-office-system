//go:build !embed

package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServeEmbeddedFrontend returns a handler that returns 404 for non-embed builds.
// In development mode, frontend is served separately by Next.js dev server.
func ServeEmbeddedFrontend() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusNotFound, "Frontend not embedded. Build with -tags embed to include frontend.")
		c.Abort()
	}
}

// HasEmbeddedFrontend returns false when frontend is not embedded (no embed build tag)
func HasEmbeddedFrontend() bool {
	return false
}