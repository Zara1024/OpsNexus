package app

import (
	"dodevops-api/api/app/controller"
	"dodevops-api/common"
	"dodevops-api/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterApplicationRoutes(router *gin.RouterGroup) {
	appCtrl := controller.NewApplicationController(common.GetDB())
	workOrderCtrl := controller.NewWorkOrderController(common.GetDB())

	router.POST("/apps", middleware.AuthMiddleware(), appCtrl.CreateApplication)
	router.GET("/apps", middleware.AuthMiddleware(), appCtrl.GetApplicationList)
	router.GET("/apps/:id", middleware.AuthMiddleware(), appCtrl.GetApplicationDetail)
	router.PUT("/apps/:id", middleware.AuthMiddleware(), appCtrl.UpdateApplication)
	router.DELETE("/apps/:id", middleware.AuthMiddleware(), appCtrl.DeleteApplication)

	router.GET("/apps/:id/jenkins-envs", middleware.AuthMiddleware(), appCtrl.GetAppJenkinsEnvs)
	router.POST("/apps/:id/jenkins-envs", middleware.AuthMiddleware(), appCtrl.AddAppJenkinsEnv)
	router.PUT("/apps/:id/jenkins-envs/:env_id", middleware.AuthMiddleware(), appCtrl.UpdateAppJenkinsEnv)
	router.DELETE("/apps/:id/jenkins-envs/:env_id", middleware.AuthMiddleware(), appCtrl.DeleteAppJenkinsEnv)

	router.GET("/apps/jenkins-servers", middleware.AuthMiddleware(), appCtrl.GetJenkinsServers)
	router.POST("/apps/jenkins-job/validate", middleware.AuthMiddleware(), appCtrl.ValidateJenkinsJob)

	router.GET("/apps/service-tree", middleware.AuthMiddleware(), appCtrl.GetServiceTree)
	router.GET("/apps/business-group-options", middleware.AuthMiddleware(), appCtrl.GetBusinessGroupOptions)
	router.GET("/apps/environment", middleware.AuthMiddleware(), appCtrl.GetAppEnvironment)

	router.GET("/apps/deployment/applications", middleware.AuthMiddleware(), appCtrl.GetApplicationsForDeployment)
	router.POST("/apps/deployment/quick", middleware.AuthMiddleware(), appCtrl.CreateQuickDeployment)
	router.POST("/apps/deployment/execute", middleware.AuthMiddleware(), appCtrl.ExecuteQuickDeployment)
	router.GET("/apps/deployment/list", middleware.AuthMiddleware(), appCtrl.GetQuickDeploymentList)
	router.GET("/apps/deployment/:id", middleware.AuthMiddleware(), appCtrl.GetQuickDeploymentDetail)
	router.DELETE("/apps/deployment/:id", middleware.AuthMiddleware(), appCtrl.DeleteQuickDeployment)
	router.GET("/apps/deployment/tasks/:task_id/log", middleware.AuthMiddleware(), appCtrl.GetTaskBuildLog)
	router.GET("/apps/deployment/tasks/:task_id/status", middleware.AuthMiddleware(), appCtrl.GetTaskStatus)

	router.GET("/work-orders/summary", middleware.AuthMiddleware(), workOrderCtrl.GetSummary)
	router.GET("/work-orders", middleware.AuthMiddleware(), workOrderCtrl.GetList)
	router.GET("/work-orders/:type/:id", middleware.AuthMiddleware(), workOrderCtrl.GetDetail)
	router.POST("/work-orders/script", middleware.AuthMiddleware(), workOrderCtrl.CreateScript)
	router.POST("/work-orders/script/:id/approve", middleware.AuthMiddleware(), workOrderCtrl.ApproveScript)
	router.POST("/work-orders/script/:id/reject", middleware.AuthMiddleware(), workOrderCtrl.RejectScript)
	router.POST("/work-orders/script/:id/execute", middleware.AuthMiddleware(), workOrderCtrl.ExecuteScript)
}

func RegisterJenkinsRoutes(router *gin.RouterGroup) {
	jenkinsCtrl := controller.NewJenkinsController(common.GetDB())

	router.GET("/jenkins/servers", middleware.AuthMiddleware(), jenkinsCtrl.GetJenkinsServers)
	router.GET("/jenkins/servers/:id", middleware.AuthMiddleware(), jenkinsCtrl.GetJenkinsServerDetail)
	router.POST("/jenkins/test-connection", middleware.AuthMiddleware(), jenkinsCtrl.TestJenkinsConnection)

	router.GET("/jenkins/:serverId/jobs", middleware.AuthMiddleware(), jenkinsCtrl.GetJobs)
	router.GET("/jenkins/:serverId/jobs/search", middleware.AuthMiddleware(), jenkinsCtrl.SearchJobs)
	router.GET("/jenkins/:serverId/jobs/:jobName", middleware.AuthMiddleware(), jenkinsCtrl.GetJobDetail)
	router.GET("/jenkins/:serverId/jobs/:jobName/parameters", middleware.AuthMiddleware(), jenkinsCtrl.GetJobParameters)

	router.POST("/jenkins/:serverId/jobs/:jobName/start", middleware.AuthMiddleware(), jenkinsCtrl.StartJob)

	router.GET("/jenkins/:serverId/jobs/:jobName/builds/:buildNumber", middleware.AuthMiddleware(), jenkinsCtrl.GetBuildDetail)
	router.POST("/jenkins/:serverId/jobs/:jobName/builds/:buildNumber/stop", middleware.AuthMiddleware(), jenkinsCtrl.StopBuild)
	router.GET("/jenkins/:serverId/jobs/:jobName/builds/:buildNumber/log", middleware.AuthMiddleware(), jenkinsCtrl.GetBuildLog)

	router.GET("/jenkins/:serverId/system-info", middleware.AuthMiddleware(), jenkinsCtrl.GetSystemInfo)
	router.GET("/jenkins/:serverId/queue", middleware.AuthMiddleware(), jenkinsCtrl.GetQueueInfo)
}
