package model

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"notepad/database"
)

type User struct {
	ID                 int64     `json:"id"`
	Username           string    `json:"username"`
	PasswordHash       string    `json:"-"`
	SecurityQuestion   string    `json:"security_question"`
	SecurityAnswerHash string    `json:"-"`
	Role               string    `json:"role"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func CreateUser(username, password, securityQuestion, securityAnswer, role string) (*User, error) {
	// 检查是否已有用户，如果没有则默认设为管理员
	hasUsers := false
	database.DB.QueryRow("SELECT COUNT(*) > 0 FROM users").Scan(&hasUsers)
	if !hasUsers {
		role = "admin"
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var answerHash string
	if securityAnswer != "" {
		normalized := strings.ToLower(strings.TrimSpace(securityAnswer))
		hash, err := bcrypt.GenerateFromPassword([]byte(normalized), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		answerHash = string(hash)
	}

	result, err := database.DB.Exec(
		"INSERT INTO users (username, password_hash, security_question, security_answer_hash, role) VALUES (?, ?, ?, ?, ?)",
		username, string(passwordHash), securityQuestion, answerHash, role,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &User{
		ID:               id,
		Username:         username,
		SecurityQuestion: securityQuestion,
		Role:             role,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}, nil
}

func AuthenticateUser(username, password string) (*User, error) {
	user := &User{}
	err := database.DB.QueryRow(
		"SELECT id, username, password_hash, security_question, role, created_at, updated_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.SecurityQuestion, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	return user, nil
}

func GetUserByID(id int64) (*User, error) {
	user := &User{}
	err := database.DB.QueryRow(
		"SELECT id, username, security_question, role, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.SecurityQuestion, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := database.DB.QueryRow(
		"SELECT id, username, security_question, security_answer_hash, role, created_at, updated_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.SecurityQuestion, &user.SecurityAnswerHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func ListUsers() ([]User, error) {
	rows, err := database.DB.Query(
		"SELECT id, username, security_question, role, created_at, updated_at FROM users ORDER BY created_at",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.SecurityQuestion, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func UpdateUserRole(id int64, role string) error {
	_, err := database.DB.Exec("UPDATE users SET role = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", role, id)
	return err
}

func UpdatePassword(id int64, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = database.DB.Exec("UPDATE users SET password_hash = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", string(hash), id)
	return err
}

func UpdatePasswordByUsername(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	result, err := database.DB.Exec("UPDATE users SET password_hash = ?, updated_at = CURRENT_TIMESTAMP WHERE username = ?", string(hash), username)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("用户不存在")
	}
	return nil
}

func DeleteUser(id int64) error {
	_, err := database.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func CountUsers() (int, error) {
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

func VerifySecurityAnswer(username, answer string) (bool, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return false, fmt.Errorf("用户不存在")
	}

	if user.SecurityAnswerHash == "" {
		return false, fmt.Errorf("未设置安全问题")
	}

	normalized := strings.ToLower(strings.TrimSpace(answer))
	err = bcrypt.CompareHashAndPassword([]byte(user.SecurityAnswerHash), []byte(normalized))
	return err == nil, nil
}

func UpdateSecurityQuestion(id int64, question, answer string) error {
	var answerHash string
	if answer != "" {
		normalized := strings.ToLower(strings.TrimSpace(answer))
		hash, err := bcrypt.GenerateFromPassword([]byte(normalized), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		answerHash = string(hash)
	}
	_, err := database.DB.Exec("UPDATE users SET security_question = ?, security_answer_hash = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", question, answerHash, id)
	return err
}

func GetAdminUser() (*User, error) {
	user := &User{}
	err := database.DB.QueryRow(
		"SELECT id, username, role, created_at FROM users WHERE role = 'admin' ORDER BY id LIMIT 1",
	).Scan(&user.ID, &user.Username, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
