package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"notepad/logger"
	"notepad/model"
)

type updateConfigRequest struct {
	Value string `json:"value" binding:"required"`
}

func ListConfigs(c *gin.Context) {
	configs, err := model.GetAllConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取配置列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"configs": configs})
}

func UpdateConfig(c *gin.Context) {
	key := c.Param("key")
	var req updateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := model.SetConfig(key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新配置失败"})
		return
	}

	logger.Audit("Config updated: %s = %s by %s", key, req.Value, c.GetString("username"))
	c.JSON(http.StatusOK, gin.H{"message": "配置更新成功"})
}

func GetPublicConfig(c *gin.Context) {
	configs, err := model.GetPublicConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取配置失败"})
		return
	}
	c.JSON(http.StatusOK, configs)
}
