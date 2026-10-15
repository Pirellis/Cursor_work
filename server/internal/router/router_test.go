package router

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-arco-admin/server/internal/config"
	"github.com/go-arco-admin/server/internal/model"
)

type envelope struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func setup(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)
	cfg := config.Load()
	cfg.DBPath = filepath.Join(t.TempDir(), "test.db")
	cfg.StaticDir = t.TempDir()
	db, err := model.Open(cfg.DBPath)
	if err != nil {
		t.Fatal(err)
	}
	return New(db, cfg)
}

func call(t *testing.T, r *gin.Engine, method, path, token string, body any) envelope {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var e envelope
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Fatalf("%s %s: bad json %q", method, path, w.Body.String())
	}
	return e
}

func login(t *testing.T, r *gin.Engine, user, pw string) string {
	e := call(t, r, "POST", "/api/user/login", "", map[string]string{"username": user, "password": pw})
	if e.Code != 20000 {
		t.Fatalf("login %s: %+v", user, e)
	}
	var d struct{ Token string }
	_ = json.Unmarshal(e.Data, &d)
	return d.Token
}

func TestAuthAndRBAC(t *testing.T) {
	r := setup(t)

	if e := call(t, r, "POST", "/api/user/login", "", map[string]string{"username": "admin", "password": "bad"}); e.Code == 20000 {
		t.Fatal("wrong password accepted")
	}
	if e := call(t, r, "POST", "/api/user/info", "", nil); e.Code != 50008 {
		t.Fatalf("missing token: want 50008, got %d", e.Code)
	}

	admin := login(t, r, "admin", "admin")
	var info struct{ Role, Name string }
	_ = json.Unmarshal(call(t, r, "POST", "/api/user/info", admin, nil).Data, &info)
	if info.Role != "admin" {
		t.Fatalf("admin role = %q", info.Role)
	}

	user := login(t, r, "user", "user")
	if e := call(t, r, "GET", "/api/system/users", user, nil); e.Code != 40300 {
		t.Fatalf("normal user reached system api: %d", e.Code)
	}
}

func TestUserCRUDAndLogs(t *testing.T) {
	r := setup(t)
	admin := login(t, r, "admin", "admin")

	e := call(t, r, "POST", "/api/system/users", admin, map[string]any{"username": "alice", "password": "secret", "roleId": 2, "nickname": "Alice"})
	if e.Code != 20000 {
		t.Fatalf("create: %+v", e)
	}
	if e := call(t, r, "POST", "/api/system/users", admin, map[string]any{"username": "alice", "password": "secret", "roleId": 2}); e.Code == 20000 {
		t.Fatal("duplicate username accepted")
	}
	var created struct{ ID uint }
	_ = json.Unmarshal(e.Data, &created)

	var page struct {
		Total int64
		List  []struct{ Username string }
	}
	_ = json.Unmarshal(call(t, r, "GET", "/api/system/users?keyword=ali", admin, nil).Data, &page)
	if page.Total != 1 || page.List[0].Username != "alice" {
		t.Fatalf("search: %+v", page)
	}

	disabled := 0
	path := "/api/system/users/" + itoa(created.ID)
	if e := call(t, r, "PUT", path, admin, map[string]any{"roleId": 2, "status": disabled}); e.Code != 20000 {
		t.Fatalf("disable: %+v", e)
	}
	if e := call(t, r, "POST", "/api/user/login", "", map[string]string{"username": "alice", "password": "secret"}); e.Code == 20000 {
		t.Fatal("disabled user logged in")
	}
	if e := call(t, r, "DELETE", "/api/system/users/1", admin, nil); e.Code == 20000 {
		t.Fatal("admin deleted itself")
	}
	if e := call(t, r, "DELETE", "/api/system/roles/1", admin, nil); e.Code == 20000 {
		t.Fatal("builtin role deleted")
	}
	if e := call(t, r, "DELETE", path, admin, nil); e.Code != 20000 {
		t.Fatalf("delete: %+v", e)
	}

	var logs struct{ Total int64 }
	_ = json.Unmarshal(call(t, r, "GET", "/api/system/logs?result=fail", admin, nil).Data, &logs)
	if logs.Total < 3 {
		t.Fatalf("expected failed operations to be logged, got %d", logs.Total)
	}
}

func itoa(n uint) string { b, _ := json.Marshal(n); return string(b) }
