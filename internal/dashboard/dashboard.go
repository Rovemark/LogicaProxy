package dashboard

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed static/*
var staticFiles embed.FS

// RegisterRoutes mounts the dashboard UI at /dashboard.
func RegisterRoutes(engine *gin.Engine) {
	sub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return
	}

	engine.GET("/dashboard", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/dashboard/")
	})

	engine.GET("/dashboard/*filepath", func(c *gin.Context) {
		path := c.Param("filepath")
		if path == "/" || path == "" {
			path = "index.html"
		}
		// Strip leading slash for fs lookup
		if len(path) > 0 && path[0] == '/' {
			path = path[1:]
		}
		if path == "" {
			path = "index.html"
		}
		data, err := fs.ReadFile(sub, path)
		if err != nil {
			c.String(http.StatusNotFound, "not found")
			return
		}
		contentType := "text/html; charset=utf-8"
		if len(path) > 3 {
			switch path[len(path)-3:] {
			case ".js":
				contentType = "application/javascript"
			case "css":
				contentType = "text/css"
			case "son":
				contentType = "application/json"
			case "svg":
				contentType = "image/svg+xml"
			case "png":
				contentType = "image/png"
			}
		}
		c.Data(http.StatusOK, contentType, data)
	})
}
