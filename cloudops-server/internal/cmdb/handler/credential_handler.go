package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/service"
	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// ListCredentials 查询凭据列表。
func (h *Handler) ListCredentials(c *gin.Context) {
	list, err := h.services.Credential.ListCredentials()
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": list})
}

// CreateCredential 创建凭据。
func (h *Handler) CreateCredential(c *gin.Context) {
	var req service.CredentialCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err := h.services.Credential.CreateCredential(&req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// UpdateCredential 更新凭据。
func (h *Handler) UpdateCredential(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "凭据ID格式错误")
		return
	}
	var req service.CredentialUpdateRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err = h.services.Credential.UpdateCredential(id, &req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// DeleteCredential 删除凭据。
func (h *Handler) DeleteCredential(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "凭据ID格式错误")
		return
	}
	if err = h.services.Credential.DeleteCredential(id); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}
