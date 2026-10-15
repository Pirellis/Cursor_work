package handler

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Pirellis/Cursor_work/server/internal/auth"
	"github.com/Pirellis/Cursor_work/server/internal/model"
	"github.com/Pirellis/Cursor_work/server/internal/resp"
)

const ctxUser = "currentUser"

func (h *Handler) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if token == "" {
			resp.Fail(c, resp.CodeInvalidToken, "未登录")
			return
		}
		claims, err := h.jwt.Parse(token)
		if errors.Is(err, auth.ErrExpired) {
			resp.Fail(c, resp.CodeTokenExpired, "登录已过期，请重新登录")
			return
		}
		if err != nil {
			resp.Fail(c, resp.CodeInvalidToken, "无效的登录凭证")
			return
		}
		var user model.User
		if err := h.db.Preload("Role").First(&user, claims.UserID).Error; err != nil {
			resp.Fail(c, resp.CodeInvalidToken, "用户不存在")
			return
		}
		if user.Status != model.StatusEnabled {
			resp.Fail(c, resp.CodeInvalidToken, "账号已被禁用")
			return
		}
		c.Set(ctxUser, &user)
		c.Next()
	}
}

func (h *Handler) RequireRole(codes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := currentUser(c)
		for _, code := range codes {
			if user.Role.Code == code {
				c.Next()
				return
			}
		}
		resp.Fail(c, resp.CodeForbidden, "无权限访问")
	}
}

// AuditLog records mutating requests so they show up in the operation log page.
func (h *Handler) AuditLog(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		username := c.GetString("loginUsername")
		if u, ok := c.Get(ctxUser); ok {
			username = u.(*model.User).Username
		}
		h.db.Create(&model.OperationLog{
			Username: username,
			Action:   action,
			Method:   c.Request.Method,
			Path:     c.Request.URL.Path,
			IP:       c.ClientIP(),
			Success:  !c.GetBool("bizFailed"),
			Latency:  time.Since(start).Milliseconds(),
		})
	}
}

func currentUser(c *gin.Context) *model.User {
	return c.MustGet(ctxUser).(*model.User)
}

func fail(c *gin.Context, msg string) {
	c.Set("bizFailed", true)
	resp.Fail(c, resp.CodeError, msg)
}
