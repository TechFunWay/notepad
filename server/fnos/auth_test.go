package fnos

import (
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"notepad/auth"
	"notepad/database"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	if err := database.Init(filepath.Join(t.TempDir(), "notepad.db"), "1.2.1"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.DB.Close() })
	auth.Init("test-secret")
}

func TestIdentityRequiresGatewayConnection(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/api/auth/fnos/login", nil)
	c.Request.Header.Set("X-Trim-Userid", "1001")
	c.Request.Header.Set("X-Trim-Username", "spoofed")
	if _, ok := IdentityFromRequest(c); ok {
		t.Fatal("TCP request must not be trusted from headers alone")
	}
}

func TestRegisterBindAndLogin(t *testing.T) {
	setupTestDB(t)
	identity := Identity{UserID: 1001, Username: "nas-user"}
	user, token, err := Bind(identity, "register", "admin", "hashed-client-password")
	if err != nil {
		t.Fatal(err)
	}
	if user.Role != "admin" || token == "" {
		t.Fatalf("unexpected registration result: %#v token=%q", user, token)
	}
	loggedIn, loginToken, err := Login(identity)
	if err != nil {
		t.Fatal(err)
	}
	if loggedIn.ID != user.ID || loginToken == "" {
		t.Fatalf("unexpected login result: %#v token=%q", loggedIn, loginToken)
	}
}

func TestNASRegisterCreatesFirstAdmin(t *testing.T) {
	setupTestDB(t)
	identity := Identity{UserID: 1002, Username: "nas-admin"}
	user, token, err := Bind(identity, "nas_register", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if user.Username != identity.Username || user.Role != "admin" || token == "" {
		t.Fatalf("unexpected NAS registration result: %#v token=%q", user, token)
	}
	if _, _, err := Login(identity); err != nil {
		t.Fatal(err)
	}
}
