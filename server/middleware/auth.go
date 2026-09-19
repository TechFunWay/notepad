package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"notepad/auth"
	"notepad/fnos"
	"notepad/logger"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString != "" {
			if claims, err := auth.ParseToken(tokenString); err == nil {
				c.Set("userID", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
				c.Next()
				return
			} else {
				logger.Error("Auth failed: invalid token from %s, err: %v", c.ClientIP(), err)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "登录已过期"})
				c.Abort()
				return
			}
		}

		// 飞牛统一网关：网关在转发前已校验过 NAS 会话，并注入 X-Trim-Userid/
		// X-Trim-Username。网关域上不能依赖应用自己的 JWT——前端在网关域根本
		// 不发送 Authorization，因为飞牛接入层会把 Authorization 当成它自己的
		// 会话 token，认不出就直接回 200 纯文本 "invalid token"，请求压根进
		// 不了应用（表现即登录成功后「加载笔记失败」）。所以网关域的登录态以
		// 「该 NAS 用户已绑定的应用账号」为准，这正是飞牛给统一网关应用的身份
		// 通道；直连端口拿不到网关标记与身份头，仍然只认自己的 JWT。
		if identity, ok := fnos.GatewayIdentity(c); ok {
			if user, ok := fnos.ResolveUser(identity); ok {
				c.Set("userID", user.ID)
				c.Set("username", user.Username)
				c.Set("role", user.Role)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		c.Abort()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			logger.Error("Access denied: non-admin user %s attempted admin action", c.GetString("username"))
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	if cookie, err := c.Cookie("token"); err == nil {
		return cookie
	}

	return ""
}
