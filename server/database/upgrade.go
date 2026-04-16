package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"notepad/logger"
)

type upgrade struct {
	version string
	fn      func(tx *sql.Tx) error
}

// upgrades 按版本号升序排列，新增版本追加在末尾
var upgrades = []upgrade{
	{"1.0.0", upgrade_v1_0_0},
	// {"1.2.0", upgrade_v1_2_0},
	// {"1.3.0", nil},  // fn 为 nil 表示空升级，仅记录版本
}

func runUpgrades(appVersion string) error {
	appVersion = normalizeVersion(appVersion)

	// 建表
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS upgrade_records (
			version TEXT PRIMARY KEY,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return fmt.Errorf("failed to create upgrade_records table: %w", err)
	}

	// 读取当前 DB 版本
	var dbVersion string
	err := DB.QueryRow("SELECT COALESCE(MAX(version), '') FROM upgrade_records").Scan(&dbVersion)
	if err != nil {
		return fmt.Errorf("failed to read current db version: %w", err)
	}
	if dbVersion == "" {
		dbVersion = "0.0.0"
	}

	// 防降级：DB 版本不能高于二进制版本
	if compareVersions(dbVersion, appVersion) > 0 {
		return fmt.Errorf("database version (%s) is newer than binary version (%s), cannot downgrade", dbVersion, appVersion)
	}

	if dbVersion == appVersion {
		logger.Info("Version %s, no upgrade needed", appVersion)
		return nil
	}

	logger.Info("Upgrading from %s to %s", dbVersion, appVersion)

	// 按顺序执行升级脚本
	for _, u := range upgrades {
		if compareVersions(u.version, dbVersion) <= 0 {
			continue // 已执行过
		}
		if compareVersions(u.version, appVersion) > 0 {
			continue // 超出当前二进制版本
		}

		tx, err := DB.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for upgrade %s: %w", u.version, err)
		}

		if u.fn != nil {
			if err := u.fn(tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("upgrade %s failed: %w", u.version, err)
			}
		}

		if _, err := tx.Exec("INSERT INTO upgrade_records (version) VALUES (?)", u.version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record upgrade %s: %w", u.version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit upgrade %s: %w", u.version, err)
		}

		logger.Info("Upgrade to %s applied", u.version)
	}

	return nil
}

func upgrade_v1_0_0(tx *sql.Tx) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			security_question TEXT DEFAULT '',
			security_answer_hash TEXT DEFAULT '',
			role TEXT NOT NULL DEFAULT 'user',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL DEFAULT 'Untitled',
			content TEXT NOT NULL DEFAULT '',
			tags TEXT DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_notes_updated_at ON notes(updated_at)`,
		`CREATE TABLE IF NOT EXISTS configs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key TEXT NOT NULL UNIQUE,
			value TEXT NOT NULL DEFAULT '',
			description TEXT NOT NULL DEFAULT '',
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`INSERT OR IGNORE INTO configs (key, value, description) VALUES
			('allow_register', 'true', '允许新用户注册'),
			('site_title', '记事本', '应用标题'),
			('max_notes_per_user', '1000', '每用户最大笔记数')`,
	}

	for _, q := range queries {
		if _, err := tx.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func normalizeVersion(v string) string {
	if len(v) > 0 && (v[0] == 'v' || v[0] == 'V') {
		return v[1:]
	}
	return v
}

// compareVersions 比较两个 "X.Y.Z" 格式的版本号
// 返回: -1 (a<b), 0 (a==b), 1 (a>b)
func compareVersions(a, b string) int {
	pa := strings.Split(a, ".")
	pb := strings.Split(b, ".")

	maxLen := len(pa)
	if len(pb) > maxLen {
		maxLen = len(pb)
	}

	for i := 0; i < maxLen; i++ {
		var va, vb int
		if i < len(pa) {
			va, _ = strconv.Atoi(pa[i])
		}
		if i < len(pb) {
			vb, _ = strconv.Atoi(pb[i])
		}
		if va < vb {
			return -1
		}
		if va > vb {
			return 1
		}
	}
	return 0
}
