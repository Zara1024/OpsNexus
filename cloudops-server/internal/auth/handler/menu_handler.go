package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/service"
	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// ListMenus 查询菜单列表。
func (h *Handler) ListMenus(c *gin.Context) {
	list, err := h.services.Menu.ListMenus()
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": list})
}

// MenuTree 查询菜单树。
func (h *Handler) MenuTree(c *gin.Context) {
	tree, err := h.services.Menu.MenuTree()
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": tree})
}

// CreateMenu 创建菜单。
func (h *Handler) CreateMenu(c *gin.Context) {
	var req service.MenuCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err := h.services.Menu.CreateMenu(&req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// UpdateMenu 更新菜单。
func (h *Handler) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "菜单ID格式错误")
		return
	}
	var req service.MenuUpdateRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err = h.services.Menu.UpdateMenu(id, &req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// DeleteMenu 删除菜单。
func (h *Handler) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "菜单ID格式错误")
		return
	}
	if err = h.services.Menu.DeleteMenu(id); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}
