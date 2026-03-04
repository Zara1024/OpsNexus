package handler

import (
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"

	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// ListConfigMaps 查询 ConfigMap 列表。
func (h *Handler) ListConfigMaps(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	res, err := h.services.Config.ListConfigMaps(clusterID, c.Param("ns"))
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": res})
}

// SaveConfigMap 保存 ConfigMap。
func (h *Handler) SaveConfigMap(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	var req corev1.ConfigMap
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "ConfigMap参数不合法")
		return
	}
	if err = h.services.Config.CreateOrUpdateConfigMap(clusterID, c.Param("ns"), &req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// DeleteConfigMap 删除 ConfigMap。
func (h *Handler) DeleteConfigMap(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	if err = h.services.Config.DeleteConfigMap(clusterID, c.Param("ns"), c.Param("name")); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// ListSecrets 查询 Secret 列表。
func (h *Handler) ListSecrets(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	decode := c.DefaultQuery("decode", "0") == "1"
	res, err := h.services.Config.ListSecrets(clusterID, c.Param("ns"), decode)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": res})
}

// SaveSecret 保存 Secret。
func (h *Handler) SaveSecret(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	var req struct {
		Name string            `json:"name" binding:"required"`
		Type string            `json:"type"`
		Data map[string]string `json:"data"`
	}
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "Secret参数不合法")
		return
	}
	item := &corev1.Secret{Data: map[string][]byte{}, Type: corev1.SecretTypeOpaque}
	item.Name = req.Name
	if req.Type != "" {
		item.Type = corev1.SecretType(req.Type)
	}
	for k, v := range req.Data {
		if b, decodeErr := base64.StdEncoding.DecodeString(v); decodeErr == nil {
			item.Data[k] = b
			continue
		}
		item.Data[k] = []byte(v)
	}
	if err = h.services.Config.CreateOrUpdateSecret(clusterID, c.Param("ns"), item, true); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// DeleteSecret 删除 Secret。
func (h *Handler) DeleteSecret(c *gin.Context) {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, "集群ID格式错误")
		return
	}
	if err = h.services.Config.DeleteSecret(clusterID, c.Param("ns"), c.Param("name")); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}
