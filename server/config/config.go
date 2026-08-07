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

	// JWT 密钥必须跨进程重启保持不变，否则所有已登录设备上的令牌会立刻失效。
	// 数据目录在 Docker、飞牛和直接部署中都会持久化，因此用它保存自动生成的密钥。
	os.MkdirAll(cfg.DataDir, 0755)
	os.MkdirAll(filepath.Join(cfg.DataDir, "logs"), 0755)
	os.MkdirAll(cfg.UploadDir, 0755)

	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWTSecret = v
	}

	if cfg.JWTSecret == "" {
		cfg.JWTSecret = loadOrCreateJWTSecret(cfg.DataDir)
	}

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
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("failed to generate JWT secret: %v", err))
	}
	return hex.EncodeToString(b)
}

func loadOrCreateJWTSecret(dataDir string) string {
	secretPath := filepath.Join(dataDir, ".jwt_secret")

	if data, err := os.ReadFile(secretPath); err == nil {
		if secret := strings.TrimSpace(string(data)); secret != "" {
			return secret
		}
	}

	secret := generateRandomSecret()
	if err := os.WriteFile(secretPath, []byte(secret+"\n"), 0600); err != nil {
		fmt.Printf("[WARN] failed to persist JWT secret: %v; tokens will become invalid after restart\n", err)
		return secret
	}

	fmt.Printf("[INFO] generated persistent JWT secret at %s\n", secretPath)
	return secret
}
