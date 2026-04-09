// pkg/db/migrate.go
package db

import (
	aimodel "dodevops-api/api/ai/model"
	appmodel "dodevops-api/api/app/model"
	cmdbmodel "dodevops-api/api/cmdb/model"
	ccmodel "dodevops-api/api/configcenter/model"
	k8smodel "dodevops-api/api/k8s/model"
	knowledgemodel "dodevops-api/api/knowledge/model"
	monitormodel "dodevops-api/api/monitor/model"
	systemmodel "dodevops-api/api/system/model"
	taskmodel "dodevops-api/api/task/model"
	toolmodel "dodevops-api/api/tool/model"
	"fmt"

	"gorm.io/gorm"
)

// 注册所有需要自动建表的 model
var models = []interface{}{
	&cmdbmodel.CmdbGroup{},
	&ccmodel.EcsAuth{},
	&ccmodel.KeyManage{},
	&ccmodel.SyncSchedule{},
	&cmdbmodel.CmdbHost{},
	&cmdbmodel.CmdbDevice{},
	&cmdbmodel.CmdbSQLRecord{},
	&cmdbmodel.CmdbSQL{},
	&cmdbmodel.CmdbSQLWorkOrder{},
	&ccmodel.AccountAuth{},
	&taskmodel.TaskTemplate{},
	&taskmodel.Task{},
	&taskmodel.TaskWork{},
	&taskmodel.TaskAnsible{},
	&taskmodel.TaskAnsibleWork{},
	&taskmodel.ConfigAnsible{},
	&monitormodel.Agent{},
	&monitormodel.MonitorHostAlertRuleEntity{},
	&monitormodel.MonitorDBAlertRuleEntity{},
	&monitormodel.MonitorAlertEventEntity{},
	&monitormodel.MonitorDBHealthSnapshotEntity{},
	&monitormodel.MonitorDomainEntity{},
	&monitormodel.MonitorDomainScheduleEntity{},
	&monitormodel.MonitorSSLCertEntity{},
	&monitormodel.MonitorSSLCertDeployLogEntity{},
	&k8smodel.KubeCluster{},
	&k8smodel.WorkloadCapacitySuggestionHistoryEntity{},
	&appmodel.Application{},
	&appmodel.JenkinsEnv{},
	&appmodel.QuickDeployment{},
	&appmodel.QuickDeploymentTask{},
	&aimodel.PromptTemplate{},
	&aimodel.AIChatHistory{},
	&aimodel.AIAssistantSessionContextEntity{},
	&aimodel.AIAssistantConfirmationEntity{},
	&aimodel.AIInspectionTemplateEntity{},
	&aimodel.AIInspectionReportEntity{},
	&knowledgemodel.KnowledgeBase{},
	&systemmodel.SysOperationLog{},
	&systemmodel.SysSessionRecording{},
	&systemmodel.SysCommandAudit{},
	&toolmodel.Tool{},
	&toolmodel.ServiceDeploy{},
	// 可以继续添加其他模型...
}

// 自动迁移所有模型
func AutoMigrate(db *gorm.DB) error {
	for _, item := range models {
		if err := db.AutoMigrate(item); err != nil {
			fmt.Printf("WARN: skip automigrate for %T: %v\n", item, err)
		}
	}
	return nil
}
