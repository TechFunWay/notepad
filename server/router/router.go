package router

import (
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"notepad/handler"
	"notepad/middleware"
	"notepad/static"
)

func Setup(uploadDir string, webDir string) *gin.Engine {
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
		api.GET("/auth/setup-status", handler.GetSetupStatus)
		api.GET("/auth/security-question", handler.GetSecurityQuestion)
		api.POST("/auth/verify-answer", handler.VerifyAnswer)
		api.POST("/auth/forgot-password", handler.ForgotPassword)
		api.GET("/public-config", handler.GetPublicConfig)
		api.GET("/version", handler.GetVersion)
		api.GET("/health", handler.Health)

		// Authenticated
		auth := api.Group("", middleware.RequireAuth())
		auth.POST("/auth/logout", handler.Logout)
		auth.POST("/auth/change-password", handler.ChangePassword)
		auth.PUT("/auth/security-question", handler.UpdateSecurityQuestion)
		auth.POST("/upload", handler.Upload)
		auth.GET("/notes", handler.ListNotes)
		auth.POST("/notes", handler.CreateNote)
		auth.GET("/notes/tags", handler.GetAllTags)
		auth.PUT("/notes/tags", handler.RenameTag)
		auth.DELETE("/notes/tags", handler.DeleteTag)
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

	// SPA static files - priority: webDir > external dist > embedded
	var fileServer http.Handler
	var useWebDir bool

	if webDir != "" {
		if _, err := os.Stat(webDir); err == nil {
			fileServer = http.FileServer(http.Dir(webDir))
			useWebDir = true
		}
	}

	if !useWebDir {
		distDir := "./static/dist"
		if _, err := os.Stat(distDir); err == nil {
			fileServer = http.FileServer(http.Dir(distDir))
		} else {
			distFS, err := fs.Sub(static.StaticFS, "dist")
			if err == nil {
				fileServer = http.FileServer(http.FS(distFS))
			}
		}
	}

	// fileExists 判断请求的静态文件是否真实存在
	fileExists := func(reqPath string) bool {
		if useWebDir {
			if _, err := os.Stat(webDir + reqPath); err == nil {
				return true
			}
			return false
		}
		distDir := "./static/dist"
		if _, err := os.Stat(distDir); err == nil {
			if _, err := os.Stat(distDir + reqPath); err == nil {
				return true
			}
			return false
		}
		if distFS, err := fs.Sub(static.StaticFS, "dist"); err == nil {
			if f, err := distFS.Open(strings.TrimPrefix(reqPath, "/")); err == nil {
				f.Close()
				return true
			}
		}
		return false
	}

	if fileServer != nil {
		r.NoRoute(func(c *gin.Context) {
			reqPath := c.Request.URL.Path
			if reqPath == "/manifest.webmanifest" {
				c.Header("Content-Type", "application/manifest+json")
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
			if strings.HasPrefix(reqPath, "/api") {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
				return
			}

			// 命中真实存在的静态文件
			if fileExists(reqPath) {
				if strings.HasPrefix(reqPath, "/assets/") {
					// 文件名带内容 hash，内容变更必然换名，可长期强缓存
					c.Header("Cache-Control", "public, max-age=31536000, immutable")
				} else {
					// 其它根目录文件（favicon 等）每次回源校验
					c.Header("Cache-Control", "no-cache")
				}
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}

			// 文件不存在且是带扩展名的静态资源请求（如升级后已不存在的旧 hash 文件）：
			// 返回 404，绝不回退 index.html，否则浏览器会把 HTML 当作 JS/CSS 执行导致白屏
			if path.Ext(reqPath) != "" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
				return
			}

			// SPA 前端路由（如 /login、/notes-list）：回退 index.html
			// 入口文件禁止缓存，确保升级后浏览器总能拿到最新的资源引用，避免白屏
			c.Header("Cache-Control", "no-cache")
			c.Request.URL.Path = "/"
			fileServer.ServeHTTP(c.Writer, c.Request)
		})
	}

	return r
}
