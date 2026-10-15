package handler

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Pirellis/Cursor_work/server/internal/model"
	"github.com/Pirellis/Cursor_work/server/internal/resp"
)

// ---------- 概览 ----------

func (h *Handler) Overview(c *gin.Context) {
	var users, active, roles, logsToday int64
	h.db.Model(&model.User{}).Count(&users)
	h.db.Model(&model.User{}).Where("status = ?", model.StatusEnabled).Count(&active)
	h.db.Model(&model.Role{}).Count(&roles)
	today := time.Now().Truncate(24 * time.Hour)
	h.db.Model(&model.OperationLog{}).Where("created_at >= ?", today).Count(&logsToday)

	type dayCount struct {
		Day   string `json:"day"`
		Count int64  `json:"count"`
	}
	trend := make([]dayCount, 0, 7)
	for i := 6; i >= 0; i-- {
		start := today.AddDate(0, 0, -i)
		var n int64
		h.db.Model(&model.OperationLog{}).
			Where("created_at >= ? AND created_at < ?", start, start.AddDate(0, 0, 1)).Count(&n)
		trend = append(trend, dayCount{Day: start.Format("01-02"), Count: n})
	}

	var recent []model.OperationLog
	h.db.Order("id desc").Limit(6).Find(&recent)

	resp.OK(c, gin.H{
		"userCount":   users,
		"activeCount": active,
		"roleCount":   roles,
		"logsToday":   logsToday,
		"trend":       trend,
		"recentLogs":  recent,
	})
}

// ---------- 用户 ----------

func (h *Handler) ListUsers(c *gin.Context) {
	page, size := paging(c)
	q := h.db.Model(&model.User{})
	if v := c.Query("username"); v != "" {
		q = q.Where("username LIKE ? OR nickname LIKE ?", "%"+v+"%", "%"+v+"%")
	}
	if v := c.Query("roleId"); v != "" {
		q = q.Where("role_id = ?", v)
	}
	if v := c.Query("status"); v != "" {
		q = q.Where("status = ?", v)
	}
	var total int64
	q.Count(&total)
	var list []model.User
	q.Preload("Role").Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list)
	resp.OK(c, resp.Page{List: list, Total: total})
}

type userReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	RoleID   uint   `json:"roleId"`
	Status   *int   `json:"status"`
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req userReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.RoleID == 0 {
		fail(c, "用户名和角色不能为空")
		return
	}
	if len(req.Password) < 5 {
		fail(c, "密码至少 5 位")
		return
	}
	var n int64
	h.db.Model(&model.User{}).Where("username = ?", req.Username).Count(&n)
	if n > 0 {
		fail(c, "用户名已存在")
		return
	}
	hash, _ := HashPassword(req.Password)
	status := model.StatusEnabled
	if req.Status != nil {
		status = *req.Status
	}
	u := model.User{
		Username: req.Username, Password: hash, Nickname: req.Nickname,
		Email: req.Email, Phone: req.Phone, RoleID: req.RoleID, Status: status,
	}
	if err := h.db.Create(&u).Error; err != nil {
		fail(c, "创建失败")
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
	if err := c.ShouldBindJSON(&req); err != nil || req.RoleID == 0 {
		fail(c, "角色不能为空")
		return
	}
	if id == currentUser(c).ID && req.Status != nil && *req.Status != model.StatusEnabled {
		fail(c, "不能禁用当前登录账号")
		return
	}
	updates := map[string]any{
		"nickname": req.Nickname, "email": req.Email, "phone": req.Phone, "role_id": req.RoleID,
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if err := h.db.Model(&model.User{ID: id}).Updates(updates).Error; err != nil {
		fail(c, "更新失败")
		return
	}
	resp.OK(c, nil)
}

func (h *Handler) ResetPassword(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Password) < 5 {
		fail(c, "密码至少 5 位")
		return
	}
	hash, _ := HashPassword(req.Password)
	h.db.Model(&model.User{ID: id}).Update("password", hash)
	resp.OK(c, nil)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if id == currentUser(c).ID {
		fail(c, "不能删除当前登录账号")
		return
	}
	h.db.Delete(&model.User{}, id)
	resp.OK(c, nil)
}

// ---------- 角色 ----------

type roleWithCount struct {
	model.Role
	UserCount int64 `json:"userCount"`
}

func (h *Handler) ListRoles(c *gin.Context) {
	page, size := paging(c)
	q := h.db.Model(&model.Role{})
	if v := c.Query("name"); v != "" {
		q = q.Where("name LIKE ? OR code LIKE ?", "%"+v+"%", "%"+v+"%")
	}
	var total int64
	q.Count(&total)
	var roles []model.Role
	q.Order("id asc").Offset((page - 1) * size).Limit(size).Find(&roles)
	list := make([]roleWithCount, len(roles))
	for i, r := range roles {
		list[i].Role = r
		h.db.Model(&model.User{}).Where("role_id = ?", r.ID).Count(&list[i].UserCount)
	}
	resp.OK(c, resp.Page{List: list, Total: total})
}

func (h *Handler) RoleOptions(c *gin.Context) {
	var roles []model.Role
	h.db.Where("status = ?", model.StatusEnabled).Order("id asc").Find(&roles)
	resp.OK(c, roles)
}

type roleReq struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Status      *int   `json:"status"`
}

func (h *Handler) CreateRole(c *gin.Context) {
	var req roleReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" || req.Code == "" {
		fail(c, "角色名称和标识不能为空")
		return
	}
	var n int64
	h.db.Model(&model.Role{}).Where("code = ?", req.Code).Count(&n)
	if n > 0 {
		fail(c, "角色标识已存在")
		return
	}
	status := model.StatusEnabled
	if req.Status != nil {
		status = *req.Status
	}
	r := model.Role{Name: req.Name, Code: req.Code, Description: req.Description, Status: status}
	h.db.Create(&r)
	resp.OK(c, r)
}

func (h *Handler) UpdateRole(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req roleReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		fail(c, "角色名称不能为空")
		return
	}
	var role model.Role
	if err := h.db.First(&role, id).Error; err != nil {
		fail(c, "角色不存在")
		return
	}
	updates := map[string]any{"name": req.Name, "description": req.Description}
	if req.Status != nil {
		if role.Code == "admin" && *req.Status != model.StatusEnabled {
			fail(c, "不能禁用超级管理员角色")
			return
		}
		updates["status"] = *req.Status
	}
	h.db.Model(&role).Updates(updates)
	resp.OK(c, nil)
}

func (h *Handler) DeleteRole(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var role model.Role
	if err := h.db.First(&role, id).Error; err != nil {
		fail(c, "角色不存在")
		return
	}
	if role.Code == "admin" {
		fail(c, "不能删除超级管理员角色")
		return
	}
	var n int64
	h.db.Model(&model.User{}).Where("role_id = ?", id).Count(&n)
	if n > 0 {
		fail(c, "该角色下仍有用户，无法删除")
		return
	}
	h.db.Delete(&role)
	resp.OK(c, nil)
}

// ---------- 日志 ----------

func (h *Handler) ListLogs(c *gin.Context) {
	page, size := paging(c)
	q := h.db.Model(&model.OperationLog{})
	if v := c.Query("username"); v != "" {
		q = q.Where("username LIKE ?", "%"+v+"%")
	}
	if v := c.Query("action"); v != "" {
		q = q.Where("action LIKE ?", "%"+v+"%")
	}
	var total int64
	q.Count(&total)
	var list []model.OperationLog
	q.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list)
	resp.OK(c, resp.Page{List: list, Total: total})
}
