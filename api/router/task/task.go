package task

import (
	"dodevops-api/api/task/controller"
	"dodevops-api/api/task/service"
	"dodevops-api/common"
	"dodevops-api/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(router *gin.RouterGroup) {
	// 任务模板
	router.POST("/template/add", middleware.AuthMiddleware(), controller.CreateTemplate)
	router.GET("/template/list", middleware.AuthMiddleware(), controller.GetAllTemplates)
	router.PUT("/template/update", middleware.AuthMiddleware(), controller.UpdateTemplate)
	router.DELETE("/template/delete", middleware.AuthMiddleware(), controller.DeleteTemplate)
	router.GET("/template/info/:id", middleware.AuthMiddleware(), controller.GetTemplateByID)
	router.GET("/template/content/:id", middleware.AuthMiddleware(), controller.GetTemplateContent)
	router.GET("/template/query/name", middleware.AuthMiddleware(), controller.GetTemplatesByName)
	router.GET("/template/query/type", middleware.AuthMiddleware(), controller.GetTemplatesByType)

	// 通用任务
	router.POST("/task/add", middleware.AuthMiddleware(), controller.CreateTask)
	router.GET("/task/get", middleware.AuthMiddleware(), controller.GetTaskByID)
	router.PUT("/task/update", middleware.AuthMiddleware(), controller.UpdateTask)
	router.DELETE("/task/delete", middleware.AuthMiddleware(), controller.DeleteTask)
	router.GET("/task/list", middleware.AuthMiddleware(), controller.ListTasks)
	router.GET("/task/list-with-details", middleware.AuthMiddleware(), controller.ListTasksWithDetails)
	router.GET("/task/query/name", middleware.AuthMiddleware(), controller.GetTasksByName)
	router.GET("/task/query/type", middleware.AuthMiddleware(), controller.GetTasksByType)
	router.GET("/task/query/status", middleware.AuthMiddleware(), controller.GetTasksByStatus)
	router.GET("/task/next-execution", middleware.AuthMiddleware(), controller.GetNextExecutionTime)
	router.GET("/task/execution-info", middleware.AuthMiddleware(), controller.GetTaskExecutionInfo)
	router.GET("/task/templates", middleware.AuthMiddleware(), controller.GetTaskTemplatesWithStatus)

	// 任务作业
	router.POST("/taskjob/start", middleware.AuthMiddleware(), controller.TaskWork().StartJob)
	router.GET("/taskjob/log", middleware.AuthMiddleware(), controller.TaskWork().GetJobLog)
	router.POST("/taskjob/stop", middleware.AuthMiddleware(), controller.TaskWork().StopJob)
	router.GET("/taskjob/status", middleware.AuthMiddleware(), controller.TaskWork().GetJobStatus)

	// 任务监控
	taskMonitorCtrl := controller.NewTaskMonitorController()
	router.GET("/task/monitor/queue/metrics", middleware.AuthMiddleware(), taskMonitorCtrl.GetQueueMetrics)
	router.GET("/task/monitor/scheduler/stats", middleware.AuthMiddleware(), taskMonitorCtrl.GetSchedulerStats)
	router.GET("/task/monitor/system/status", middleware.AuthMiddleware(), taskMonitorCtrl.GetSystemStatus)
	router.GET("/task/monitor/queue/details", middleware.AuthMiddleware(), taskMonitorCtrl.GetQueueDetails)
	router.POST("/task/monitor/queue/clear-failed", middleware.AuthMiddleware(), taskMonitorCtrl.ClearFailedQueue)
	router.POST("/task/monitor/queue/retry-failed", middleware.AuthMiddleware(), taskMonitorCtrl.RetryFailedTasks)
	router.POST("/task/monitor/scheduled/pause", middleware.AuthMiddleware(), taskMonitorCtrl.PauseScheduledTask)
	router.POST("/task/monitor/scheduled/resume", middleware.AuthMiddleware(), taskMonitorCtrl.ResumeScheduledTask)
	router.POST("/task/monitor/scheduled/reset", middleware.AuthMiddleware(), taskMonitorCtrl.ResetScheduledTaskStatus)
	router.GET("/task/monitor/task/status", middleware.AuthMiddleware(), taskMonitorCtrl.GetTaskStatus)

	// Ansible 任务
	taskAnsibleCtrl := controller.NewTaskAnsibleController(service.NewTaskAnsibleService(common.GetDB()))
	router.GET("/task/ansiblelist", middleware.AuthMiddleware(), taskAnsibleCtrl.List)
	router.POST("/task/ansible", middleware.AuthMiddleware(), taskAnsibleCtrl.CreateTask)
	router.POST("/task/k8s", middleware.AuthMiddleware(), taskAnsibleCtrl.CreateK8sTask)
	router.GET("/task/ansible/:id", middleware.AuthMiddleware(), taskAnsibleCtrl.GetTask)
	router.GET("/task/ansible/:id/history", middleware.AuthMiddleware(), taskAnsibleCtrl.GetTaskHistory)
	router.GET("/task/ansible/:id/history/:historyId", middleware.AuthMiddleware(), taskAnsibleCtrl.GetTaskHistoryDetail)
	router.PUT("/task/ansible/:id", middleware.AuthMiddleware(), taskAnsibleCtrl.UpdateTask)
	router.POST("/task/ansible/:id/start", middleware.AuthMiddleware(), taskAnsibleCtrl.StartTask)
	router.DELETE("/task/ansible/:id", middleware.AuthMiddleware(), taskAnsibleCtrl.DeleteTask)
	router.GET("/task/ansible/:id/log/:work_id", middleware.AuthMiddleware(), taskAnsibleCtrl.GetJobLog)
	router.GET("/task/ansible/query/name", middleware.AuthMiddleware(), taskAnsibleCtrl.GetTasksByName)
	router.GET("/task/ansible/query/type", middleware.AuthMiddleware(), taskAnsibleCtrl.GetTasksByType)
}

// RegisterWebSocketRoutes registers task WebSocket routes.
func RegisterWebSocketRoutes(router *gin.RouterGroup) {
	wsCtrl := controller.NewWebSocketController(service.NewTaskAnsibleService(common.GetDB()))
	router.GET("/ws/task/ansible/:id/log/:work_id", wsCtrl.GetJobLogWS)
}
