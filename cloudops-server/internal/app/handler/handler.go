package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/app/service"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// Handler app 模块处理器。
type Handler struct {
	services *service.Services
}

func New(services *service.Services) *Handler { return &Handler{services: services} }

func (h *Handler) RegisterRoutes(v1 *gin.RouterGroup, jwtMiddleware gin.HandlerFunc, permissionMiddleware func(string) gin.HandlerFunc) {
	dashboard := v1.Group("/dashboard")
	dashboard.Use(jwtMiddleware)
	{
		dashboard.GET("/overview", permissionMiddleware("dashboard:view"), h.GetDashboardOverview)
		dashboard.GET("/trends", permissionMiddleware("dashboard:view"), h.GetDashboardTrends)
	}
}

func (h *Handler) GetDashboardOverview(c *gin.Context) {
	data, err := h.services.Dashboard.Overview()
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetDashboardTrends(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	data, err := h.services.Dashboard.Trends(days)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, data)
}
