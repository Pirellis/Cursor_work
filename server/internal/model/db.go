package model

import (
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const DefaultAvatar = ""

func Open(path string) (*gorm.DB, error) {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	db, err := gorm.Open(sqlite.Open(path+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Role{}, &User{}, &OperationLog{}); err != nil {
		return nil, err
	}
	return db, seed(db)
}

func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}

// seed creates the two roles the Arco Pro routes reference (admin, user)
// and one account for each, only when the database is empty.
func seed(db *gorm.DB) error {
	var count int64
	db.Model(&Role{}).Count(&count)
	if count > 0 {
		return nil
	}
	return db.Transaction(func(tx *gorm.DB) error {
		admin := Role{Code: "admin", Name: "超级管理员", Description: "拥有全部菜单与系统管理权限", Status: 1, Builtin: true}
		user := Role{Code: "user", Name: "普通用户", Description: "仅可访问基础业务页面", Status: 1, Builtin: true}
		if err := tx.Create(&[]*Role{&admin, &user}).Error; err != nil {
			return err
		}
		accounts := []struct {
			name, pw, nick, email, job string
			role                       uint
		}{
			{"admin", "admin", "管理员", "admin@example.com", "平台负责人", admin.ID},
			{"user", "user", "演示用户", "user@example.com", "运营专员", user.ID},
		}
		for _, a := range accounts {
			hash, err := HashPassword(a.pw)
			if err != nil {
				return err
			}
			u := User{Username: a.name, PasswordHash: hash, Nickname: a.nick, Email: a.email, Job: a.job,
				Avatar: DefaultAvatar, Status: 1, RoleID: a.role, Introduction: "Go + Arco Design Pro 通用后台"}
			if err := tx.Create(&u).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
