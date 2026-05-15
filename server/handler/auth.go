package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"notepad/auth"
	"notepad/logger"
	"notepad/model"
)

type registerRequest struct {
	Username         string `json:"username" binding:"required,min=3,max=50"`
	Password         string `json:"password" binding:"required,min=6"`
	SecurityQuestion string `json:"security_question"`
	SecurityAnswer   string `json:"security_answer"`
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type updateSecurityQuestionRequest struct {
	SecurityQuestion string `json:"security_question"`
	SecurityAnswer   string `json:"security_answer"`
}

type forgotPasswordRequest struct {
	Username       string `json:"username" binding:"required"`
	SecurityAnswer string `json:"security_answer" binding:"required"`
	NewPassword    string `json:"new_password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	count, _ := model.CountUsers()

	if count > 0 {
		publicCfg, _ := model.GetPublicConfigs()
		if publicCfg["allow_register"] != "true" {
			c.JSON(http.StatusForbidden, gin.H{"error": "管理员已关闭注册"})
			return
		}
	}

	role := "user"
	if count == 0 {
		role = "admin"
	}

	user, err := model.CreateUser(req.Username, req.Password, req.SecurityQuestion, req.SecurityAnswer, role)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	logger.Audit("User registered: %s (role: %s)", req.Username, role)
	token, _ := auth.GenerateToken(user.ID, user.Username, user.Role)
	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user":  user,
	})
}

func Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入用户名和密码"})
		return
	}

	user, err := model.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		logger.Audit("Login failed: %s from %s", req.Username, c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	logger.Audit("Login success: %s from %s", req.Username, c.ClientIP())
	token, _ := auth.GenerateToken(user.ID, user.Username, user.Role)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "已登出"})
}

func GetSecurityQuestion(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供用户名"})
		return
	}

	user, err := model.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"security_question": user.SecurityQuestion})
}

type verifyAnswerRequest struct {
	Username       string `json:"username" binding:"required"`
	SecurityAnswer string `json:"security_answer" binding:"required"`
}

func VerifyAnswer(c *gin.Context) {
	var req verifyAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	valid, err := model.VerifySecurityAnswer(req.Username, req.SecurityAnswer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "安全答案错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "验证成功"})
}

func ForgotPassword(c *gin.Context) {
	var req forgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// Verify security answer first
	valid, err := model.VerifySecurityAnswer(req.Username, req.SecurityAnswer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "安全答案错误"})
		return
	}

	if err := model.UpdatePasswordByUsername(req.Username, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码重置失败"})
		return
	}

	logger.Audit("Password reset: %s from %s", req.Username, c.ClientIP())
	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}

func ChangePassword(c *gin.Context) {
	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	user, err := model.AuthenticateUser(username.(string), req.CurrentPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "当前密码错误"})
		return
	}

	if err := model.UpdatePassword(user.ID, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码修改失败"})
		return
	}

	_ = userID // used for auth check
	logger.Audit("Password changed: %s", username)
	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

func UpdateSecurityQuestion(c *gin.Context) {
	var req updateSecurityQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	if err := model.UpdateSecurityQuestion(userID.(int64), req.SecurityQuestion, req.SecurityAnswer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "安全问题修改失败"})
		return
	}

	logger.Audit("Security question updated: %s", username)
	c.JSON(http.StatusOK, gin.H{"message": "安全问题修改成功"})
}
