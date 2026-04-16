package cmd

import (
	"fmt"
	"path/filepath"

	"notepad/auth"
	"notepad/config"
	"notepad/database"
	"notepad/handler"
	"notepad/logger"
	"notepad/router"
)

var appVersion = "dev"

func StartServer(port int, dataDir string) {
	cfg := config.Load(port, dataDir)

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

	handler.SetUploadDir(filepath.Join(cfg.DataDir, "upload"))

	r := router.Setup(filepath.Join(cfg.DataDir, "upload"))
	r.Run(fmt.Sprintf(":%d", cfg.Port))
}

func SetVersionInfo(version, buildTime, gitCommit string) {
	appVersion = version
	handler.Version = version
	handler.BuildTime = buildTime
	handler.GitCommit = gitCommit
}
