package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"notepad/auth"
	"notepad/config"
	"notepad/database"
	"notepad/handler"
	"notepad/logger"
	"notepad/router"
)

var appVersion = "dev"

func StartServer(port int, dataDir, webDir, uploadDir, shareDirs string) {
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

	gin.SetMode(gin.ReleaseMode)
	r := router.Setup(uploadPath, cfg.WebDir)
	fmt.Printf("记事本已启动，访问地址：http://0.0.0.0:%d\n", cfg.Port)
	r.Run(fmt.Sprintf(":%d", cfg.Port))
}

func SetVersionInfo(version, buildTime, gitCommit string) {
	appVersion = version
	handler.Version = version
	handler.BuildTime = buildTime
	handler.GitCommit = gitCommit
}
