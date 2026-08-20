package fnos

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"notepad/auth"
	"notepad/database"
	"notepad/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var ErrNotBound = errors.New("此飞牛 NAS 账号尚未绑定应用账号")

type Identity struct {
	UserID   int64
	Username string
	IsAdmin  bool
}

func IdentityFromRequest(c *gin.Context) (Identity, bool) {
	if !IsGatewayRequest(c.Request.Context()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请从飞牛桌面中的应用入口使用 NAS 登录"})
		return Identity{}, false
	}
	uid, err := strconv.ParseInt(c.GetHeader("X-Trim-Userid"), 10, 64)
	username := c.GetHeader("X-Trim-Username")
	if err != nil || uid <= 0 || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到飞牛 NAS 登录信息"})
		return Identity{}, false
	}
	return Identity{UserID: uid, Username: username, IsAdmin: c.GetHeader("X-Trim-Isadmin") == "true"}, true
}

func Login(identity Identity) (*model.User, string, error) {
	user := &model.User{}
	err := database.DB.QueryRow(`SELECT u.id, u.username, u.security_question, u.role, u.created_at, u.updated_at
		FROM users u JOIN fnos_bindings f ON f.user_id = u.id WHERE f.fnos_user_id = ?`, identity.UserID).
		Scan(&user.ID, &user.Username, &user.SecurityQuestion, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, "", ErrNotBound
	}
	if err != nil {
		return nil, "", err
	}
	if _, err := database.DB.Exec("UPDATE fnos_bindings SET fnos_username = ?, updated_at = CURRENT_TIMESTAMP WHERE fnos_user_id = ?", identity.Username, identity.UserID); err != nil {
		return nil, "", err
	}
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	return user, token, err
}

func Bind(identity Identity, mode, username, password string) (*model.User, string, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback()
	var count int
	if err := tx.QueryRow("SELECT COUNT(*) FROM fnos_bindings WHERE fnos_user_id = ?", identity.UserID).Scan(&count); err != nil {
		return nil, "", err
	}
	if count > 0 {
		return nil, "", fmt.Errorf("此飞牛 NAS 账号已绑定其他应用账号")
	}

	user := &model.User{}
	switch mode {
	case "nas_register":
		if err := tx.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
			return nil, "", err
		}
		if count > 0 {
			return nil, "", fmt.Errorf("仅首次初始化可直接创建飞牛管理员")
		}
		randomPassword := make([]byte, 32)
		if _, err := rand.Read(randomPassword); err != nil {
			return nil, "", err
		}
		hash, err := bcrypt.GenerateFromPassword(randomPassword, bcrypt.DefaultCost)
		if err != nil {
			return nil, "", err
		}
		result, err := tx.Exec("INSERT INTO users (username, password_hash, security_question, security_answer_hash, role) VALUES (?, ?, '', '', 'admin')", identity.Username, string(hash))
		if err != nil {
			return nil, "", err
		}
		user.ID, _ = result.LastInsertId()
		user.Username = identity.Username
		user.Role = "admin"
		user.CreatedAt = time.Now()
		user.UpdatedAt = user.CreatedAt
	case "register":
		if err := tx.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
			return nil, "", err
		}
		if count > 0 {
			var allow string
			if err := tx.QueryRow("SELECT value FROM configs WHERE key = 'allow_register'").Scan(&allow); err != nil || allow != "true" {
				return nil, "", fmt.Errorf("管理员已关闭注册")
			}
		}
		var exists int
		if err := tx.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&exists); err != nil {
			return nil, "", err
		}
		if exists > 0 {
			return nil, "", fmt.Errorf("用户名已存在")
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, "", err
		}
		role := "user"
		if count == 0 {
			role = "admin"
		}
		result, err := tx.Exec("INSERT INTO users (username, password_hash, security_question, security_answer_hash, role) VALUES (?, ?, '', '', ?)", username, string(hash), role)
		if err != nil {
			return nil, "", err
		}
		user.ID, _ = result.LastInsertId()
		user.Username = username
		user.Role = role
		user.CreatedAt = time.Now()
		user.UpdatedAt = user.CreatedAt
	case "bind":
		var stored string
		err := tx.QueryRow("SELECT id, username, password_hash, security_question, role, created_at, updated_at FROM users WHERE username = ?", username).
			Scan(&user.ID, &user.Username, &stored, &user.SecurityQuestion, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(stored), []byte(password)) != nil {
			return nil, "", fmt.Errorf("用户名或密码错误")
		}
	default:
		return nil, "", fmt.Errorf("无效的绑定方式")
	}
	if err := tx.QueryRow("SELECT COUNT(*) FROM fnos_bindings WHERE user_id = ?", user.ID).Scan(&count); err != nil {
		return nil, "", err
	}
	if count > 0 {
		return nil, "", fmt.Errorf("该应用账号已绑定其他飞牛 NAS 账号")
	}
	if _, err := tx.Exec("INSERT INTO fnos_bindings (user_id, fnos_user_id, fnos_username) VALUES (?, ?, ?)", user.ID, identity.UserID, identity.Username); err != nil {
		return nil, "", err
	}
	if err := tx.Commit(); err != nil {
		return nil, "", err
	}
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	return user, token, err
}

func RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/identity", func(c *gin.Context) {
		identity, ok := IdentityFromRequest(c)
		if !ok {
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"fnos_user_id": identity.UserID,
			"fnos_username": identity.Username,
			"is_admin":      identity.IsAdmin,
		})
	})
	group.POST("/login", func(c *gin.Context) {
		identity, ok := IdentityFromRequest(c)
		if !ok {
			return
		}
		user, token, err := Login(identity)
		if errors.Is(err, ErrNotBound) {
			c.JSON(http.StatusOK, gin.H{"binding_required": true, "fnos_username": identity.Username})
			return
		}
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "飞牛 NAS 登录失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
	})
	group.POST("/bind", func(c *gin.Context) {
		identity, ok := IdentityFromRequest(c)
		if !ok {
			return
		}
		var req struct {
			Mode     string `json:"mode"`
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.Mode == "" || (req.Mode != "nas_register" && (req.Username == "" || req.Password == "")) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请输入应用账号和密码"})
			return
		}
		user, token, err := Bind(identity, req.Mode, req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
	})
}
