package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// ListNodes 查询节点列表。
func (h *Handler) ListNodes(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	res, err := h.services.Node.ListNodes(id)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": res})
}

// CordonNode 封锁节点。
func (h *Handler) CordonNode(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	if err = h.services.Node.CordonNode(clusterID, c.Param("name")); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// UncordonNode 解封节点。
func (h *Handler) UncordonNode(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	if err = h.services.Node.UncordonNode(clusterID, c.Param("name")); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// DrainNode 排水节点。
func (h *Handler) DrainNode(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	if err = h.services.Node.DrainNode(clusterID, c.Param("name")); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}
