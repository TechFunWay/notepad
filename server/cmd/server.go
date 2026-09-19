package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"notepad/auth"
	"notepad/config"
	"notepad/database"
	"notepad/fnos"
	"notepad/handler"
	"notepad/logger"
	"notepad/router"
)

var appVersion = "dev"

func StartServer(port int, dataDir, webDir, uploadDir, shareDirs string, fnOSApp bool, gatewaySocket, gatewayPrefix string) {
	cfg := config.Load(port, dataDir, webDir, uploadDir, shareDirs)

	if err := logger.Init(cfg.LogDir()); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}
	defer logger.Close()

	logger.Info("Starting Notepad server on port %d", cfg.Port)

	if err := database.Init(cfg.DBPath(), appVersion); err != nil {
		logger.Error("Failed to initialize database: %v", err)
		return
	}

	auth.Init(cfg.JWTSecret)

	uploadPath := cfg.UploadDir
	if uploadPath == "" {
		uploadPath = filepath.Join(cfg.DataDir, "upload")
	}

	handler.SetUploadDir(uploadPath)
	// 部署形态与数据目录：匿名统计上报用
	handler.StatsDeviceType = handler.DetectDeviceType(fnOSApp)
	handler.SetDataDir(cfg.DataDir)

	gin.SetMode(gin.ReleaseMode)
	fnOSConfig := fnos.Config{Enabled: fnOSApp, Socket: gatewaySocket, Prefix: gatewayPrefix}
	r := router.Setup(uploadPath, cfg.WebDir, fnOSConfig)
	fmt.Printf("记事本已启动，访问地址：http://0.0.0.0:%d\n", cfg.Port)
	listeners, err := fnos.Listen(fmt.Sprintf(":%d", cfg.Port), fnOSConfig)
	if err != nil {
		logger.Error("Failed to listen: %v", err)
		return
	}
	for _, listener := range listeners {
		defer listener.Close()
	}
	if fnOSConfig.Enabled {
		defer os.Remove(fnOSConfig.Socket)
	}

	servers := make([]*http.Server, 0, len(listeners))
	for _, listener := range listeners {
		handler := http.Handler(r)
		if fnOSConfig.Enabled && !listener.TrustedGateway {
			gatewayURL := strings.TrimSuffix(fnOSConfig.Prefix, "/") + "/"
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.URL.Path == "/" {
					http.Redirect(w, req, gatewayURL, http.StatusTemporaryRedirect)
					return
				}
				r.ServeHTTP(w, req)
			})
		}
		server := fnos.NewHTTPServer(handler, listener.TrustedGateway)
		servers = append(servers, server)
		go func(server *http.Server, listener fnos.Listener) {
			if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
				logger.Error("Server error: %v", err)
			}
		}(server, listener)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, server := range servers {
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Server shutdown: %v", err)
		}
	}
}

func SetVersionInfo(version, buildTime, gitCommit string) {
	appVersion = version
	handler.Version = version
	handler.BuildTime = buildTime
	handler.GitCommit = gitCommit
}
