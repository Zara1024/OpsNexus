package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/repository"
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/service"
	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// ListHosts 查询主机列表。
func (h *Handler) ListHosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	groupID, _ := strconv.ParseInt(c.DefaultQuery("group_id", "0"), 10, 64)
	statusInt, _ := strconv.Atoi(c.DefaultQuery("status", "0"))
	filter := &repository.HostListFilter{
		Page:     page,
		PageSize: pageSize,
		GroupID:  groupID,
		Status:   int16(statusInt),
		Keyword:  c.Query("keyword"),
		InnerIP:  c.Query("inner_ip"),
		OuterIP:  c.Query("outer_ip"),
	}
	res, err := h.services.Host.ListHosts(filter)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, res)
}

// GetHost 查询主机详情。
func (h *Handler) GetHost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "主机ID格式错误")
		return
	}
	res, err := h.services.Host.GetHost(id)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, res)
}

// CreateHost 创建主机。
func (h *Handler) CreateHost(c *gin.Context) {
	var req service.HostCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err := h.services.Host.CreateHost(&req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// UpdateHost 更新主机。
func (h *Handler) UpdateHost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "主机ID格式错误")
		return
	}
	var req service.HostUpdateRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err = h.services.Host.UpdateHost(id, &req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// DeleteHost 删除主机。
func (h *Handler) DeleteHost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "主机ID格式错误")
		return
	}
	if err = h.services.Host.DeleteHost(id); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// TestHostSSH 测试 SSH 连接。
func (h *Handler) TestHostSSH(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "主机ID格式错误")
		return
	}
	if err = h.services.Host.TestSSHConnection(id); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"message": "SSH连接成功"})
}

// BatchOperateHosts 批量操作主机。
func (h *Handler) BatchOperateHosts(c *gin.Context) {
	var req service.HostBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err := h.services.Host.BatchOperateHosts(&req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}
