package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/service"
	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// ListRoles 查询角色列表。
func (h *Handler) ListRoles(c *gin.Context) {
	list, err := h.services.Role.ListRoles()
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": list})
}

// CreateRole 创建角色。
func (h *Handler) CreateRole(c *gin.Context) {
	var req service.RoleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err := h.services.Role.CreateRole(&req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// UpdateRole 更新角色。
func (h *Handler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "角色ID格式错误")
		return
	}
	var req service.RoleUpdateRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err = h.services.Role.UpdateRole(id, &req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// DeleteRole 删除角色。
func (h *Handler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "角色ID格式错误")
		return
	}
	if err = h.services.Role.DeleteRole(id); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// AssignRoleMenus 分配角色菜单。
func (h *Handler) AssignRoleMenus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "角色ID格式错误")
		return
	}
	var req service.AssignMenusRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err = h.services.Role.AssignMenus(id, &req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}
