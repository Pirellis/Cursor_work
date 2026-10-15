package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Pirellis/Cursor_work/server/internal/auth"
	"github.com/Pirellis/Cursor_work/server/internal/handler"
	"github.com/Pirellis/Cursor_work/server/internal/model"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	addr := env("ADDR", ":8080")
	dsn := env("DB_PATH", "data/admin.db")
	secret := env("JWT_SECRET", "change-me-in-production")
	webDir := env("WEB_DIR", "../web/dist")

	if err := os.MkdirAll(filepath.Dir(dsn), 0o755); err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
	if err != nil {
		log.Fatal(err)
	}
	if err := model.Migrate(db); err != nil {
		log.Fatal(err)
	}
	if err := seed(db); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	handler.New(db, auth.NewManager(secret, 24*time.Hour)).Register(r)
	serveSPA(r, webDir)

	log.Printf("Go Arco Admin listening on %s", addr)
	log.Fatal(r.Run(addr))
}

// serveSPA serves the built frontend when present, falling back to index.html for client routes.
func serveSPA(r *gin.Engine, dir string) {
	index := filepath.Join(dir, "index.html")
	if _, err := os.Stat(index); err != nil {
		return
	}
	fs := http.Dir(dir)
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"code": 40400, "msg": "接口不存在"})
			return
		}
		if f, err := fs.Open(c.Request.URL.Path); err == nil {
			stat, _ := f.Stat()
			f.Close()
			if stat != nil && !stat.IsDir() {
				c.FileFromFS(c.Request.URL.Path, fs)
				return
			}
		}
		c.File(index)
	})
}

func seed(db *gorm.DB) error {
	var n int64
	db.Model(&model.Role{}).Count(&n)
	if n > 0 {
		return nil
	}
	admin := model.Role{Name: "超级管理员", Code: "admin", Description: "拥有系统全部权限"}
	user := model.Role{Name: "普通用户", Code: "user", Description: "仅可访问业务页面"}
	if err := db.Create(&[]*model.Role{&admin, &user}).Error; err != nil {
		return err
	}
	adminPw, _ := handler.HashPassword("admin")
	userPw, _ := handler.HashPassword("user")
	users := []model.User{
		{Username: "admin", Password: adminPw, Nickname: "管理员", Email: "admin@example.com", RoleID: admin.ID, Status: model.StatusEnabled,
			Avatar: "//lf1-xgcdn-tos.pstatp.com/obj/vcloud/vadmin/start.8e0e4855ee346a46ccff8ff3e24db27b.png"},
		{Username: "user", Password: userPw, Nickname: "演示用户", Email: "user@example.com", RoleID: user.ID, Status: model.StatusEnabled},
	}
	return db.Create(&users).Error
}
