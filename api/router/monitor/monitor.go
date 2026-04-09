package monitor

import (
	"dodevops-api/api/monitor/controller"
	"dodevops-api/middleware"

	"github.com/gin-gonic/gin"
)

func InitMonitorRouter(r *gin.RouterGroup) {
	monitorController := controller.NewMonitorController()
	agentController := controller.NewAgentController()
	alertController := controller.NewMonitorAlertController()
	automationController := controller.NewMonitorAutomationController()

	monitorGroup := r.Group("/monitor")
	monitorGroup.Use(middleware.AuthMiddleware())

	// Host monitoring
	monitorGroup.GET("/host/:id", monitorController.GetHostMetrics)
	monitorGroup.GET("/hosts", monitorController.BatchGetHostMetrics)
	monitorGroup.GET("/hosts/:id/history", monitorController.GetHostMetricHistory)
	monitorGroup.GET("/hosts/:id/all-metrics", monitorController.GetHostAllMetricsHistory)
	monitorGroup.GET("/hosts/:id/top-processes", monitorController.GetTopProcesses)
	monitorGroup.GET("/hosts/:id/ports", monitorController.GetHostPorts)

	// Agent management
	monitorGroup.POST("/agent/deploy", agentController.DeployAgent)
	monitorGroup.DELETE("/agent/uninstall", agentController.UninstallAgent)
	monitorGroup.GET("/agent/status/:id", agentController.GetAgentStatus)
	monitorGroup.POST("/agent/restart/:id", agentController.RestartAgent)
	monitorGroup.GET("/agent/list", agentController.GetAgentList)
	monitorGroup.GET("/agent/statistics", agentController.GetAgentStatistics)
	monitorGroup.DELETE("/agent/delete/:id", agentController.DeleteAgent)

	// Alert center
	monitorGroup.GET("/alerts/summary", alertController.GetAlertSummary)
	monitorGroup.GET("/alerts/incidents", alertController.GetIncidentList)
	monitorGroup.GET("/alerts/webhooks", alertController.GetWebhookLogList)
	monitorGroup.GET("/alerts/notify/logs", alertController.GetNotifyLogList)
	monitorGroup.GET("/alerts/notify/robots", alertController.GetNotifyRobotList)
	monitorGroup.POST("/alerts/notify/robots", alertController.CreateNotifyRobot)
	monitorGroup.PUT("/alerts/notify/robots/:id", alertController.UpdateNotifyRobot)
	monitorGroup.POST("/alerts/notify/robots/:id/status", alertController.UpdateNotifyRobotStatus)
	monitorGroup.POST("/alerts/notify/robots/:id/test", alertController.TestNotifyRobot)
	monitorGroup.DELETE("/alerts/notify/robots/:id", alertController.DeleteNotifyRobot)
	monitorGroup.GET("/alerts/sources", alertController.GetAlertSourceList)
	monitorGroup.POST("/alerts/sources", alertController.CreateAlertSource)
	monitorGroup.PUT("/alerts/sources/:id", alertController.UpdateAlertSource)
	monitorGroup.POST("/alerts/sources/:id/status", alertController.UpdateAlertSourceStatus)
	monitorGroup.DELETE("/alerts/sources/:id", alertController.DeleteAlertSource)
	monitorGroup.GET("/alerts/alertmanager/status", alertController.GetAlertManagerStatus)
	monitorGroup.GET("/alerts/alertmanager/silences", alertController.GetAlertManagerSilenceList)
	monitorGroup.POST("/alerts/alertmanager/silences", alertController.CreateAlertManagerSilence)
	monitorGroup.DELETE("/alerts/alertmanager/silences/:silenceId", alertController.DeleteAlertManagerSilence)
	monitorGroup.GET("/alerts/alertmanager/receivers", alertController.GetAlertManagerReceiverList)

	// Monitor automation
	monitorGroup.GET("/automation/overview", automationController.GetAutomationOverview)
	monitorGroup.GET("/automation/events", automationController.GetAutomationEventList)
	monitorGroup.GET("/automation/host-alert/templates", automationController.GetHostAlertTemplates)
	monitorGroup.GET("/automation/host-alert/rules", automationController.GetHostAlertRuleList)
	monitorGroup.POST("/automation/host-alert/rules", automationController.CreateHostAlertRule)
	monitorGroup.PUT("/automation/host-alert/rules/:id", automationController.UpdateHostAlertRule)
	monitorGroup.POST("/automation/host-alert/rules/:id/status", automationController.UpdateHostAlertRuleStatus)
	monitorGroup.DELETE("/automation/host-alert/rules/:id", automationController.DeleteHostAlertRule)
	monitorGroup.POST("/automation/host-alert/scan", automationController.ScanHostAlerts)
	monitorGroup.GET("/automation/db-alert/rules", automationController.GetDBAlertRuleList)
	monitorGroup.POST("/automation/db-alert/rules", automationController.CreateDBAlertRule)
	monitorGroup.PUT("/automation/db-alert/rules/:id", automationController.UpdateDBAlertRule)
	monitorGroup.POST("/automation/db-alert/rules/:id/status", automationController.UpdateDBAlertRuleStatus)
	monitorGroup.DELETE("/automation/db-alert/rules/:id", automationController.DeleteDBAlertRule)
	monitorGroup.GET("/automation/db-alert/snapshots", automationController.GetDBHealthSnapshots)
	monitorGroup.POST("/automation/db-alert/scan", automationController.ScanDBAlerts)
	monitorGroup.GET("/automation/ssl/domains", automationController.GetSSLDomainList)
	monitorGroup.POST("/automation/ssl/domains", automationController.CreateSSLDomain)
	monitorGroup.PUT("/automation/ssl/domains/:id", automationController.UpdateSSLDomain)
	monitorGroup.POST("/automation/ssl/domains/:id/status", automationController.UpdateSSLDomainStatus)
	monitorGroup.DELETE("/automation/ssl/domains/:id", automationController.DeleteSSLDomain)
	monitorGroup.GET("/automation/ssl/schedules", automationController.GetSSLDomainScheduleList)
	monitorGroup.PUT("/automation/ssl/schedules/:id", automationController.UpdateSSLDomainSchedule)
	monitorGroup.POST("/automation/ssl/scan", automationController.ScanSSLDomains)
	monitorGroup.GET("/automation/ssl/certs", automationController.GetSSLCertList)
	monitorGroup.GET("/automation/ssl/deploy/logs", automationController.GetSSLCertDeployLogs)
	monitorGroup.POST("/automation/ssl/deploy", automationController.DeploySSLCert)
}
