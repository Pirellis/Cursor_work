package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	StatusDisabled = 0
	StatusEnabled  = 1
)

type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:64;not null" json:"name"`
	Code        string    `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Description string    `gorm:"size:255" json:"description"`
	Status      int       `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type User struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Username    string     `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password    string     `gorm:"size:128;not null" json:"-"`
	Nickname    string     `gorm:"size:64" json:"nickname"`
	Email       string     `gorm:"size:128" json:"email"`
	Phone       string     `gorm:"size:32" json:"phone"`
	Avatar      string     `gorm:"size:255" json:"avatar"`
	RoleID      uint       `json:"roleId"`
	Role        Role       `json:"role"`
	Status      int        `gorm:"default:1" json:"status"`
	LastLoginAt *time.Time `json:"lastLoginAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type OperationLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:64;index" json:"username"`
	Action    string    `gorm:"size:128" json:"action"`
	Method    string    `gorm:"size:16" json:"method"`
	Path      string    `gorm:"size:255" json:"path"`
	IP        string    `gorm:"size:64" json:"ip"`
	Success   bool      `json:"success"`
	Latency   int64     `json:"latency"`
	CreatedAt time.Time `gorm:"index" json:"createdAt"`
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Role{}, &User{}, &OperationLog{})
}
