package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/service"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/middleware"
)

// Handler 聚合auth模块HTTP处理器。
type Handler struct {
	services *service.Services
}

// New 创建auth模块处理器。
func New(services *service.Services) *Handler {
	return &Handler{services: services}
}

// RegisterRoutes 注册auth模块路由。
func (h *Handler) RegisterRoutes(v1 *gin.RouterGroup, jwtMiddleware gin.HandlerFunc, permissionMiddleware func(string) gin.HandlerFunc, auditMiddleware gin.HandlerFunc) {
	authGroup := v1.Group("/auth")
	{
		authGroup.GET("/captcha", h.GetCaptcha)
		authGroup.POST("/login", h.Login)
		authGroup.POST("/refresh", h.Refresh)
		authGroup.POST("/logout", jwtMiddleware, h.Logout)
		authGroup.GET("/userinfo", jwtMiddleware, h.UserInfo)
	}

	systemGroup := v1.Group("/system")
	systemGroup.Use(jwtMiddleware)
	{
		users := systemGroup.Group("/users")
		{
			users.GET("", permissionMiddleware("system:user:list"), h.ListUsers)
			users.POST("", permissionMiddleware("system:user:create"), auditMiddleware, h.CreateUser)
			users.PUT("/:id", permissionMiddleware("system:user:update"), auditMiddleware, h.UpdateUser)
			users.DELETE("/:id", permissionMiddleware("system:user:delete"), auditMiddleware, h.DeleteUser)
			users.PUT("/:id/password", permissionMiddleware("system:user:password"), auditMiddleware, h.ChangePassword)
			users.PUT("/:id/roles", permissionMiddleware("system:user:assign_roles"), auditMiddleware, h.AssignRoles)
		}

		roles := systemGroup.Group("/roles")
		{
			roles.GET("", permissionMiddleware("system:role:list"), h.ListRoles)
			roles.POST("", permissionMiddleware("system:role:create"), auditMiddleware, h.CreateRole)
			roles.PUT("/:id", permissionMiddleware("system:role:update"), auditMiddleware, h.UpdateRole)
			roles.DELETE("/:id", permissionMiddleware("system:role:delete"), auditMiddleware, h.DeleteRole)
			roles.PUT("/:id/menus", permissionMiddleware("system:role:assign_menus"), auditMiddleware, h.AssignRoleMenus)
		}

		menus := systemGroup.Group("/menus")
		{
			menus.GET("", permissionMiddleware("system:menu:list"), h.ListMenus)
			menus.GET("/tree", permissionMiddleware("system:menu:list"), h.MenuTree)
			menus.POST("", permissionMiddleware("system:menu:create"), auditMiddleware, h.CreateMenu)
			menus.PUT("/:id", permissionMiddleware("system:menu:update"), auditMiddleware, h.UpdateMenu)
			menus.DELETE("/:id", permissionMiddleware("system:menu:delete"), auditMiddleware, h.DeleteMenu)
		}

		departments := systemGroup.Group("/departments")
		{
			departments.GET("", permissionMiddleware("system:department:list"), h.ListDepartments)
			departments.POST("", permissionMiddleware("system:department:create"), auditMiddleware, h.CreateDepartment)
			departments.PUT("/:id", permissionMiddleware("system:department:update"), auditMiddleware, h.UpdateDepartment)
			departments.DELETE("/:id", permissionMiddleware("system:department:delete"), auditMiddleware, h.DeleteDepartment)
		}
	}

	_ = middleware.GetRequestID
}
