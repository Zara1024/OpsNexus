package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// ListServices 查询 Service 列表。
func (h *Handler) ListServices(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	res, err := h.services.Network.ListServices(clusterID, c.Param("ns"))
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": res})
}

// ListIngresses 查询 Ingress 列表。
func (h *Handler) ListIngresses(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	res, err := h.services.Network.ListIngresses(clusterID, c.Param("ns"))
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": res})
}
