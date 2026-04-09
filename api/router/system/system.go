package system

import (
	"dodevops-api/api/system/controller"

	"github.com/gin-gonic/gin"
)

// RegisterSystemRoutes registers system-related routes.
func RegisterSystemRoutes(router *gin.RouterGroup) {
	ldapCtrl := controller.NewLDAPController()

	router.POST("/post/add", controller.CreateSysPost)
	router.GET("/post/list", controller.GetSysPostList)
	router.GET("/post/info", controller.GetSysPostById)
	router.PUT("/post/update", controller.UpdateSysPost)
	router.DELETE("/post/delete", controller.DeleteSysPostById)
	router.DELETE("/post/batch/delete", controller.BatchDeleteSysPost)
	router.PUT("/post/updateStatus", controller.UpdateSysPostStatus)
	router.GET("/post/vo/list", controller.QuerySysPostVoList)

	router.GET("/dept/list", controller.GetSysDeptList)
	router.POST("/dept/add", controller.CreateSysDept)
	router.GET("/dept/info", controller.GetSysDeptById)
	router.PUT("/dept/update", controller.UpdateSysDept)
	router.DELETE("/dept/delete", controller.DeleteSysDeptById)
	router.GET("/dept/vo/list", controller.QuerySysDeptVoList)
	router.GET("/dept/users", controller.GetDeptUsers)

	router.POST("/menu/add", controller.CreateSysMenu)
	router.GET("/menu/vo/list", controller.QuerySysMenuVoList)
	router.GET("/menu/info", controller.GetSysMenu)
	router.PUT("/menu/update", controller.UpdateSysMenu)
	router.DELETE("/menu/delete", controller.DeleteSysMenu)
	router.GET("/menu/list", controller.GetSysMenuList)

	router.POST("/role/add", controller.CreateSysRole)
	router.GET("/role/info", controller.GetSysRoleById)
	router.PUT("/role/update", controller.UpdateSysRole)
	router.DELETE("/role/delete", controller.DeleteSysRoleById)
	router.PUT("/role/updateStatus", controller.UpdateSysRoleStatus)
	router.GET("/role/list", controller.GetSysRoleList)
	router.GET("/role/vo/list", controller.QuerySysRoleVoList)
	router.GET("/role/vo/idList", controller.QueryRoleMenuIdList)
	router.PUT("/role/assignPermissions", controller.AssignPermissions)

	router.POST("/admin/add", controller.CreateSysAdmin)
	router.GET("/admin/info", controller.GetSysAdminInfo)
	router.PUT("/admin/update", controller.UpdateSysAdmin)
	router.DELETE("/admin/delete", controller.DeleteSysAdminById)
	router.PUT("/admin/updateStatus", controller.UpdateSysAdminStatus)
	router.PUT("/admin/updatePassword", controller.ResetSysAdminPassword)
	router.GET("/admin/list", controller.GetSysAdminList)
	router.POST("/upload", controller.Upload)
	router.PUT("/admin/updatePersonal", controller.UpdatePersonal)
	router.PUT("/admin/updatePersonalPassword", controller.UpdatePersonalPassword)

	router.GET("/sysLoginInfo/list", controller.GetSysLoginInfoList)
	router.DELETE("/sysLoginInfo/batch/delete", controller.BatchDeleteSysLoginInfo)
	router.DELETE("/sysLoginInfo/delete", controller.DeleteSysLoginInfoById)
	router.DELETE("/sysLoginInfo/clean", controller.CleanSysLoginInfo)
	router.GET("/sysOperationLog/list", controller.GetSysOperationLogList)
	router.DELETE("/sysOperationLog/delete", controller.DeleteSysOperationLogById)
	router.DELETE("/sysOperationLog/batch/delete", controller.BatchDeleteSysOperationLog)
	router.DELETE("/sysOperationLog/clean", controller.CleanSysOperationLog)
	router.GET("/terminalAudit/summary", controller.GetTerminalAuditSummary)
	router.GET("/terminalAudit/list", controller.GetTerminalAuditSessionList)
	router.GET("/terminalAudit/session/:sessionId", controller.GetTerminalAuditSessionDetail)
	router.GET("/terminalAudit/session/:sessionId/playback", controller.GetTerminalAuditSessionPlayback)
	router.GET("/terminalAudit/session/:sessionId/download", controller.DownloadTerminalAuditRecording)

	router.GET("/system/ldap/config", ldapCtrl.GetConfig)
	router.PUT("/system/ldap/config", ldapCtrl.UpdateConfig)
	router.POST("/system/ldap/test", ldapCtrl.TestConfig)
}
