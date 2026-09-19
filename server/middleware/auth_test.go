package middleware

import (
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"notepad/auth"
	"notepad/database"
	"notepad/fnos"
)

func setupAuthTest(t *testing.T) {
	t.Helper()
	if err := database.Init(filepath.Join(t.TempDir(), "notepad.db"), "1.3.1"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.DB.Close() })
	auth.Init("test-secret")
}

func newAuthRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/me", RequireAuth(), func(c *gin.Context) {
		c.JSON(200, gin.H{"id": c.MustGet("userID").(int64), "username": c.GetString("username"), "role": c.GetString("role")})
	})
	return r
}

func TestRequireAuthFallsBackToGatewayIdentity(t *testing.T) {
	setupAuthTest(t)
	identity := fnos.Identity{UserID: 3001, Username: "nas-user"}
	if _, _, err := fnos.Bind(identity, "register", "walker", "pw"); err != nil {
		t.Fatal(err)
	}
	r := newAuthRouter()

	// 网关 socket 上的请求：不带 Authorization，凭网关注入的身份头解析出
	// 已绑定的应用账号。
	req := httptest.NewRequest("GET", "/me", nil)
	req = req.WithContext(fnos.GatewayContext(req.Context()))
	req.Header.Set("X-Trim-Userid", "3001")
	req.Header.Set("X-Trim-Username", "nas-user")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("gateway request without Authorization must authenticate, got %d: %s", w.Code, w.Body.String())
	}

	// 直连端口的 TCP 请求：同样的一组身份头属于伪造，必须拒绝。
	req = httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("X-Trim-Userid", "3001")
	req.Header.Set("X-Trim-Username", "nas-user")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Fatalf("TCP request must not be trusted from headers alone, got %d", w.Code)
	}

	// 直连端口的 TCP 请求：应用自己的 JWT 仍然有效。
	user, token, err := fnos.Login(identity)
	if err != nil || user == nil {
		t.Fatal(err)
	}
	req = httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("JWT on direct port must keep working, got %d: %s", w.Code, w.Body.String())
	}

	// 网关身份未绑定应用账号时不能放行。
	req = httptest.NewRequest("GET", "/me", nil)
	req = req.WithContext(fnos.GatewayContext(req.Context()))
	req.Header.Set("X-Trim-Userid", "9999")
	req.Header.Set("X-Trim-Username", "stranger")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Fatalf("unbound gateway identity must be rejected, got %d", w.Code)
	}
}
