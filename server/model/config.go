package model

import (
	"time"

	"notepad/database"
)

type Config struct {
	ID          int64     `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GetConfig(key string) (string, error) {
	var value string
	err := database.DB.QueryRow("SELECT value FROM configs WHERE key = ?", key).Scan(&value)
	return value, err
}

func GetAllConfigs() ([]Config, error) {
	rows, err := database.DB.Query("SELECT id, key, value, description, updated_at FROM configs ORDER BY key")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []Config
	for rows.Next() {
		var c Config
		if err := rows.Scan(&c.ID, &c.Key, &c.Value, &c.Description, &c.UpdatedAt); err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}
	return configs, nil
}

func SetConfig(key, value string) error {
	_, err := database.DB.Exec(
		"INSERT INTO configs (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP",
		key, value,
	)
	return err
}

func GetPublicConfigs() (map[string]string, error) {
	rows, err := database.DB.Query("SELECT key, value FROM configs WHERE key IN ('allow_register', 'site_title')")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		result[k] = v
	}
	return result, nil
}
