package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/audit/service"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// Handler 审计模块处理器。
type Handler struct {
	services *service.Services
}

func New(services *service.Services) *Handler { return &Handler{services: services} }

func (h *Handler) RegisterRoutes(v1 *gin.RouterGroup, jwtMiddleware gin.HandlerFunc, permissionMiddleware func(string) gin.HandlerFunc) {
	auditGroup := v1.Group("/audit")
	auditGroup.Use(jwtMiddleware)
	{
		auditGroup.GET("/logs", permissionMiddleware("audit:log:list"), h.ListLogs)
		auditGroup.GET("/logs/export", permissionMiddleware("audit:log:export"), h.ExportLogs)
		auditGroup.GET("/login-logs", permissionMiddleware("audit:login:list"), h.ListLoginLogs)
	}
}

func (h *Handler) ListLogs(c *gin.Context) {
	var query service.ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		handleAppError(c, err)
		return
	}
	list, total, err := h.services.Log.List(query)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

func (h *Handler) ExportLogs(c *gin.Context) {
	var query service.ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		handleAppError(c, err)
		return
	}
	csvContent, err := h.services.Log.ExportCSV(query)
	if err != nil {
		handleAppError(c, err)
		return
	}
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=audit_logs.csv")
	c.String(http.StatusOK, csvContent)
}

func (h *Handler) ListLoginLogs(c *gin.Context) {
	var query service.ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		handleAppError(c, err)
		return
	}
	list, total, err := h.services.Log.ListLoginLogs(query)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}
