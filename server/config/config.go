package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	Port      int
	DataDir   string
	JWTSecret string
}

func (c *Config) DBPath() string {
	return filepath.Join(c.DataDir, "db", "notepad.db")
}

func (c *Config) LogDir() string {
	return filepath.Join(c.DataDir, "logs")
}

func Load(port int, dataDir string) *Config {
	cfg := &Config{
		Port:    port,
		DataDir: dataDir,
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

	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWTSecret = v
	}

	if cfg.JWTSecret == "" {
		cfg.JWTSecret = generateRandomSecret()
		fmt.Println("[WARN] JWT_SECRET not set, using random secret (tokens invalid on restart)")
	}

	// 创建数据目录结构
	os.MkdirAll(filepath.Join(cfg.DataDir, "db"), 0755)
	os.MkdirAll(filepath.Join(cfg.DataDir, "logs"), 0755)
	os.MkdirAll(filepath.Join(cfg.DataDir, "upload"), 0755)

	return cfg
}

func generateRandomSecret() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
