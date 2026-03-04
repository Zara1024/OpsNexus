package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/k8s/service"
	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/middleware"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// ListClusters 查询集群列表。
func (h *Handler) ListClusters(c *gin.Context) {
	res, err := h.services.Cluster.ListClusters()
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": res})
}

// CreateCluster 创建集群。
func (h *Handler) CreateCluster(c *gin.Context) {
	var req service.ClusterCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	createdBy := int64(0)
	if claims, ok := middleware.GetJWTClaims(c); ok {
		createdBy = claims.UserID
	}
	if err := h.services.Cluster.CreateCluster(&req, createdBy); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// GetCluster 查询集群详情。
func (h *Handler) GetCluster(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	res, err := h.services.Cluster.GetClusterDetail(id)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, res)
}

// UpdateCluster 更新集群。
func (h *Handler) UpdateCluster(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	var req service.ClusterUpdateRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err = h.services.Cluster.UpdateCluster(id, &req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// DeleteCluster 删除集群。
func (h *Handler) DeleteCluster(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	if err = h.services.Cluster.DeleteCluster(id); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// TestCluster 测试集群连接。
func (h *Handler) TestCluster(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	if err = h.services.Cluster.TestCluster(id); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"message": "连接成功"})
}
