package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/go-arco-admin/server/internal/config"
	"github.com/go-arco-admin/server/internal/handler"
	mw "github.com/go-arco-admin/server/internal/middleware"
	"github.com/go-arco-admin/server/internal/resp"
)

func New(db *gorm.DB, cfg config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	h := &handler.Handler{DB: db, Cfg: cfg}

	api := r.Group("/api", mw.OperationLog(db))
	api.GET("/health", func(c *gin.Context) { resp.OK(c, gin.H{"status": "up"}) })
	api.POST("/user/login", h.Login)

	authed := api.Group("", mw.Auth(db, cfg.JWTSecret))
	authed.POST("/user/logout", h.Logout)
	authed.POST("/user/info", h.Info)
	authed.PUT("/user/profile", h.UpdateProfile)
	authed.PUT("/user/password", h.ChangePassword)

	sys := authed.Group("/system", mw.RequireRoles("admin"))
	sys.GET("/overview", h.Overview)
	sys.GET("/users", h.ListUsers)
	sys.POST("/users", h.CreateUser)
	sys.PUT("/users/:id", h.UpdateUser)
	sys.DELETE("/users/:id", h.DeleteUser)
	sys.GET("/roles", h.ListRoles)
	sys.POST("/roles", h.CreateRole)
	sys.PUT("/roles/:id", h.UpdateRole)
	sys.DELETE("/roles/:id", h.DeleteRole)
	sys.GET("/logs", h.ListLogs)

	serveSPA(r, cfg.StaticDir)
	return r
}

// serveSPA serves the built frontend (if present) and falls back to
// index.html so that history-mode routes work on refresh.
func serveSPA(r *gin.Engine, dir string) {
	index := filepath.Join(dir, "index.html")
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, resp.Body{Code: 40400, Status: "fail", Msg: "接口不存在"})
			return
		}
		p := filepath.Join(dir, filepath.Clean("/"+c.Request.URL.Path))
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			c.File(p)
			return
		}
		if _, err := os.Stat(index); err == nil {
			c.File(index)
			return
		}
		c.String(http.StatusNotFound, "frontend not built: run `pnpm build` in web/")
	})
}
