package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	mw "github.com/go-arco-admin/server/internal/middleware"
	"github.com/go-arco-admin/server/internal/model"
	"github.com/go-arco-admin/server/internal/resp"
)

func paging(c *gin.Context) (int, int) {
	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if current < 1 {
		current = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	return (current - 1) * size, size
}

func idParam(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.Fail(c, resp.CodeBadRequest, "无效的 ID")
		return 0, false
	}
	return uint(id), true
}

func like(s string) string { return "%" + strings.TrimSpace(s) + "%" }

// ---------- users ----------

type userReq struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Nickname     string `json:"nickname"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Job          string `json:"job"`
	Introduction string `json:"introduction"`
	Status       *int   `json:"status"`
	RoleID       uint   `json:"roleId"`
}

func (h *Handler) ListUsers(c *gin.Context) {
	q := h.DB.Model(&model.User{})
	if kw := c.Query("keyword"); kw != "" {
		q = q.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ? OR phone LIKE ?", like(kw), like(kw), like(kw), like(kw))
	}
	if s := c.Query("status"); s != "" {
		q = q.Where("status = ?", s)
	}
	if r := c.Query("roleId"); r != "" {
		q = q.Where("role_id = ?", r)
	}
	var page resp.Page[model.User]
	q.Count(&page.Total)
	offset, size := paging(c)
	q.Preload("Role").Order("id DESC").Offset(offset).Limit(size).Find(&page.List)
	resp.OK(c, page)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req userReq
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Username) == "" || len(req.Password) < 5 {
		resp.Fail(c, resp.CodeBadRequest, "用户名必填，密码至少 5 位")
		return
	}
	if !h.roleExists(req.RoleID) {
		resp.Fail(c, resp.CodeBadRequest, "请选择有效的角色")
		return
	}
	var n int64
	h.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&n)
	if n > 0 {
		resp.Fail(c, resp.CodeError, "用户名已存在")
		return
	}
	hash, _ := model.HashPassword(req.Password)
	u := model.User{Username: strings.TrimSpace(req.Username), PasswordHash: hash, Nickname: req.Nickname, Email: req.Email,
		Phone: req.Phone, Job: req.Job, Introduction: req.Introduction, RoleID: req.RoleID, Status: 1, Avatar: model.DefaultAvatar}
	if req.Status != nil {
		u.Status = *req.Status
	}
	if err := h.DB.Create(&u).Error; err != nil {
		resp.Fail(c, resp.CodeError, "创建失败")
		return
	}
	resp.OK(c, u)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req userReq
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail(c, resp.CodeBadRequest, "参数错误")
		return
	}
	var u model.User
	if err := h.DB.First(&u, id).Error; err != nil {
		resp.Fail(c, resp.CodeError, "用户不存在")
		return
	}
	if !h.roleExists(req.RoleID) {
		resp.Fail(c, resp.CodeBadRequest, "请选择有效的角色")
		return
	}
	self := mw.CurrentUser(c)
	if self.ID == u.ID && ((req.Status != nil && *req.Status != 1) || req.RoleID != u.RoleID) {
		resp.Fail(c, resp.CodeError, "不能禁用自己或修改自己的角色")
		return
	}
	updates := map[string]any{"nickname": req.Nickname, "email": req.Email, "phone": req.Phone, "job": req.Job,
		"introduction": req.Introduction, "role_id": req.RoleID}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Password != "" {
		if len(req.Password) < 5 {
			resp.Fail(c, resp.CodeBadRequest, "密码至少 5 位")
			return
		}
		updates["password_hash"], _ = model.HashPassword(req.Password)
	}
	h.DB.Model(&u).Updates(updates)
	resp.OK(c, nil)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if mw.CurrentUser(c).ID == id {
		resp.Fail(c, resp.CodeError, "不能删除当前登录账号")
		return
	}
	if h.DB.Delete(&model.User{}, id).RowsAffected == 0 {
		resp.Fail(c, resp.CodeError, "用户不存在")
		return
	}
	resp.OK(c, nil)
}

func (h *Handler) roleExists(id uint) bool {
	var n int64
	h.DB.Model(&model.Role{}).Where("id = ?", id).Count(&n)
	return n > 0
}

// ---------- roles ----------

type roleReq struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      *int   `json:"status"`
}

func (h *Handler) ListRoles(c *gin.Context) {
	var roles []model.Role
	q := h.DB.Order("id ASC")
	if kw := c.Query("keyword"); kw != "" {
		q = q.Where("code LIKE ? OR name LIKE ?", like(kw), like(kw))
	}
	q.Find(&roles)
	type row struct {
		RoleID uint
		N      int64
	}
	var counts []row
	h.DB.Model(&model.User{}).Select("role_id, count(*) as n").Group("role_id").Scan(&counts)
	byRole := map[uint]int64{}
	for _, r := range counts {
		byRole[r.RoleID] = r.N
	}
	for i := range roles {
		roles[i].UserCount = byRole[roles[i].ID]
	}
	resp.OK(c, roles)
}

func (h *Handler) CreateRole(c *gin.Context) {
	var req roleReq
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.Name) == "" {
		resp.Fail(c, resp.CodeBadRequest, "角色编码与名称必填")
		return
	}
	r := model.Role{Code: strings.TrimSpace(req.Code), Name: req.Name, Description: req.Description, Status: 1}
	if req.Status != nil {
		r.Status = *req.Status
	}
	if err := h.DB.Create(&r).Error; err != nil {
		resp.Fail(c, resp.CodeError, "角色编码已存在")
		return
	}
	resp.OK(c, r)
}

func (h *Handler) UpdateRole(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req roleReq
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		resp.Fail(c, resp.CodeBadRequest, "角色名称必填")
		return
	}
	var r model.Role
	if err := h.DB.First(&r, id).Error; err != nil {
		resp.Fail(c, resp.CodeError, "角色不存在")
		return
	}
	updates := map[string]any{"name": req.Name, "description": req.Description}
	if !r.Builtin {
		if req.Code != "" {
			updates["code"] = strings.TrimSpace(req.Code)
		}
		if req.Status != nil {
			updates["status"] = *req.Status
		}
	}
	if err := h.DB.Model(&r).Updates(updates).Error; err != nil {
		resp.Fail(c, resp.CodeError, "角色编码已存在")
		return
	}
	resp.OK(c, nil)
}

func (h *Handler) DeleteRole(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var r model.Role
	if err := h.DB.First(&r, id).Error; err != nil {
		resp.Fail(c, resp.CodeError, "角色不存在")
		return
	}
	if r.Builtin {
		resp.Fail(c, resp.CodeError, "内置角色不可删除")
		return
	}
	var n int64
	h.DB.Model(&model.User{}).Where("role_id = ?", id).Count(&n)
	if n > 0 {
		resp.Fail(c, resp.CodeError, "该角色下仍有用户，无法删除")
		return
	}
	h.DB.Delete(&r)
	resp.OK(c, nil)
}

// ---------- logs & overview ----------

func (h *Handler) ListLogs(c *gin.Context) {
	q := h.DB.Model(&model.OperationLog{})
	if kw := c.Query("keyword"); kw != "" {
		q = q.Where("username LIKE ? OR path LIKE ?", like(kw), like(kw))
	}
	if m := c.Query("method"); m != "" {
		q = q.Where("method = ?", m)
	}
	switch c.Query("result") {
	case "success":
		q = q.Where("code = ?", resp.CodeOK)
	case "fail":
		q = q.Where("code <> ?", resp.CodeOK)
	}
	var page resp.Page[model.OperationLog]
	q.Count(&page.Total)
	offset, size := paging(c)
	q.Order("id DESC").Offset(offset).Limit(size).Find(&page.List)
	resp.OK(c, page)
}

func (h *Handler) Overview(c *gin.Context) {
	var users, active, roles, todayOps int64
	h.DB.Model(&model.User{}).Count(&users)
	h.DB.Model(&model.User{}).Where("status = 1").Count(&active)
	h.DB.Model(&model.Role{}).Count(&roles)
	today := time.Now().Truncate(24 * time.Hour)
	h.DB.Model(&model.OperationLog{}).Where("created_at >= ?", today).Count(&todayOps)
	resp.OK(c, gin.H{"users": users, "activeUsers": active, "roles": roles, "todayOperations": todayOps})
}
