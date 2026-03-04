package handler

import (
	"bufio"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/k8s/service"
	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// ListDeployments 查询 Deployment 列表。
func (h *Handler) ListDeployments(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	res, err := h.services.Workload.ListDeployments(clusterID, c.Param("ns"))
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": res})
}

// GetDeployment 查询 Deployment 详情。
func (h *Handler) GetDeployment(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	res, err := h.services.Workload.GetDeployment(clusterID, c.Param("ns"), c.Param("name"))
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, res)
}

// CreateDeployment 创建 Deployment。
func (h *Handler) CreateDeployment(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	var req appsv1.Deployment
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "Deployment参数不合法")
		return
	}
	if err = h.services.Workload.CreateDeployment(clusterID, c.Param("ns"), &req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// UpdateDeployment 更新 Deployment。
func (h *Handler) UpdateDeployment(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	name := c.Param("name")
	var req appsv1.Deployment
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "Deployment参数不合法")
		return
	}
	if err = h.services.Workload.UpdateDeployment(clusterID, c.Param("ns"), name, &req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// DeleteDeployment 删除 Deployment。
func (h *Handler) DeleteDeployment(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	if err = h.services.Workload.DeleteDeployment(clusterID, c.Param("ns"), c.Param("name")); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// ScaleDeployment 伸缩 Deployment。
func (h *Handler) ScaleDeployment(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	var req service.DeploymentScaleRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "伸缩参数不合法")
		return
	}
	if err = h.services.Workload.ScaleDeployment(clusterID, c.Param("ns"), c.Param("name"), req.Replicas); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// RestartDeployment 重启 Deployment。
func (h *Handler) RestartDeployment(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	if err = h.services.Workload.RestartDeployment(clusterID, c.Param("ns"), c.Param("name")); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// ListPods 查询 Pod 列表。
func (h *Handler) ListPods(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	res, err := h.services.Workload.ListPods(clusterID, c.Param("ns"))
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": res})
}

// GetPod 查询 Pod 详情。
func (h *Handler) GetPod(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	res, err := h.services.Workload.GetPod(clusterID, c.Param("ns"), c.Param("pod"))
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, res)
}

// DeletePod 删除 Pod。
func (h *Handler) DeletePod(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	if err = h.services.Workload.DeletePod(clusterID, c.Param("ns"), c.Param("pod")); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// GetPodLogs 获取 Pod 日志（支持 SSE）。
func (h *Handler) GetPodLogs(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	container := c.Query("container")
	tailLines, _ := strconv.ParseInt(c.DefaultQuery("tail_lines", "200"), 10, 64)
	logs, err := h.services.Workload.GetPodLogs(clusterID, c.Param("ns"), c.Param("pod"), container, tailLines)
	if err != nil {
		handleAppError(c, err)
		return
	}
	if c.DefaultQuery("stream", "0") == "1" {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		scanner := bufio.NewScanner(strings.NewReader(logs))
		for scanner.Scan() {
			_, _ = c.Writer.WriteString(service.BuildPodLogSSEEvent(scanner.Text()))
			c.Writer.Flush()
		}
		return
	}
	response.Success(c, gin.H{"logs": logs})
}
