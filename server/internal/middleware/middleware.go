package middleware

import (
	"bytes"
	"errors"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"github.com/go-arco-admin/server/internal/model"
	"github.com/go-arco-admin/server/internal/resp"
)

const ctxUser = "currentUser"

type Claims struct {
	UserID uint `json:"uid"`
	jwt.RegisteredClaims
}

func SignToken(secret string, userID uint, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{UserID: userID, RegisteredClaims: jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// Auth validates the bearer token and loads the user on every request so that
// disabling an account or changing its role takes effect immediately.
func Auth(db *gorm.DB, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if raw == "" {
			resp.Fail(c, resp.CodeInvalidToken, "未登录")
			return
		}
		claims := &Claims{}
		_, err := jwt.ParseWithClaims(raw, claims, func(*jwt.Token) (any, error) { return []byte(secret), nil },
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
		if errors.Is(err, jwt.ErrTokenExpired) {
			resp.Fail(c, resp.CodeTokenExpired, "登录已过期，请重新登录")
			return
		}
		if err != nil {
			resp.Fail(c, resp.CodeInvalidToken, "无效的登录凭证")
			return
		}
		var u model.User
		if err := db.Preload("Role").First(&u, claims.UserID).Error; err != nil || u.Status != 1 {
			resp.Fail(c, resp.CodeInvalidToken, "账号不存在或已被禁用")
			return
		}
		c.Set(ctxUser, &u)
		c.Next()
	}
}

func CurrentUser(c *gin.Context) *model.User {
	u, _ := c.Get(ctxUser)
	user, _ := u.(*model.User)
	return user
}

func RequireRoles(codes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := CurrentUser(c)
		for _, code := range codes {
			if u != nil && u.Role.Code == code {
				c.Next()
				return
			}
		}
		resp.Fail(c, resp.CodeForbidden, "没有权限执行此操作")
	}
}

type bodyCapture struct {
	gin.ResponseWriter
	buf bytes.Buffer
}

func (w *bodyCapture) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

// OperationLog records every mutating request together with the business
// code from the response envelope, since HTTP status is always 200.
func OperationLog(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			c.Next()
			return
		}
		start := time.Now()
		w := &bodyCapture{ResponseWriter: c.Writer}
		c.Writer = w
		c.Next()

		var body struct {
			Code int `json:"code"`
		}
		_ = json.Unmarshal(w.buf.Bytes(), &body)
		entry := model.OperationLog{
			Method: c.Request.Method, Path: c.Request.URL.Path, Status: c.Writer.Status(), Code: body.Code,
			LatencyMs: time.Since(start).Milliseconds(), IP: c.ClientIP(), UserAgent: truncate(c.Request.UserAgent(), 255),
		}
		if u := CurrentUser(c); u != nil {
			entry.UserID, entry.Username = u.ID, u.Username
		} else if name, ok := c.Get("loginUsername"); ok {
			entry.Username, _ = name.(string)
		}
		db.Create(&entry)
	}
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}
