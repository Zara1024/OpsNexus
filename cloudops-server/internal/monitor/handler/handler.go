package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/monitor/service"
	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/middleware"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// Handler 聚合 monitor 模块处理器。
type Handler struct {
	services *service.Services
}

func New(services *service.Services) *Handler { return &Handler{services: services} }

func (h *Handler) RegisterRoutes(v1 *gin.RouterGroup, jwtMiddleware gin.HandlerFunc, permissionMiddleware func(string) gin.HandlerFunc, auditMiddleware gin.HandlerFunc) {
	monitor := v1.Group("/monitor")
	monitor.Use(jwtMiddleware)
	{
		rules := monitor.Group("/alert-rules")
		{
			rules.GET("", permissionMiddleware("monitor:rule:list"), h.ListAlertRules)
			rules.POST("", permissionMiddleware("monitor:rule:create"), auditMiddleware, h.CreateAlertRule)
			rules.PUT("/:id", permissionMiddleware("monitor:rule:update"), auditMiddleware, h.UpdateAlertRule)
			rules.DELETE("/:id", permissionMiddleware("monitor:rule:delete"), auditMiddleware, h.DeleteAlertRule)
			rules.PUT("/:id/toggle", permissionMiddleware("monitor:rule:toggle"), auditMiddleware, h.ToggleAlertRule)
		}

		alerts := monitor.Group("/alerts")
		{
			alerts.GET("", permissionMiddleware("monitor:alert:list"), h.ListActiveAlerts)
			alerts.GET("/history", permissionMiddleware("monitor:alert:history"), h.ListAlertHistory)
			alerts.PUT("/:id/ack", permissionMiddleware("monitor:alert:ack"), auditMiddleware, h.AckAlert)
			alerts.PUT("/:id/silence", permissionMiddleware("monitor:alert:silence"), auditMiddleware, h.SilenceAlert)
		}

		channels := monitor.Group("/channels")
		{
			channels.GET("", permissionMiddleware("monitor:channel:list"), h.ListChannels)
			channels.POST("", permissionMiddleware("monitor:channel:create"), auditMiddleware, h.CreateChannel)
			channels.POST("/:id/test", permissionMiddleware("monitor:channel:test"), auditMiddleware, h.TestChannel)
		}
	}

	_ = middleware.GetRequestID
}

func (h *Handler) ListAlertRules(c *gin.Context) {
	list, err := h.services.AlertRule.List()
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": list})
}

func (h *Handler) CreateAlertRule(c *gin.Context) {
	var req service.AlertRuleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err := h.services.AlertRule.Create(&req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) UpdateAlertRule(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req service.AlertRuleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err := h.services.AlertRule.Update(id, &req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) DeleteAlertRule(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.services.AlertRule.Delete(id); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) ToggleAlertRule(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err := h.services.AlertRule.Toggle(id, req.Enabled); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) ListActiveAlerts(c *gin.Context) {
	list, err := h.services.Alert.ListActive()
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": list})
}

func (h *Handler) ListAlertHistory(c *gin.Context) {
	list, err := h.services.Alert.ListHistory()
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": list})
}

func (h *Handler) AckAlert(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	claims, _ := middleware.GetJWTClaims(c)
	userID := int64(0)
	if claims != nil {
		userID = claims.UserID
	}
	if err := h.services.Alert.Ack(id, userID); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) SilenceAlert(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		Minutes int `json:"minutes"`
	}
	_ = c.ShouldBindJSON(&req)
	if err := h.services.Alert.Silence(id, req.Minutes); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) ListChannels(c *gin.Context) {
	list, err := h.services.Channel.List()
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": list})
}

func (h *Handler) CreateChannel(c *gin.Context) {
	var req service.ChannelCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	if err := h.services.Channel.Create(&req); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) TestChannel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.services.Channel.Test(id); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"message": "渠道测试成功"})
}
