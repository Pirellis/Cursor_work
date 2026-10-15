package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/Pirellis/Cursor_work/server/internal/auth"
	"github.com/Pirellis/Cursor_work/server/internal/model"
	"github.com/Pirellis/Cursor_work/server/internal/resp"
)

type Handler struct {
	db  *gorm.DB
	jwt *auth.Manager
}

func New(db *gorm.DB, jwt *auth.Manager) *Handler {
	return &Handler{db: db, jwt: jwt}
}

func (h *Handler) Register(r *gin.Engine) {
	api := r.Group("/api")

	api.POST("/user/login", h.AuditLog("登录"), h.Login)

	authed := api.Group("", h.Auth())
	authed.POST("/user/logout", h.AuditLog("退出登录"), h.Logout)
	authed.POST("/user/info", h.Info)
	authed.PUT("/user/profile", h.AuditLog("修改个人资料"), h.UpdateProfile)
	authed.PUT("/user/password", h.AuditLog("修改密码"), h.ChangePassword)

	sys := authed.Group("/system", h.RequireRole("admin"))
	sys.GET("/overview", h.Overview)

	sys.GET("/users", h.ListUsers)
	sys.POST("/users", h.AuditLog("新增用户"), h.CreateUser)
	sys.PUT("/users/:id", h.AuditLog("编辑用户"), h.UpdateUser)
	sys.PUT("/users/:id/password", h.AuditLog("重置密码"), h.ResetPassword)
	sys.DELETE("/users/:id", h.AuditLog("删除用户"), h.DeleteUser)

	sys.GET("/roles", h.ListRoles)
	sys.GET("/roles/options", h.RoleOptions)
	sys.POST("/roles", h.AuditLog("新增角色"), h.CreateRole)
	sys.PUT("/roles/:id", h.AuditLog("编辑角色"), h.UpdateRole)
	sys.DELETE("/roles/:id", h.AuditLog("删除角色"), h.DeleteRole)

	sys.GET("/logs", h.ListLogs)
}

func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

func paging(c *gin.Context) (page, size int) {
	page, _ = strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ = strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	return
}

func idParam(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		fail(c, "非法的 ID")
		return 0, false
	}
	return uint(id), true
}

// ---------- 当前用户 ----------

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		fail(c, "用户名和密码不能为空")
		return
	}
	c.Set("loginUsername", req.Username)
	var user model.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil ||
		bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		fail(c, "账号或者密码错误")
		return
	}
	if user.Status != model.StatusEnabled {
		fail(c, "账号已被禁用")
		return
	}
	token, err := h.jwt.Issue(user.ID, user.Username)
	if err != nil {
		fail(c, "签发凭证失败")
		return
	}
	now := time.Now()
	h.db.Model(&user).Update("last_login_at", &now)
	resp.OK(c, gin.H{"token": token})
}

func (h *Handler) Logout(c *gin.Context) {
	resp.OK(c, nil)
}

// Info returns the shape consumed by the Pro template's user store.
func (h *Handler) Info(c *gin.Context) {
	u := currentUser(c)
	name := u.Nickname
	if name == "" {
		name = u.Username
	}
	resp.OK(c, gin.H{
		"id":               u.ID,
		"name":             name,
		"username":         u.Username,
		"avatar":           u.Avatar,
		"email":            u.Email,
		"phone":            u.Phone,
		"role":             u.Role.Code,
		"roleName":         u.Role.Name,
		"jobName":          u.Role.Name,
		"organizationName": "Go Arco Admin",
		"locationName":     "云端",
		"accountId":        strconv.Itoa(int(u.ID)),
		"registrationDate": u.CreatedAt.Format("2006-01-02 15:04:05"),
		"certification":    1,
	})
}

type profileReq struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req profileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "参数错误")
		return
	}
	u := currentUser(c)
	h.db.Model(u).Updates(map[string]any{
		"nickname": req.Nickname, "email": req.Email, "phone": req.Phone, "avatar": req.Avatar,
	})
	resp.OK(c, nil)
}

type passwordReq struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var req passwordReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.NewPassword) < 5 {
		fail(c, "新密码至少 5 位")
		return
	}
	u := currentUser(c)
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.OldPassword)) != nil {
		fail(c, "原密码不正确")
		return
	}
	hash, _ := HashPassword(req.NewPassword)
	h.db.Model(u).Update("password", hash)
	resp.OK(c, nil)
}
