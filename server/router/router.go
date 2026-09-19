package router

import (
	"io/fs"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"notepad/fnos"
	"notepad/handler"
	"notepad/middleware"
	"notepad/static"
)

func Setup(uploadDir string, webDir string, fnOSConfig fnos.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	appGroup := r.Group("")
	appPrefix := ""
	if fnOSConfig.Enabled {
		appPrefix = strings.TrimSuffix(fnOSConfig.Prefix, "/")
		appGroup = r.Group(appPrefix)
	}

	// 上传文件静态服务
	appGroup.Static("/uploads", uploadDir)

	api := appGroup.Group("/api")
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
		if fnOSConfig.Enabled {
			fnos.RegisterRoutes(api.Group("/auth/fnos"))
		}

		// Authenticated
		auth := api.Group("", middleware.RequireAuth())
		auth.POST("/auth/logout", handler.Logout)
		// 赞赏支持计数（匿名设备统计，仅设备级信息，不含用户数据）
		auth.POST("/donate/support", handler.DonateSupport)
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

	// serveIndex 返回入口 HTML，并把 __FNOS_GATEWAY_FLAG__ 替换成「这次文档
	// 请求是不是从飞牛统一网关 socket 进来的」。只有服务端知道答案（两条监听
	// 器分别是网关 socket 与直连端口），而前端必须据此决定带不带应用自己的
	// Authorization：飞牛接入层会把请求里的 Authorization 当成它自己的会话
	// token，认不出就直接回 200 纯文本 "invalid token"，请求根本到不了应用。
	// 占位符写在字符串字面量里，开发服务器不替换时恒为 false，见
	// web/index.html 与 web/src/utils/gateway.js。
	readIndexHTML := func() ([]byte, bool) {
		if useWebDir {
			data, err := os.ReadFile(webDir + "/index.html")
			return data, err == nil
		}
		distDir := "./static/dist"
		if _, err := os.Stat(distDir + "/index.html"); err == nil {
			data, err := os.ReadFile(distDir + "/index.html")
			return data, err == nil
		}
		if distFS, err := fs.Sub(static.StaticFS, "dist"); err == nil {
			data, err := fs.ReadFile(distFS, "index.html")
			return data, err == nil
		}
		return nil, false
	}
	serveIndex := func(c *gin.Context) {
		data, ok := readIndexHTML()
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		page := strings.ReplaceAll(string(data), "__FNOS_GATEWAY_FLAG__", strconv.FormatBool(fnos.IsGatewayRequest(c.Request.Context())))
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(page))
	}

	if fileServer != nil {
		r.NoRoute(func(c *gin.Context) {
			reqPath := c.Request.URL.Path
			if fnOSConfig.Enabled {
				if reqPath != appPrefix && !strings.HasPrefix(reqPath, appPrefix+"/") {
					c.Status(http.StatusNotFound)
					return
				}
				reqPath = strings.TrimPrefix(reqPath, appPrefix)
				if reqPath == "" {
					reqPath = "/"
				}
			}
			serveFile := func(filePath string) {
				original := c.Request.URL.Path
				c.Request.URL.Path = filePath
				fileServer.ServeHTTP(c.Writer, c.Request)
				c.Request.URL.Path = original
			}
			if reqPath == "/manifest.webmanifest" {
				c.Header("Content-Type", "application/manifest+json")
				serveFile(reqPath)
				return
			}
			if strings.HasPrefix(reqPath, "/api") {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
				return
			}

			// 入口 HTML 一律走 serveIndex 注入网关标志，且禁止缓存，确保升级
			// 后浏览器总能拿到最新的资源引用与正确的网关标志
			if reqPath == "/" || reqPath == "/index.html" {
				serveIndex(c)
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
				serveFile(reqPath)
				return
			}

			// 文件不存在且是带扩展名的静态资源请求（如升级后已不存在的旧 hash 文件）：
			// 返回 404，绝不回退 index.html，否则浏览器会把 HTML 当作 JS/CSS 执行导致白屏
			if path.Ext(reqPath) != "" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
				return
			}

			// SPA 前端路由（如 /login、/notes-list）：回退 index.html
			serveIndex(c)
		})
	}

	return r
}
