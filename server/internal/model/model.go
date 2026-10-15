package model

import "time"

type Base struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Role struct {
	Base
	Code        string `gorm:"uniqueIndex;size:32;not null" json:"code"`
	Name        string `gorm:"size:64;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Status      int    `gorm:"default:1" json:"status"`
	Builtin     bool   `json:"builtin"`
	UserCount   int64  `gorm:"-" json:"userCount"`
}

type User struct {
	Base
	Username     string     `gorm:"uniqueIndex;size:64;not null" json:"username"`
	PasswordHash string     `gorm:"size:128;not null" json:"-"`
	Nickname     string     `gorm:"size:64" json:"nickname"`
	Email        string     `gorm:"size:128" json:"email"`
	Phone        string     `gorm:"size:32" json:"phone"`
	Avatar       string     `gorm:"size:255" json:"avatar"`
	Job          string     `gorm:"size:64" json:"job"`
	Introduction string     `gorm:"size:255" json:"introduction"`
	Status       int        `gorm:"default:1" json:"status"`
	RoleID       uint       `json:"roleId"`
	Role         Role       `json:"role"`
	LastLoginAt  *time.Time `json:"lastLoginAt"`
}

type OperationLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"createdAt"`
	UserID    uint      `json:"userId"`
	Username  string    `gorm:"size:64;index" json:"username"`
	Method    string    `gorm:"size:8" json:"method"`
	Path      string    `gorm:"size:255" json:"path"`
	Status    int       `json:"status"`
	Code      int       `json:"code"`
	LatencyMs int64     `json:"latencyMs"`
	IP        string    `gorm:"size:64" json:"ip"`
	UserAgent string    `gorm:"size:255" json:"userAgent"`
}
