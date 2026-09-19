package handler

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 家族统一统计通道：赞赏支持计数上报到门户接收端。
// STATS_ENDPOINT 环境变量可覆盖（本地测试用）。
const defaultStatsEndpoint = "https://techfunway.wycto.cn/api/apps.online/refresh"

// statsHTTPClient 统计上报专用客户端，超时宁短勿长，不拖住用户请求
var statsHTTPClient = &http.Client{Timeout: 5 * time.Second}

// StatsDeviceType 部署形态标识（fnos/docker，裸二进制为空），启动时注入
var StatsDeviceType = ""

// statsDataDir 数据目录，用于持久化设备标识兜底文件，启动时注入
var statsDataDir = ""

// statsPayload 统计上报字段。隐私红线：不含主机名、用户名与任何用户数据。
type statsPayload struct {
	AppName    string  `json:"app_name"`
	Version    string  `json:"version"`
	DeviceType string  `json:"device_type,omitempty"`
	DeviceID   string  `json:"device_id"`
	OS         string  `json:"os"`
	Arch       string  `json:"arch"`
	Event      string  `json:"event,omitempty"`
	Amount     float64 `json:"amount,omitempty"`
}

// SetDataDir 注入数据目录（设备标识兜底文件所在）
func SetDataDir(dataDir string) {
	statsDataDir = dataDir
}

// DetectDeviceType 部署形态判定：飞牛网关模式优先，其次容器，裸二进制为空
func DetectDeviceType(fnosEnabled bool) string {
	if fnosEnabled {
		return "fnos"
	}
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return "docker"
	}
	return ""
}

// statsEndpoint 统计上报地址
func statsEndpoint() string {
	if e := os.Getenv("STATS_ENDPOINT"); e != "" {
		return e
	}
	return defaultStatsEndpoint
}

// statsDeviceID 匿名设备标识：优先系统机器标识（卸载重装应用不变、重装系统才变），
// 读不到时在数据目录持久化随机 ID 兜底。仅上报哈希结果，不含原始标识。
func statsDeviceID() string {
	source := machineSignature()
	if source == "" {
		source = persistentDeviceID()
	}
	sum := md5.Sum([]byte("notepad|" + source))
	return hex.EncodeToString(sum[:])
}

// machineSignature 读取操作系统安装时生成的机器标识
func machineSignature() string {
	switch runtime.GOOS {
	case "linux":
		for _, p := range []string{"/etc/machine-id", "/var/lib/dbus/machine-id"} {
			if b, err := os.ReadFile(p); err == nil {
				if s := string(bytes.TrimSpace(b)); s != "" {
					return s
				}
			}
		}
	case "darwin":
		out, err := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice").Output()
		if err == nil {
			for _, line := range strings.Split(string(out), "\n") {
				if !strings.Contains(line, "IOPlatformUUID") {
					continue
				}
				if i := strings.Index(line, "= \""); i >= 0 {
					return strings.Trim(line[i:], "= \"\n\r")
				}
			}
		}
	case "windows":
		out, err := exec.Command("reg", "query", `HKLM\SOFTWARE\Microsoft\Cryptography`, "/v", "MachineGuid").Output()
		if err == nil {
			fields := strings.Fields(string(out))
			if len(fields) >= 3 {
				return fields[len(fields)-1]
			}
		}
	}
	return ""
}

// persistentDeviceID 数据目录持久化随机 ID，容器重建、主机改名后仍保持同一身份
func persistentDeviceID() string {
	if statsDataDir == "" {
		hostname, _ := os.Hostname()
		return hostname + runtime.GOOS + runtime.GOARCH
	}
	path := filepath.Join(statsDataDir, "device.id")
	if b, err := os.ReadFile(path); err == nil {
		if s := string(bytes.TrimSpace(b)); s != "" {
			return s
		}
	}
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		hostname, _ := os.Hostname()
		return hostname + runtime.GOOS + runtime.GOARCH
	}
	id := hex.EncodeToString(buf)
	_ = os.WriteFile(path, []byte(id), 0o600)
	return id
}

// DonateSupport 用户在赞赏弹窗点击「已支持」后上报一次匿名支持计数。
// 金额仅作接收端统计标注，不做支付核验；上报失败不影响应用任何功能。
func DonateSupport(c *gin.Context) {
	var payload struct {
		Amount float64 `json:"amount"`
	}
	_ = c.ShouldBindJSON(&payload)
	if payload.Amount < 0 {
		payload.Amount = 0
	}

	req := statsPayload{
		AppName:    "notepad",
		Version:    Version,
		DeviceType: StatsDeviceType,
		DeviceID:   statsDeviceID(),
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		Event:      "donate_support",
		Amount:     payload.Amount,
	}
	body, _ := json.Marshal(req)
	resp, err := statsHTTPClient.Post(statsEndpoint(), "application/json", bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "支持计数发送失败，请稍后再试"})
		return
	}
	defer resp.Body.Close()
	// 统计服务返回非 2xx 视为发送失败（http.Post 不把 4xx/5xx 当作 err）
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.JSON(http.StatusBadGateway, gin.H{"error": "支持计数发送失败，请稍后再试"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
