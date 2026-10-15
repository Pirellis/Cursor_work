package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/go-arco-admin/server/internal/config"
	mw "github.com/go-arco-admin/server/internal/middleware"
	"github.com/go-arco-admin/server/internal/model"
	"github.com/go-arco-admin/server/internal/resp"
)

type Handler struct {
	DB  *gorm.DB
	Cfg config.Config
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		resp.Fail(c, resp.CodeBadRequest, "用户名和密码不能为空")
		return
	}
	c.Set("loginUsername", req.Username)
	var u model.User
	if err := h.DB.Preload("Role").Where("username = ?", req.Username).First(&u).Error; err != nil ||
		!model.CheckPassword(u.PasswordHash, req.Password) {
		resp.Fail(c, resp.CodeError, "账号或者密码错误")
		return
	}
	if u.Status != 1 {
		resp.Fail(c, resp.CodeError, "账号已被禁用，请联系管理员")
		return
	}
	token, err := mw.SignToken(h.Cfg.JWTSecret, u.ID, h.Cfg.TokenTTL)
	if err != nil {
		resp.Fail(c, resp.CodeError, "签发凭证失败")
		return
	}
	now := time.Now()
	h.DB.Model(&u).Update("last_login_at", &now)
	resp.OK(c, gin.H{"token": token})
}

func (h *Handler) Logout(c *gin.Context) { resp.OK(c, nil) }

// Info returns the user in the shape of the Arco Pro `UserState` store.
func (h *Handler) Info(c *gin.Context) {
	u := mw.CurrentUser(c)
	resp.OK(c, gin.H{
		"id":               u.ID,
		"name":             displayName(u),
		"username":         u.Username,
		"avatar":           u.Avatar,
		"email":            u.Email,
		"phone":            u.Phone,
		"job":              u.Job,
		"jobName":          u.Job,
		"introduction":     u.Introduction,
		"registrationDate": u.CreatedAt.Format("2006-01-02 15:04:05"),
		"accountId":        strconv.FormatUint(uint64(u.ID), 10),
		"certification":    1,
		"role":             u.Role.Code,
		"roleName":         u.Role.Name,
	})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req struct {
		Nickname     string `json:"nickname"`
		Email        string `json:"email"`
		Phone        string `json:"phone"`
		Job          string `json:"job"`
		Introduction string `json:"introduction"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail(c, resp.CodeBadRequest, "参数错误")
		return
	}
	u := mw.CurrentUser(c)
	h.DB.Model(u).Select("nickname", "email", "phone", "job", "introduction").Updates(model.User{
		Nickname: req.Nickname, Email: req.Email, Phone: req.Phone, Job: req.Job, Introduction: req.Introduction,
	})
	resp.OK(c, nil)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.NewPassword) < 5 {
		resp.Fail(c, resp.CodeBadRequest, "新密码至少 5 位")
		return
	}
	u := mw.CurrentUser(c)
	if !model.CheckPassword(u.PasswordHash, req.OldPassword) {
		resp.Fail(c, resp.CodeError, "原密码错误")
		return
	}
	hash, _ := model.HashPassword(req.NewPassword)
	h.DB.Model(u).Update("password_hash", hash)
	resp.OK(c, nil)
}

func displayName(u *model.User) string {
	if u.Nickname != "" {
		return u.Nickname
	}
	return u.Username
}
