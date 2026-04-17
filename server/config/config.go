package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Port      int
	DataDir   string
	WebDir    string
	UploadDir string
	ShareDirs string
	JWTSecret string
}

func (c *Config) DBPath() string {
	return filepath.Join(c.DataDir, "notepad.db")
}

func (c *Config) LogDir() string {
	return filepath.Join(c.DataDir, "logs")
}

func Load(port int, dataDir, webDir, uploadDir, shareDirs string) *Config {
	cfg := &Config{
		Port:      port,
		DataDir:   dataDir,
		WebDir:    webDir,
		UploadDir: uploadDir,
		ShareDirs: shareDirs,
	}

	if cfg.Port <= 0 {
		if v := os.Getenv("PORT"); v != "" {
			if p, err := strconv.Atoi(v); err == nil && p > 0 {
				cfg.Port = p
			}
		}
		if cfg.Port <= 0 {
			cfg.Port = 8904
		}
	}

	if cfg.DataDir == "" {
		if v := os.Getenv("DATA_DIR"); v != "" {
			cfg.DataDir = v
		} else {
			cfg.DataDir = "./data"
		}
	}

	if cfg.WebDir == "" {
		if v := os.Getenv("WEB_DIR"); v != "" {
			cfg.WebDir = v
		} else {
			cfg.WebDir = "./www"
		}
	}

	if cfg.UploadDir == "" {
		if v := os.Getenv("UPLOAD_DIR"); v != "" {
			cfg.UploadDir = v
		} else {
			cfg.UploadDir = filepath.Join(cfg.DataDir, "upload")
		}
	}

	if cfg.ShareDirs == "" {
		if v := os.Getenv("SHARE_DIRS"); v != "" {
			cfg.ShareDirs = v
		}
	}

	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWTSecret = v
	}

	if cfg.JWTSecret == "" {
		cfg.JWTSecret = generateRandomSecret()
		fmt.Println("[WARN] JWT_SECRET not set, using random secret (tokens invalid on restart)")
	}

	// 创建数据目录结构
	os.MkdirAll(cfg.DataDir, 0755)
	os.MkdirAll(filepath.Join(cfg.DataDir, "logs"), 0755)
	os.MkdirAll(cfg.UploadDir, 0755)

	return cfg
}

func (c *Config) ShareDirPaths() []string {
	if c.ShareDirs == "" {
		return nil
	}
	parts := strings.Split(c.ShareDirs, ":")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func generateRandomSecret() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
