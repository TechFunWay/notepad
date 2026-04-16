package router

import (
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"notepad/handler"
	"notepad/middleware"
	"notepad/static"
)

func Setup(uploadDir string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// 上传文件静态服务
	r.Static("/uploads", uploadDir)

	api := r.Group("/api")
	{
		// Public
		api.POST("/auth/register", handler.Register)
		api.POST("/auth/login", handler.Login)
		api.GET("/auth/security-question", handler.GetSecurityQuestion)
		api.POST("/auth/forgot-password", handler.ForgotPassword)
		api.GET("/public-config", handler.GetPublicConfig)
		api.GET("/version", handler.GetVersion)
		api.GET("/health", handler.Health)

		// Authenticated
		auth := api.Group("", middleware.RequireAuth())
		auth.POST("/auth/logout", handler.Logout)
		auth.POST("/auth/change-password", handler.ChangePassword)
		auth.POST("/upload", handler.Upload)
		auth.GET("/notes", handler.ListNotes)
		auth.POST("/notes", handler.CreateNote)
		auth.GET("/notes/tags", handler.GetAllTags)
		auth.GET("/notes/:id", handler.GetNote)
		auth.PUT("/notes/:id", handler.UpdateNote)
		auth.DELETE("/notes/:id", handler.DeleteNote)

		// Admin
		admin := auth.Group("", middleware.RequireAdmin())
		admin.GET("/users", handler.ListUsers)
		admin.POST("/users", handler.CreateUser)
		admin.PUT("/users/:id", handler.UpdateUser)
		admin.DELETE("/users/:id", handler.DeleteUser)
		admin.GET("/configs", handler.ListConfigs)
		admin.PUT("/configs/:key", handler.UpdateConfig)
	}

	// SPA static files - prioritize external dist directory for development
	distDir := "./static/dist"
	if _, err := os.Stat(distDir); err == nil {
		// Use external dist directory if it exists
		fileServer := http.FileServer(http.Dir(distDir))
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api") {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
				return
			}

			// Try to serve the file directly
			if _, err := os.Stat(distDir + path); err == nil {
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}

			// SPA fallback: serve index.html
			c.Request.URL.Path = "/"
			fileServer.ServeHTTP(c.Writer, c.Request)
		})
	} else {
		// Fall back to embedded filesystem
		distFS, err := fs.Sub(static.StaticFS, "dist")
		if err == nil {
			fileServer := http.FileServer(http.FS(distFS))
			r.NoRoute(func(c *gin.Context) {
				path := c.Request.URL.Path
				if strings.HasPrefix(path, "/api") {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}

				// Try to serve the file directly
				if _, err := distFS.Open(strings.TrimPrefix(path, "/")); err == nil {
					fileServer.ServeHTTP(c.Writer, c.Request)
					return
				}

				// SPA fallback: serve index.html
				c.Request.URL.Path = "/"
				fileServer.ServeHTTP(c.Writer, c.Request)
			})
		}
	}

	return r
}
