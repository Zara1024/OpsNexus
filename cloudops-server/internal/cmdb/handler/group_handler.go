package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/service"
	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// ListGroups 查询分组树。
func (h *Handler) ListGroups(c *gin.Context) {
	list, err := h.services.Group.ListGroupTree()
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": list})
}

// CreateGroup 创建分组。
func (h *Handler) CreateGroup(c *gin.Context) {
	var req service.GroupCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err := h.services.Group.CreateGroup(&req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// UpdateGroup 更新分组。
func (h *Handler) UpdateGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "分组ID格式错误")
		return
	}
	var req service.GroupUpdateRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err = h.services.Group.UpdateGroup(id, &req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// DeleteGroup 删除分组。
func (h *Handler) DeleteGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "分组ID格式错误")
		return
	}
	if err = h.services.Group.DeleteGroup(id); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}
