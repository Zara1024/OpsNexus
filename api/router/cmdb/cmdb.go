package cmdb

import (
	"dodevops-api/api/cmdb/controller"
	"dodevops-api/api/cmdb/service"

	"github.com/gin-gonic/gin"
)

// RegisterCmdbRoutes registers CMDB routes.
func RegisterCmdbRoutes(router *gin.RouterGroup) {
	sqlWorkOrderCtrl := controller.NewSQLWorkOrderController()
	deviceCtrl := controller.NewCmdbDeviceController()

	// Group management
	router.POST("/cmdb/groupadd", controller.CreateCmdbGroup)
	router.GET("/cmdb/grouplist", controller.GetAllCmdbGroups)
	router.GET("/cmdb/grouplistwithhosts", controller.GetAllCmdbGroupsWithHosts)
	router.PUT("/cmdb/groupupdate", controller.UpdateCmdbGroup)
	router.DELETE("/cmdb/groupdelete", controller.DeleteCmdbGroup)
	router.GET("/cmdb/groupbyname", controller.GetCmdbGroupByName)

	// Host management
	router.POST("/cmdb/hostcreate", controller.NewCmdbHostController().CreateCmdbHost)
	router.PUT("/cmdb/hostupdate", controller.NewCmdbHostController().UpdateCmdbHost)
	router.DELETE("/cmdb/hostdelete", controller.NewCmdbHostController().DeleteCmdbHost)
	router.GET("/cmdb/hostlist", controller.NewCmdbHostController().GetCmdbHostListWithPage)
	router.GET("/cmdb/hostinfo", controller.NewCmdbHostController().GetCmdbHostById)
	router.GET("/cmdb/hostgroup", controller.NewCmdbHostController().GetCmdbHostsByGroupId)
	router.GET("/cmdb/hostbyname", controller.NewCmdbHostController().GetCmdbHostsByHostNameLike)
	router.GET("/cmdb/hostbyip", controller.NewCmdbHostController().GetCmdbHostsByIP)
	router.GET("/cmdb/hostbystatus", controller.NewCmdbHostController().GetCmdbHostsByStatus)
	router.POST("/cmdb/hostimport", controller.NewCmdbHostController().ImportHostsFromExcel)
	router.GET("/cmdb/hosttemplate", controller.NewCmdbHostController().DownloadHostTemplate)
	router.POST("/cmdb/hostsync", controller.NewCmdbHostController().SyncHostInfo)
	router.POST("/cmdb/hostconnectivity", controller.NewCmdbHostController().BatchTestHostConnectivity)

	// Network device management
	router.POST("/cmdb/device", deviceCtrl.CreateDevice)
	router.PUT("/cmdb/device", deviceCtrl.UpdateDevice)
	router.DELETE("/cmdb/device", deviceCtrl.DeleteDevice)
	router.GET("/cmdb/devicelist", deviceCtrl.ListDevices)
	router.GET("/cmdb/device/info", deviceCtrl.GetDevice)
	router.POST("/cmdb/device/connectivity", deviceCtrl.BatchTestDeviceConnectivity)

	// Cloud host management
	router.POST("/cmdb/hostcloudcreatealiyun", controller.NewCmdbHostCloudController().CreateAliyunHost)
	router.POST("/cmdb/hostcloudcreatetencent", controller.NewCmdbHostCloudController().CreateTencentHost)
	router.POST("/cmdb/hostcloudcreatebaidu", controller.NewCmdbHostCloudController().CreateBaiduHost)
	router.GET("/cmdb/hostssh/connect/:id", controller.NewCmdbHostSSHController(service.GetCmdbHostSSHService()).ConnectTerminal)
	router.GET("/cmdb/hostssh/command-risk/:id", controller.NewCmdbHostSSHController(service.GetCmdbHostSSHService()).PreviewCommandRisk)
	router.GET("/cmdb/hostssh/command/:id", controller.NewCmdbHostSSHController(service.GetCmdbHostSSHService()).ExecuteCommand)
	router.POST("/cmdb/hostssh/upload/:id", controller.NewCmdbHostSSHController(service.GetCmdbHostSSHService()).UploadFile)

	// SQL execution
	router.POST("/cmdb/sql/select", controller.GetCmdbSQLRecordController().ExecuteSelect)
	router.POST("/cmdb/sql", controller.GetCmdbSQLRecordController().ExecuteInsert)
	router.PUT("/cmdb/sql", controller.GetCmdbSQLRecordController().ExecuteUpdate)
	router.DELETE("/cmdb/sql", controller.GetCmdbSQLRecordController().ExecuteDelete)
	router.POST("/cmdb/sql/execute", controller.GetCmdbSQLRecordController().ExecuteSQL)
	router.POST("/cmdb/sql/databaselist", controller.GetCmdbSQLRecordController().ListDatabases)

	// SQL log management
	router.GET("/cmdb/sqlLog/list", controller.GetCmdbSqlLogList)
	router.DELETE("/cmdb/sqlLog/delete", controller.DeleteCmdbSqlLogById)
	router.DELETE("/cmdb/sqlLog/batch/delete", controller.BatchDeleteCmdbSqlLog)
	router.DELETE("/cmdb/sqlLog/clean", controller.CleanCmdbSqlLog)
	router.GET("/cmdb/sql/work-orders/summary", sqlWorkOrderCtrl.GetSummary)
	router.GET("/cmdb/sql/work-orders", sqlWorkOrderCtrl.List)
	router.GET("/cmdb/sql/work-orders/:id", sqlWorkOrderCtrl.Detail)
	router.POST("/cmdb/sql/work-orders", sqlWorkOrderCtrl.Create)
	router.POST("/cmdb/sql/work-orders/:id/approve", sqlWorkOrderCtrl.Approve)
	router.POST("/cmdb/sql/work-orders/:id/reject", sqlWorkOrderCtrl.Reject)
	router.POST("/cmdb/sql/work-orders/:id/execute", sqlWorkOrderCtrl.Execute)

	// Database management
	router.POST("/cmdb/database", controller.NewCmdbSQLController().CreateDatabase)
	router.PUT("/cmdb/database", controller.NewCmdbSQLController().UpdateDatabase)
	router.DELETE("/cmdb/database", controller.NewCmdbSQLController().DeleteDatabase)
	router.GET("/cmdb/database/info", controller.NewCmdbSQLController().GetDatabase)
	router.GET("/cmdb/databaselist", controller.NewCmdbSQLController().ListDatabases)
	router.GET("/cmdb/database/byname", controller.NewCmdbSQLController().GetDatabasesByName)
	router.GET("/cmdb/database/bytype", controller.NewCmdbSQLController().GetDatabasesByType)
}
