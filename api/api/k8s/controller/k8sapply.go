package controller

import (
	"strconv"

	"dodevops-api/api/k8s/model"
	"dodevops-api/api/k8s/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type K8sApplyController struct {
	service *service.K8sApplyService
}

func NewK8sApplyController(db *gorm.DB) *K8sApplyController {
	return &K8sApplyController{
		service: service.NewK8sApplyService(db),
	}
}

// ApplyManifest applies Kubernetes manifests with dry-run and server-side apply support.
// @Summary Apply Kubernetes manifest
// @Description Apply one or more Kubernetes manifests to the target cluster.
// @Tags K8s YAML Management
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer user token"
// @Param id path int true "Cluster ID"
// @Param request body model.ApplyManifestRequest true "Apply manifest request"
// @Success 200 {object} result.Result{data=model.ApplyManifestResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/apply [post]
func (ctrl *K8sApplyController) ApplyManifest(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	var req model.ApplyManifestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctrl.service.ApplyManifest(c, uint(clusterID), &req)
}
