package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"notepad/database"
	"notepad/logger"
	"notepad/model"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

type updateUserRequest struct {
	Role     string `json:"role" binding:"omitempty,oneof=admin user"`
	Password string `json:"password" binding:"omitempty,min=6"`
}

func ListUsers(c *gin.Context) {
	users, err := model.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 检查是否已有管理员
	adminCount := 0
	database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&adminCount)
	if adminCount > 0 && req.Role == "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已存在管理员，无法创建新管理员"})
		return
	}

	// 默认角色为普通用户
	if req.Role == "" {
		req.Role = "user"
	}

	user, err := model.CreateUser(req.Username, req.Password, "", "", req.Role)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	logger.Audit("User created: %s (role: %s) by admin %s", req.Username, req.Role, c.GetString("username"))
	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 获取当前用户信息
	currentUserID, _ := c.Get("userID")
	isCurrentUser := id == currentUserID.(int64)
	currentUser, _ := model.GetUserByID(currentUserID.(int64))

	// 如果是当前登录管理员尝试将自己改为普通用户，拒绝
	if isCurrentUser && currentUser != nil && currentUser.Role == "admin" && req.Role == "user" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能将自己降为普通用户"})
		return
	}

	// 检查是否已有管理员，防止通过编辑将唯一管理员改为普通用户
	if req.Role == "user" {
		adminCount := 0
		database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&adminCount)
		user, _ := model.GetUserByID(id)
		if adminCount == 1 && user != nil && user.Role == "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不能移除唯一管理员的角色"})
			return
		}
	}

	if req.Role != "" {
		if err := model.UpdateUserRole(id, req.Role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户角色失败"})
			return
		}
	}

	if req.Password != "" {
		if err := model.UpdatePassword(id, req.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
			return
		}
	}

	logger.Audit("User updated: id=%d by admin %s", id, c.GetString("username"))
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	currentUserID, _ := c.Get("userID")
	if id == currentUserID.(int64) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除自己"})
		return
	}

	if err := model.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	logger.Audit("User deleted: id=%d by admin %s", id, c.GetString("username"))
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
