package service

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	aiModel "dodevops-api/api/ai/model"
	appModel "dodevops-api/api/app/model"
	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	k8sModel "dodevops-api/api/k8s/model"
	monitorDao "dodevops-api/api/monitor/dao"
	monitorModel "dodevops-api/api/monitor/model"
	"dodevops-api/pkg/db"
	"github.com/gin-gonic/gin"
)

var (
	assistantTemplatePattern      = regexp.MustCompile(`(?i)(巡检模板|template)`)
	assistantReportArchivePattern = regexp.MustCompile(`(?i)(巡检报告列表|报告列表|report archive|历史报告)`)
	assistantAlertPattern         = regexp.MustCompile(`(?i)(告警|alert|incident|静默|alertmanager)`)
	assistantWorkloadPattern      = regexp.MustCompile(`(?i)(k8s|集群|cluster|namespace|命名空间|deployment|workload|pod|hpa|扩缩容)`)
	assistantWorkOrderPattern     = regexp.MustCompile(`(?i)(工单|sql工单|work order)`)
	assistantDeploymentPattern    = regexp.MustCompile(`(?i)(发布|deployment list|快速发布|release)`)
)

func (s *AIService) attachAssistantDefaults(response aiModel.AIAssistantChatResponse, sessionID string, userID uint, intent string) aiModel.AIAssistantChatResponse {
	response.Context = s.loadAssistantContext(sessionID, userID)
	response.AvailableTemplates = s.listInspectionTemplates(20)
	response.RecentReports = s.listRecentInspectionReports(userID, 8)
	if response.ToolSteps == nil {
		response.ToolSteps = []aiModel.AIAssistantToolStep{}
	}
	if response.Context != nil {
		response.Context.LastIntent = intent
	}
	return response
}

func appendAssistantToolStep(response *aiModel.AIAssistantChatResponse, tool, status, summary, detail string) {
	if response == nil {
		return
	}
	response.ToolSteps = append(response.ToolSteps, aiModel.AIAssistantToolStep{
		Tool:    tool,
		Status:  status,
		Summary: summary,
		Detail:  detail,
	})
}

func enrichContextWithHost(response *aiModel.AIAssistantChatResponse, host *aiModel.AIAssistantHost) {
	if response == nil || host == nil {
		return
	}
	if response.Context == nil {
		response.Context = &aiModel.AIAssistantContext{}
	}
	response.Context.CurrentScope = "host"
	response.Context.CurrentHostID = host.ID
	response.Context.CurrentHostName = host.HostName
	response.Context.Summary = fmt.Sprintf("当前聚焦主机 %s (%s)", host.HostName, firstNonEmpty(host.SSHIP, host.PrivateIP, host.PublicIP))
}

func (s *AIService) handleTemplateIntent(response aiModel.AIAssistantChatResponse, message string) aiModel.AIAssistantChatResponse {
	templates := s.listInspectionTemplates(20)
	response.AvailableTemplates = templates
	appendAssistantToolStep(&response, "inspection_template_list", "success", fmt.Sprintf("加载了 %d 个巡检模板", len(templates)), "")
	if len(templates) == 0 {
		response.AssistantMessage = "当前还没有可用的巡检模板。"
		return response
	}

	lines := []string{"当前可用巡检模板如下："}
	for _, item := range templates {
		lines = append(lines, fmt.Sprintf("- %s：%s", item.Name, item.Description))
	}
	response.AssistantMessage = strings.Join(lines, "\n")
	response.Actions = make([]aiModel.AIAssistantAction, 0, len(templates))
	for _, item := range templates {
		response.Actions = append(response.Actions, aiModel.AIAssistantAction{
			Label:   item.Name,
			Message: fmt.Sprintf("使用%s为当前主机生成巡检报告", item.Name),
			Kind:    "inspection",
		})
	}
	return response
}

func (s *AIService) handleReportArchiveIntent(response aiModel.AIAssistantChatResponse, userID uint) aiModel.AIAssistantChatResponse {
	reports := s.listRecentInspectionReports(userID, 12)
	response.RecentReports = reports
	appendAssistantToolStep(&response, "inspection_report_archive", "success", fmt.Sprintf("加载了 %d 份历史巡检报告", len(reports)), "")
	if len(reports) == 0 {
		response.AssistantMessage = "当前还没有历史巡检报告。先对一台主机执行巡检后，这里会自动沉淀报告。"
		return response
	}

	lines := []string{"最近巡检报告："}
	for _, item := range reports {
		lines = append(lines, fmt.Sprintf("- [%s] %s / %s", item.CreatedAt, item.TemplateName, item.TargetName))
	}
	response.AssistantMessage = strings.Join(lines, "\n")
	return response
}

func (s *AIService) handleAlertIntent(response aiModel.AIAssistantChatResponse, message string) aiModel.AIAssistantChatResponse {
	summary, err := monitorDao.GetMonitorAlertSummary()
	if err != nil {
		appendAssistantToolStep(&response, "alert_summary", "failed", "获取告警摘要失败", err.Error())
		response.AssistantMessage = "告警中心摘要获取失败，请稍后重试。"
		return response
	}
	appendAssistantToolStep(&response, "alert_summary", "success", "已获取告警中心摘要", "")

	incidentQuery := monitorModel.MonitorIncidentQuery{
		Keyword:  "",
		Status:   -1,
		PageSize: 5,
		PageNum:  1,
	}
	if response.Context != nil {
		if response.Context.CurrentNamespace != "" {
			incidentQuery.Namespace = response.Context.CurrentNamespace
		}
		if response.Context.CurrentWorkloadName != "" {
			incidentQuery.WorkloadName = response.Context.CurrentWorkloadName
		}
	}
	if strings.TrimSpace(message) != "" {
		incidentQuery.Keyword = extractAlertKeyword(message)
	}

	incidents, total, incErr := monitorDao.GetMonitorIncidentList(incidentQuery)
	if incErr == nil {
		appendAssistantToolStep(&response, "alert_incident_list", "success", fmt.Sprintf("匹配到 %d 条告警事件", total), "")
	} else {
		appendAssistantToolStep(&response, "alert_incident_list", "failed", "获取告警事件失败", incErr.Error())
	}

	lines := []string{
		"告警中心摘要：",
		fmt.Sprintf("- 总 incident：%d", summary.TotalIncidents),
		fmt.Sprintf("- 未恢复 incident：%d", summary.OpenIncidents),
		fmt.Sprintf("- 处理中 incident：%d", summary.ProcessingIncidents),
		fmt.Sprintf("- 严重 webhook：%d", summary.CriticalWebhookLogs),
	}
	if len(incidents) > 0 {
		lines = append(lines, "", "最近匹配到的告警：")
		for _, item := range incidents {
			lines = append(lines, fmt.Sprintf("- [%s] %s / %s", item.AlertLevel, item.BusinessLine, truncateForReport(item.AlertDesc, 80)))
		}
	}
	response.AssistantMessage = strings.Join(lines, "\n")
	if response.Context == nil {
		response.Context = &aiModel.AIAssistantContext{}
	}
	response.Context.CurrentScope = "alert"
	response.Context.Summary = fmt.Sprintf("当前告警中心有 %d 条未恢复事件", summary.OpenIncidents)
	return response
}

func extractAlertKeyword(message string) string {
	message = strings.TrimSpace(message)
	if message == "" {
		return ""
	}
	for _, keyword := range []string{"告警", "alert", "incident", "查看", "分析", "一下", "当前"} {
		message = strings.ReplaceAll(message, keyword, "")
	}
	return strings.TrimSpace(message)
}

func (s *AIService) handleWorkloadIntent(response aiModel.AIAssistantChatResponse, message string) aiModel.AIAssistantChatResponse {
	clusterName, namespaceName, workloadName := extractWorkloadHints(message, response.Context)
	appendAssistantToolStep(&response, "workload_hint_parse", "success", "已解析工作负载查询线索", fmt.Sprintf("cluster=%s namespace=%s workload=%s", clusterName, namespaceName, workloadName))

	var snapshots []k8sModel.WorkloadCapacitySuggestionHistoryEntity
	query := db.Db.Model(&k8sModel.WorkloadCapacitySuggestionHistoryEntity{})
	if clusterName != "" {
		query = query.Where("cluster_name LIKE ?", "%"+clusterName+"%")
	}
	if namespaceName != "" {
		query = query.Where("namespace_name LIKE ?", "%"+namespaceName+"%")
	}
	if workloadName != "" {
		query = query.Where("workload_name LIKE ?", "%"+workloadName+"%")
	}
	if err := query.Order("created_at DESC").Limit(5).Find(&snapshots).Error; err != nil {
		appendAssistantToolStep(&response, "workload_capacity_snapshot", "failed", "获取工作负载快照失败", err.Error())
		response.AssistantMessage = "暂时无法获取工作负载治理快照。"
		return response
	}
	appendAssistantToolStep(&response, "workload_capacity_snapshot", "success", fmt.Sprintf("获取到 %d 条工作负载快照", len(snapshots)), "")

	if len(snapshots) == 0 {
		var clusters []k8sModel.KubeCluster
		_ = db.Db.Order("updated_at DESC").Limit(5).Find(&clusters).Error
		lines := []string{"当前没有命中工作负载治理快照。"}
		if len(clusters) > 0 {
			lines = append(lines, "", "可用集群：")
			for _, item := range clusters {
				lines = append(lines, fmt.Sprintf("- %s（%s）", item.Name, item.GetStatusText()))
			}
		}
		response.AssistantMessage = strings.Join(lines, "\n")
		return response
	}

	lines := []string{"最近工作负载治理快照："}
	top := snapshots[0]
	for _, item := range snapshots {
		lines = append(lines, fmt.Sprintf("- [%s] %s/%s %s 风险=%s", item.CreatedAt.Format("2006-01-02 15:04:05"), item.ClusterName, item.NamespaceName, item.WorkloadName, item.RiskLevel))
	}
	response.AssistantMessage = strings.Join(lines, "\n")
	if response.Context == nil {
		response.Context = &aiModel.AIAssistantContext{}
	}
	response.Context.CurrentScope = "workload"
	response.Context.CurrentClusterID = top.ClusterID
	response.Context.CurrentClusterName = top.ClusterName
	response.Context.CurrentNamespace = top.NamespaceName
	response.Context.CurrentWorkloadType = top.WorkloadType
	response.Context.CurrentWorkloadName = top.WorkloadName
	response.Context.Summary = fmt.Sprintf("当前聚焦工作负载 %s/%s/%s", top.ClusterName, top.NamespaceName, top.WorkloadName)
	return response
}

func extractWorkloadHints(message string, context *aiModel.AIAssistantContext) (string, string, string) {
	clusterName := ""
	namespaceName := ""
	workloadName := ""
	if context != nil {
		clusterName = context.CurrentClusterName
		namespaceName = context.CurrentNamespace
		workloadName = context.CurrentWorkloadName
	}

	fields := strings.Fields(strings.NewReplacer("，", " ", ",", " ", "/", " ", "：", " ", ":", " ").Replace(message))
	for index, token := range fields {
		lower := strings.ToLower(token)
		if lower == "cluster" || token == "集群" {
			if index+1 < len(fields) {
				clusterName = fields[index+1]
			}
		}
		if lower == "namespace" || token == "命名空间" {
			if index+1 < len(fields) {
				namespaceName = fields[index+1]
			}
		}
		if lower == "deployment" || lower == "workload" || token == "工作负载" || token == "deployment" {
			if index+1 < len(fields) {
				workloadName = fields[index+1]
			}
		}
	}
	return strings.TrimSpace(clusterName), strings.TrimSpace(namespaceName), strings.TrimSpace(workloadName)
}

func (s *AIService) handleWorkOrderIntent(response aiModel.AIAssistantChatResponse, message string) aiModel.AIAssistantChatResponse {
	type workOrderDigest struct {
		Title   string
		Type    string
		Status  int
		Created string
	}

	var quickDeploys []appModel.QuickDeployment
	_ = db.Db.Order("created_at DESC").Limit(5).Find(&quickDeploys).Error

	var sqlOrders []struct {
		ID        uint
		Title     string
		OrderNo   string
		Status    int
		CreatedAt time.Time
	}
	_ = db.Db.Table("cmdb_sql_work_order").Select("id, title, order_no, status, created_at").Order("created_at DESC").Limit(5).Scan(&sqlOrders).Error

	entries := make([]workOrderDigest, 0, len(quickDeploys)+len(sqlOrders))
	for _, item := range quickDeploys {
		entries = append(entries, workOrderDigest{
			Title:   item.Title,
			Type:    "quick_deploy",
			Status:  item.Status,
			Created: item.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	for _, item := range sqlOrders {
		entries = append(entries, workOrderDigest{
			Title:   firstNonEmpty(item.OrderNo, item.Title),
			Type:    "sql_work_order",
			Status:  item.Status,
			Created: item.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	appendAssistantToolStep(&response, "work_order_digest", "success", fmt.Sprintf("加载了 %d 条工单摘要", len(entries)), "")
	if len(entries) == 0 {
		response.AssistantMessage = "当前没有可用的工单摘要。"
		return response
	}

	lines := []string{"最近工单摘要："}
	for _, item := range entries {
		lines = append(lines, fmt.Sprintf("- [%s] %s / status=%d / %s", item.Type, item.Title, item.Status, item.Created))
	}
	response.AssistantMessage = strings.Join(lines, "\n")
	if response.Context == nil {
		response.Context = &aiModel.AIAssistantContext{}
	}
	response.Context.CurrentScope = "workorder"
	response.Context.Summary = "当前聚焦工单中心摘要"
	return response
}

func (s *AIService) handleDeploymentIntent(response aiModel.AIAssistantChatResponse, message string) aiModel.AIAssistantChatResponse {
	var deployments []appModel.QuickDeployment
	query := db.Db.Model(&appModel.QuickDeployment{})
	keyword := strings.TrimSpace(strings.NewReplacer("发布", "", "deployment", "", "release", "").Replace(strings.ToLower(message)))
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if err := query.Order("created_at DESC").Limit(6).Find(&deployments).Error; err != nil {
		appendAssistantToolStep(&response, "deployment_list", "failed", "获取发布记录失败", err.Error())
		response.AssistantMessage = "暂时无法获取发布记录。"
		return response
	}

	appendAssistantToolStep(&response, "deployment_list", "success", fmt.Sprintf("加载了 %d 条发布记录", len(deployments)), "")
	if len(deployments) == 0 {
		response.AssistantMessage = "当前没有可用的发布记录。"
		return response
	}

	lines := []string{"最近发布记录："}
	top := deployments[0]
	for _, item := range deployments {
		lines = append(lines, fmt.Sprintf("- %s / status=%d / %s", item.Title, item.Status, item.CreatedAt.Format("2006-01-02 15:04:05")))
	}
	response.AssistantMessage = strings.Join(lines, "\n")
	if response.Context == nil {
		response.Context = &aiModel.AIAssistantContext{}
	}
	response.Context.CurrentScope = "deployment"
	response.Context.CurrentDeploymentID = top.ID
	response.Context.Summary = fmt.Sprintf("当前聚焦发布单 %s", top.Title)
	return response
}

func isRiskyAssistantAction(message string) bool {
	lower := strings.ToLower(message)
	for _, token := range []string{
		"重启", "restart", "reboot", "删除", "delete", "停止", "stop",
		"rollback", "回滚", "扩容", "缩容", "scale", "pause", "resume",
		"kill", "rm ", "systemctl restart", "systemctl stop", "kubectl delete",
	} {
		if strings.Contains(lower, token) {
			return true
		}
	}
	return false
}

func (s *AIService) handleRiskyIntent(c *gin.Context, response aiModel.AIAssistantChatResponse, sessionID string, userID uint, message string) aiModel.AIAssistantChatResponse {
	hosts, _ := resolveAssistantHosts(message)
	if len(hosts) == 1 {
		host := hosts[0]
		command := extractAssistantCommand(message)
		if command == "" {
			command = strings.TrimSpace(message)
		}
		pending := s.createPendingConfirmation(
			sessionID,
			userID,
			"host",
			"host_risky_command",
			host.ID,
			host.HostName,
			command,
			map[string]interface{}{
				"hostId":   host.ID,
				"hostName": host.HostName,
				"command":  command,
			},
			fmt.Sprintf("高风险操作待确认：将在主机 %s 上执行 `%s`", host.HostName, command),
		)
		response.PendingConfirmation = pending
		response.HostMatches = toAssistantHosts([]cmdbModel.CmdbHost{host})
		response.AssistantMessage = fmt.Sprintf("我识别到这是高风险操作，已经为你创建确认任务。在确认前我不会执行 `%s`。", command)
		appendAssistantToolStep(&response, "risky_action_guard", "pending", "已创建高风险确认任务", command)
		response.Actions = []aiModel.AIAssistantAction{
			{Label: "确认执行", Message: fmt.Sprintf("确认执行 %d", pending.ID), Kind: "confirm"},
			{Label: "取消任务", Message: fmt.Sprintf("取消执行 %d", pending.ID), Kind: "confirm"},
		}
		enrichContextWithHost(&response, &response.HostMatches[0])
		return response
	}

	response.AssistantMessage = "我识别到这是高风险操作，但当前还缺少明确目标或尚不支持自动执行。建议先定位到具体主机，再走确认流。"
	appendAssistantToolStep(&response, "risky_action_guard", "failed", "高风险操作未能进入确认流", "")
	return response
}

func resolveAssistantHostsWithContext(message string, context *aiModel.AIAssistantContext) ([]cmdbModel.CmdbHost, string) {
	hosts, hint := resolveAssistantHosts(message)
	if len(hosts) > 0 {
		return hosts, hint
	}
	if context != nil && context.CurrentHostID > 0 && (strings.Contains(message, "这台") || strings.Contains(message, "当前主机") || strings.Contains(message, "刚才那台")) {
		hostDao := cmdbDao.NewCmdbHostDao()
		host, err := hostDao.GetCmdbHostById(context.CurrentHostID)
		if err == nil {
			return []cmdbModel.CmdbHost{host}, "session_context"
		}
	}
	return hosts, hint
}
