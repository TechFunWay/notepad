package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// processStartedAt 进程启动时间，前端用 started_at 区分「本次运行」与「重启后」：
// 刷新页面不重置，应用重启后变化，赞赏弹窗据此判断是否再次提示
var processStartedAt = time.Now()

func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":    Version,
		"build_time": BuildTime,
		"git_commit": GitCommit,
		"started_at": processStartedAt.Format(time.RFC3339),
	})
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
