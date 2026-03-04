package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/service"
)

// Handler 聚合 cmdb 模块处理器。
type Handler struct {
	services *service.Services
}

// New 创建 cmdb 模块处理器。
func New(services *service.Services) *Handler {
	return &Handler{services: services}
}

// RegisterRoutes 注册 cmdb HTTP 路由。
func (h *Handler) RegisterRoutes(v1 *gin.RouterGroup, jwtMiddleware gin.HandlerFunc, permissionMiddleware func(string) gin.HandlerFunc, auditMiddleware gin.HandlerFunc) {
	cmdb := v1.Group("/cmdb")
	cmdb.Use(jwtMiddleware)
	{
		hosts := cmdb.Group("/hosts")
		{
			hosts.GET("", permissionMiddleware("cmdb:host:list"), h.ListHosts)
			hosts.GET("/:id", permissionMiddleware("cmdb:host:detail"), h.GetHost)
			hosts.POST("", permissionMiddleware("cmdb:host:create"), auditMiddleware, h.CreateHost)
			hosts.PUT("/:id", permissionMiddleware("cmdb:host:update"), auditMiddleware, h.UpdateHost)
			hosts.DELETE("/:id", permissionMiddleware("cmdb:host:delete"), auditMiddleware, h.DeleteHost)
			hosts.POST("/:id/test", permissionMiddleware("cmdb:host:test"), h.TestHostSSH)
			hosts.POST("/batch", permissionMiddleware("cmdb:host:batch"), auditMiddleware, h.BatchOperateHosts)
		}

		groups := cmdb.Group("/groups")
		{
			groups.GET("", permissionMiddleware("cmdb:group:list"), h.ListGroups)
			groups.POST("", permissionMiddleware("cmdb:group:create"), auditMiddleware, h.CreateGroup)
			groups.PUT("/:id", permissionMiddleware("cmdb:group:update"), auditMiddleware, h.UpdateGroup)
			groups.DELETE("/:id", permissionMiddleware("cmdb:group:delete"), auditMiddleware, h.DeleteGroup)
		}

		credentials := cmdb.Group("/credentials")
		{
			credentials.GET("", permissionMiddleware("cmdb:credential:list"), h.ListCredentials)
			credentials.POST("", permissionMiddleware("cmdb:credential:create"), auditMiddleware, h.CreateCredential)
			credentials.PUT("/:id", permissionMiddleware("cmdb:credential:update"), auditMiddleware, h.UpdateCredential)
			credentials.DELETE("/:id", permissionMiddleware("cmdb:credential:delete"), auditMiddleware, h.DeleteCredential)
		}

		cmdb.GET("/ssh-records", permissionMiddleware("cmdb:ssh_record:list"), h.ListSSHRecords)
	}
}

// RegisterWSRoutes 注册 cmdb WebSocket 路由。
func (h *Handler) RegisterWSRoutes(r *gin.Engine, jwtMiddleware gin.HandlerFunc) {
	r.GET("/ws/cmdb/terminal/:host_id", jwtMiddleware, h.WebTerminal)
}
