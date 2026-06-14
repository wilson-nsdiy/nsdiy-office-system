//go:build embed

package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed all:dist
var frontendFS embed.FS

// ServeEmbeddedFrontend returns a Gin middleware that serves embedded frontend static files.
// For SPA routes (paths that don't match a static file), it serves index.html.
func ServeEmbeddedFrontend() gin.HandlerFunc {
	distFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		panic("failed to get dist subdirectory: " + err.Error())
	}
	fileServer := http.FileServer(http.FS(distFS))

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip API routes
		if shouldBypassFrontend(path) {
			c.Next()
			return
		}

		// Clean path and handle root
		cleanPath := strings.TrimPrefix(path, "/")
		if cleanPath == "" {
			cleanPath = "index.html"
		}

		// Check if file exists in embedded FS
		if file, err := distFS.Open(cleanPath); err == nil {
			_ = file.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		// SPA fallback: serve index.html for unmatched routes
		serveIndexHTML(c, distFS)
	}
}

func shouldBypassFrontend(path string) bool {
	trimmed := strings.TrimSpace(path)
	return strings.HasPrefix(trimmed, "/api/") || trimmed == "/api/health"
}

func serveIndexHTML(c *gin.Context, fsys fs.FS) {
	file, err := fsys.Open("index.html")
	if err != nil {
		c.String(http.StatusNotFound, "Frontend not found")
		c.Abort()
		return
	}
	defer func() { _ = file.Close() }()

	content, err := fs.ReadFile(fsys, "index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read index.html")
		c.Abort()
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	c.Abort()
}

// HasEmbeddedFrontend returns true when frontend is embedded (embed build tag)
func HasEmbeddedFrontend() bool {
	return true
}